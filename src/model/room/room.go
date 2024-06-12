package room

import (
	"errors"
	"fmt"

	"github.com/jianshao/poker_counter/prisma/db"
	"github.com/jianshao/poker_counter/src/model/records"
	"github.com/jianshao/poker_counter/src/model/user"
	"github.com/jianshao/poker_counter/src/view"
)

const (
	INVALID_ROOM_ID = 0
)

// 1. owner create room
func CreateRoom(userId int) (*RoomInfo, error) {
	// 用户在同一时间只能创建一个房间内
	room, err := view.GetOpenRoom(userId)
	if err == db.ErrNotFound {
		// 在数据库中创建一个房间
		roomId := generateRoomId()
		if _, err := view.CreateOneRoom(roomId, userId); err != nil {
			return nil, err
		}

		// 将房间信息载入进程
		return loadRoom(roomId)
	}

	if err != nil {
		return nil, err
	}
	return nil, errors.New(fmt.Sprintf("已经创建房间：%d", room.RoomID))
}

func CheckRoom(roomId int) *RoomInfo {
	return getActiveRoom(roomId)
}

func CloseRoom(roomId, userId int) error {
	room := getActiveRoom(roomId)
	if room == nil {
		return nil
	}

	if room.Owner != userId {
		return errors.New("only room owner can close room")
	}

	count := 0
	for _, playerId := range room.Players {
		if user.IsUserPlaying(roomId, playerId) {
			count += 1
		}
	}
	if count > 0 {
		return errors.New(fmt.Sprintf("仍有%d个玩家没有提交剩余积分数量，房间不可关闭。", count))
	}

	room.Status = RoomStatus_Close
	// 设置过期时间,防止长时间占用
	setRoom2Redis(room, 24*3600)
	// 更新数据库
	view.CloseRoom(roomId, userId)
	return nil
}

// 不在任何房间的用户才能进入指定房间
func EntryRoom(roomId, userId int) (bool, error) {
	// 先检查房间是否活跃
	room := getActiveRoom(roomId)
	if room == nil {
		return false, errors.New("room not activc")
	}

	// 构建下层数据
	err := user.EntryRoom(roomId, userId)
	if err != nil {
		return false, err
	}

	// 构建本层数据
	room.Players[userId] = userId
	setRoom2Redis(room, 0)
	return true, nil
}

// 进入房间之后才能加入该房间的游戏
func JoinGame(roomId, userId int) (bool, error) {
	// 先检查房间是否活跃
	roomInfo := getActiveRoom(roomId)
	if roomInfo == nil {
		return false, errors.New("room not existed")
	}

	if _, ok := roomInfo.Players[userId]; !ok {
		return false, errors.New("user not in this room")
	}

	if err := user.JoinGame(roomId, userId); err != nil {
		return false, err
	}
	return true, nil
}

func QuitGame(roomId, userId int) (bool, error) {
	roomInfo := getActiveRoom(roomId)
	if roomInfo == nil {
		return false, errors.New("room not existed")
	}

	if err := user.QuitGame(roomId, userId); err != nil {
		return false, err
	}

	return true, nil
}

// 退出房间不会导致数据变化
func LeaveRoom(roomId, userId int) (bool, error) {
	// 房间不是活跃状态
	room := getActiveRoom(roomId)
	if room == nil {
		return true, nil
	}

	// 清理下层数据
	if err := user.LeaveRoom(roomId, userId); err != nil {
		return false, err
	}

	// 清理本层数据
	// room.Players[userId] = userId
	delete(room.Players, userId)
	setRoom2Redis(room, 0)
	return true, nil
}

func GetRoomInfo(roomId int) (*RoomInfo, error) {
	roomInfo := getActiveRoom(roomId)
	if roomInfo == nil {
		return nil, errors.New("room not existed")
	}
	return roomInfo, nil
}

func ApplyBuyIn(roomId, userId, score, applyType int) (*records.ApplyScore, error) {
	room := getActiveRoom(roomId)
	if room == nil {
		return nil, errors.New("room not exist")
	}

	return user.ApplyBuyIn(roomId, userId, score, applyType)
}

func ConfirmBuyIn(roomId, owner, applyId, status int) (*records.ApplyScore, error) {
	room := getActiveRoom(roomId)
	if room == nil {
		return nil, errors.New("room not exist")
	}
	if room.Owner != owner {
		return nil, errors.New("only room owner can confirm applies")
	}
	return user.ConfirmBuyIn(applyId, status)
}

func GetAllScoreApplies(roomId int) ([]records.ApplyScore, error) {
	room := getActiveRoom(roomId)
	if room == nil {
		return nil, errors.New("room not exist")
	}

	return user.GetAllScoreApplies(roomId)
}
