package user

import (
	"encoding/json"
	"fmt"

	"github.com/jianshao/poker_counter/src/utils"
	"github.com/jianshao/poker_counter/src/view"
)

type PlayerInfo struct {
	// 基础信息，存在database中
	Id         int
	Name       string
	OpenId     string
	CurrRoomId int
	Rooms      map[int]*UserRoomInfo
}

// 用户在房间内的动态数据
type UserRoomInfo struct {
	Status     int
	CurrScore  int
	FinalScore int
	JoinTime   string
	ExitTime   string
	ApplyList  map[int]int
}

var (
	gUserMap = map[int]*PlayerInfo{}
)

const (
	USER_STATUS_WATCHING = 0
	USER_STATUS_PLAYING  = 1
	USER_STATUS_QUIT     = 2
)

func GetUser(userId int) *PlayerInfo {
	if user, ok := gUserMap[userId]; ok {
		return user
	}
	// TODO:
	return loadUser(userId)
}

// 从database中拉取用户数据
func loadUserFromData(userId int) (*PlayerInfo, error) {
	user, err := view.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	return &PlayerInfo{
		Id:    user.ID,
		Name:  user.Name,
		Rooms: map[int]*UserRoomInfo{},
	}, nil
}

func buildUserKey(userId int) string {
	return fmt.Sprintf("User:%d", userId)
}

func loadUserFromRedis(userId int) (*PlayerInfo, error) {
	key := buildUserKey(userId)
	userStr, err := utils.GetString(key)
	if err != nil {
		return nil, err
	}

	var player PlayerInfo
	err = json.Unmarshal([]byte(userStr), &player)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func setUser2Redis(user *PlayerInfo, timeout int) error {
	key := buildUserKey(user.Id)
	userStr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return utils.SetString(key, string(userStr), timeout)
}

// 载入完成需要保证，本地缓存、redis、database中都有相同的数据
func loadUser(userId int) *PlayerInfo {
	// 如果本地缓存有，则代表redis和database中有
	if user, ok := gUserMap[userId]; ok {
		return user
	}

	// redis中有，获取到之后需要保存到本地缓存
	user, err := loadUserFromRedis(userId)
	if err == nil {
		gUserMap[user.Id] = user
		return user
	}

	user, err = loadUserFromData(userId)
	if err == nil {
		// 保存到本地缓存和redis
		gUserMap[user.Id] = user
		setUser2Redis(user, 0)
	}
	return user
}
