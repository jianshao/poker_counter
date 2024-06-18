package view

import (
	"context"

	"github.com/jianshao/poker_counter/prisma/db"
	"github.com/jianshao/poker_counter/src/utils"
)

var (
	int2Status = map[int]db.ScoreRecordStatus{
		0: "APPLY",
		1: "ACCEPT",
		2: "REJECT",
	}
	Status2int = map[db.ScoreRecordStatus]int{
		"APPLY":  0,
		"ACCEPT": 1,
		"REJECT": 2,
	}
	int2Type = map[int]db.ScoreRecordType{
		0: "BUYIN",
		1: "CASHOUT",
	}
	Type2int = map[db.ScoreRecordType]int{
		"BUYIN":   0,
		"CASHOUT": 1,
	}
)

func InsertScoreApply(roomId, userId, score, applyType int) (*db.ScoreRecordsModel, error) {
	client := utils.GetPrismaClient()
	return client.ScoreRecords.CreateOne(
		db.ScoreRecords.UID.Set(userId),
		db.ScoreRecords.RoomID.Set(roomId),
		db.ScoreRecords.Score.Set(score),
		db.ScoreRecords.Type.Set(int2Type[applyType]),
	).Exec(context.Background())
}

func UpdateScoreApply(applyId, status int) (*db.ScoreRecordsModel, error) {
	client := utils.GetPrismaClient()
	return client.ScoreRecords.FindUnique(
		db.ScoreRecords.ID.Equals(applyId),
	).Update(
		db.ScoreRecords.Status.Set(int2Status[status]),
	).Exec(context.Background())
}

func GetScoreRecords(roomId, status int) ([]db.ScoreRecordsModel, error) {
	client := utils.GetPrismaClient()
	return client.ScoreRecords.FindMany(
		db.ScoreRecords.Status.Equals(int2Status[status]),
		db.ScoreRecords.RoomID.Equals(roomId),
	).OrderBy(db.ScoreRecords.UpdatedTime.Order(db.SortOrderDesc)).Exec(context.Background())
}

func GetScoreRecordById(id int) (*db.ScoreRecordsModel, error) {
	client := utils.GetPrismaClient()
	return client.ScoreRecords.FindUnique(
		db.ScoreRecords.ID.Equals(id),
	).Exec(context.Background())
}
