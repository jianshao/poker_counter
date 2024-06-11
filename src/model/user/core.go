package user

import (
	"encoding/json"
	"fmt"

	"github.com/jianshao/poker_counter/src/utils"
	"github.com/jianshao/poker_counter/src/view"
)

type PlayerInfo struct {
	// 基础信息，存在database中
	Id     int
	Name   string
	OpenId string
	// 运行时数据，存在本地缓存或redis
	CurrRoomId int
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

func GetUser(userId int) *PlayerInfo {
	if user, ok := gUserMap[userId]; ok {
		return user
	}
	// TODO:
	return loadUser(userId)
}

func GetActiveUser(userId int) *PlayerInfo {
	user := GetUser(userId)
	if user.Status == 1 {
		return user
	}
	return nil
}

// 从database中拉取用户数据
func loadUserFromData(userId int) (*PlayerInfo, error) {
	user, err := view.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	return &PlayerInfo{
		Id:        user.ID,
		Name:      user.Name,
		ApplyList: make(map[int]int),
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

func setUser2Redis(user *PlayerInfo) error {
	key := buildUserKey(user.Id)
	userStr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return utils.SetString(key, string(userStr))
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
		setUser2Redis(user)
	}
	return user
}
