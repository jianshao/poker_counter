package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/jianshao/poker_counter/src/model/records"
	"github.com/jianshao/poker_counter/src/view"
)

func UserCheck(openId string) (*PlayerInfo, error) {
	// 直接查询数据库，以确定用户是否存在
	user, err := view.GetUserInfoByOpenid(openId)
	if err != nil {
		return nil, err
	}
	return &PlayerInfo{
		Id:   user.ID,
		Name: user.Name,
	}, nil
}

func UserLogin(userId int) *PlayerInfo {
	// 将用户的信息载入到进程中，以备后面使用
	return loadUser(userId)
}

func UserRegister(name, openId string) (*PlayerInfo, error) {
	// 在database中插入一条记录
	user, err := view.CreateOneUser(name, openId)
	if err != nil {
		return nil, err
	}

	return loadUser(user.ID), nil
}

func EntryRoom(roomId, userId int) error {
	// 先检查用户是否存在
	user := GetUser(userId)
	if user == nil {
		return errors.New("user not exist")
	}

	// 用户已经在该房间内了
	if user.CurrRoomId == roomId {
		return nil
	}

	// 用户已经在其他房间玩游戏了
	if user.CurrRoomId != 0 && user.Rooms[user.CurrRoomId].Status == USER_STATUS_PLAYING {
		return errors.New(fmt.Sprintf("已经在游戏中，请先在房间%d中退出游戏", user.CurrRoomId))
	}

	// 将用户的房间信息更新到本地缓存
	user.CurrRoomId = roomId
	// 如果用户的房间信息不存在，就初始化一下
	if _, ok := user.Rooms[roomId]; !ok {
		user.Rooms[roomId] = &UserRoomInfo{
			ApplyList: make(map[int]int),
		}
	}
	setUser2Redis(user, 0)

	return nil
}

func LeaveRoom(roomId, userId int) error {
	// 先检查用户是否存在
	user := GetUser(userId)
	if user == nil {
		return errors.New("user not exist")
	}

	// 用户当前不在任何房间
	if user.CurrRoomId == 0 {
		return nil
	}

	// 退出的房间号不对
	if user.CurrRoomId != roomId {
		return errors.New("user not in this room")
	}

	if user.Rooms[user.CurrRoomId].Status == USER_STATUS_PLAYING {
		return errors.New("user is playing, quit first")
	}

	user.CurrRoomId = 0
	setUser2Redis(user, 0)
	return nil
}

func JoinGame(roomId, userId int) error {
	// 先检查用户是否存在
	user := GetUser(userId)
	if user == nil {
		return errors.New("user not exist")
	}

	if user.CurrRoomId != roomId {
		return errors.New("user not in this room")
	}

	if user.Rooms[user.CurrRoomId].Status == USER_STATUS_PLAYING {
		return nil
	}

	user.Rooms[user.CurrRoomId].Status = USER_STATUS_PLAYING
	if user.Rooms[user.CurrRoomId].JoinTime == "" {
		user.Rooms[user.CurrRoomId].JoinTime = time.Now().Format("2006-01-02 15:04:05")
	}
	setUser2Redis(user, 0)
	return nil
}

func QuitGame(roomId, userId int) error {
	user := GetUser(userId)
	if user == nil {
		return errors.New("user not exist")
	}

	if user.CurrRoomId != roomId {
		return errors.New("user not in this room")
	}

	if user.Rooms[user.CurrRoomId].Status != USER_STATUS_PLAYING {
		return errors.New("user not playing")
	}

	user.Rooms[user.CurrRoomId].Status = USER_STATUS_QUIT
	user.Rooms[user.CurrRoomId].ExitTime = time.Now().Format("2006-01-02 15:04:05")
	setUser2Redis(user, 0)
	return nil
}

func addName2Apply(apply *records.ApplyScore) {
	user := GetUser(apply.UserId)
	apply.Name = user.Name
}

func ApplyBuyIn(roomId, userId, score, applyType int) (*records.ApplyScore, error) {
	user := GetUser(userId)
	if user == nil {
		return nil, errors.New("user not exist")
	}

	if user.CurrRoomId != roomId {
		return nil, errors.New("user not in this room")
	}

	if user.Rooms[user.CurrRoomId].Status != USER_STATUS_PLAYING {
		return nil, errors.New("user not playing")
	}

	apply, err := records.ApplyBuyIn(roomId, userId, score, applyType)
	if err != nil {
		return nil, err
	}

	user.Rooms[user.CurrRoomId].ApplyList[apply.Id] = apply.Id
	setUser2Redis(user, 0)

	addName2Apply(apply)
	return apply, nil
}

func ConfirmBuyIn(applyId, status int) (*records.ApplyScore, error) {
	apply, err := records.ConfirmBuyIn(applyId, status)
	if err != nil {
		return nil, err
	}

	// 更新用户的分数状态
	user := GetUser(apply.UserId)
	// 申请类型：0-申请买入，1-申请结算
	if apply.ApplyType == 0 {
		// 确认状态：0-未确认，1-同意，2-拒绝
		if status == 1 {
			user.Rooms[apply.RoomId].CurrScore += apply.Score
		}
	} else {
		if status == 1 {
			user.Rooms[apply.RoomId].FinalScore += apply.Score
		}
	}

	setUser2Redis(user, 0)
	addName2Apply(apply)
	return apply, nil
}

func GetAllScoreApplies(roomId int) ([]records.ApplyScore, error) {
	applies, err := records.GetApplyScoreAll(roomId)
	if err != nil {
		return nil, err
	}

	for i, _ := range applies {
		addName2Apply(&applies[i])
	}
	return applies, nil
}

func ClearUnusedRooms(users, rooms map[int]int) {
	for userId, _ := range users {
		user := GetUser(userId)
		if user == nil {
			continue
		}

		for roomId, _ := range rooms {
			if _, ok := user.Rooms[roomId]; ok {
				delete(user.Rooms, roomId)
				if user.CurrRoomId == roomId {
					user.CurrRoomId = 0
				}
			}
		}
		setUser2Redis(user, 0)
	}
	return
}
