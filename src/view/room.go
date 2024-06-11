package view

// 所有针对数据库的操作都放在本文件内

import (
	"context"
	"errors"

	"github.com/jianshao/poker_counter/prisma/db"
	"github.com/jianshao/poker_counter/src/utils"
)

var (
	RoomStatus2int = map[db.RoomStatus]int{
		"OPEN":   0,
		"CLOSED": 1,
	}
	int2RoomStatus = map[int]db.RoomStatus{
		0: "OPEN",
		1: "CLOSED",
	}
)

// 往数据库中创建一个房间
func CreateOneRoom(roomId, owner int) (*db.RoomModel, error) {
	if roomId == 0 || owner == 0 {
		return nil, errors.New("params error")
	}
	client := utils.GetPrismaClient()

	room, err := client.Room.CreateOne(
		db.Room.RoomID.Set(roomId),
		db.Room.Owner.Set(owner),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return room, nil
}

func GetRoomByRoomId(roomId, status int) (*db.RoomModel, error) {
	client := utils.GetPrismaClient()
	return client.Room.FindFirst(
		db.Room.RoomID.Equals(roomId),
		db.Room.Status.Equals(int2RoomStatus[status]),
	).Exec(context.Background())
}

func GetOpenRoom(owner int) (*db.RoomModel, error) {
	client := utils.GetPrismaClient()
	return client.Room.FindFirst(
		db.Room.Owner.Equals(owner),
		db.Room.Status.Equals("OPEN"),
	).Exec(context.Background())
}

func CloseRoom(roomId, owner int) error {
	client := utils.GetPrismaClient()
	client.Room.FindMany(
		db.Room.Owner.Equals(owner),
		db.Room.Status.Equals("OPEN"),
		db.Room.RoomID.Equals(roomId),
	).Update(
		db.Room.Status.Set("CLOSED"),
	).Exec(context.Background())
	return nil
}

func GetAllOpenRooms() ([]db.RoomModel, error) {
	client := utils.GetPrismaClient()
	return client.Room.FindMany(
		db.Room.Status.Equals("OPEN"),
	).Exec(context.Background())
}
