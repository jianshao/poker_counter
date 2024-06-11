package room

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jianshao/poker_counter/src/utils"
	"github.com/jianshao/poker_counter/src/view"
)

type RoomInfo struct {
	RoomId  int
	Owner   int
	Status  int
	Players map[int]int
}

const (
	RoomStatus_Open = iota
	RoomStatus_Close
)

var (
	gRoomMap = map[int]*RoomInfo{}
)

func initRoom() error {
	// 为服务初始化资源
	gRoomMap = map[int]*RoomInfo{}
	// 从redis获取房间号开始位置
	return nil
}

func getRoom(roomId int) *RoomInfo {
	if roomInfo, ok := gRoomMap[roomId]; ok {
		return roomInfo
	}
	// TODO: 如果大量访问不存在的房间会导致资源浪费
	// 执行重载数据，有可能进程重启导致缓存数据丢失
	room, _ := loadRoom(roomId)
	return room
}

// 获取活跃房间
func getActiveRoom(roomId int) *RoomInfo {
	room := getRoom(roomId)
	if room != nil && room.Status == RoomStatus_Open {
		return room
	}
	return nil
}

func IsOwner(roomId, userId int) bool {
	roomInfo := getActiveRoom(roomId)
	if roomInfo.Owner == userId {
		return true
	}
	return false
}

// 生成房间号，生成规则：以星期为周期，每周从1开始，每次生成房间号，房间号为星期*1000+递增的房间号
func generateRoomId() int {
	weekday := int(time.Now().Weekday()) + 1
	roomId, err := utils.Inc("room_Id")
	if err != nil {
		return 0
	}
	return weekday*1000 + roomId%1000
}

// 载入房间信息
func loadRoom(roomId int) (*RoomInfo, error) {
	// 先从本地缓存获取
	if room, ok := gRoomMap[roomId]; ok {
		return room, nil
	}

	// 再从redis中获取
	room, err := loadRoomFromRedis(roomId)
	if err != nil {
		if err != redis.ErrNil {
			return nil, err
		}
	} else if room != nil {
		gRoomMap[roomId] = room
		return room, nil
	}

	// 从数据库中获取房间信息
	room, err = loadRoomFromData(roomId)
	if err != nil {
		// 如果数据不存在会进入这个分支
		return nil, err
	} else {
		gRoomMap[roomId] = room
		setRoom2Redis(room)
	}
	return room, nil
}

func loadRoomFromData(roomId int) (*RoomInfo, error) {
	room, err := view.GetRoomByRoomId(roomId, 0)
	if err != nil {
		return nil, err
	}
	return &RoomInfo{
		RoomId:  room.RoomID,
		Owner:   room.Owner,
		Status:  view.RoomStatus2int[room.Status],
		Players: map[int]int{},
	}, nil
}

func buildRoomKey(roomId int) string {
	return fmt.Sprintf("room:%d", roomId)
}

// 从redis载入房间信息
func loadRoomFromRedis(roomId int) (*RoomInfo, error) {
	// 先访问redis，看有没有该房间
	roomKey := buildRoomKey(roomId)
	roomInfo, err := utils.GetString(roomKey)
	if err != nil {
		return nil, err
	}

	var room RoomInfo
	err = json.Unmarshal([]byte(roomInfo), &room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func setRoom2Redis(room *RoomInfo) error {
	key := buildRoomKey(room.RoomId)
	roomStr, err := json.Marshal(room)
	if err != nil {
		return err
	}
	return utils.SetString(key, string(roomStr))
}
