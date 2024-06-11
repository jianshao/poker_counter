package records

import (
	"errors"

	"github.com/jianshao/poker_counter/prisma/db"
	"github.com/jianshao/poker_counter/src/view"
)

func buildApplyScore(apply *db.ScoreRecordsModel) *ApplyScore {
	record := &ApplyScore{
		Id:          apply.ID,
		RoomId:      apply.RoomID,
		UserId:      apply.UID,
		Score:       apply.Score,
		Status:      view.Status2int[apply.Status],
		ApplyType:   view.Type2int[apply.Type],
		ApplyTime:   apply.CreatedTime.String(),
		ConfirmTime: apply.UpdatedTime.String(),
	}
	if view.Status2int[apply.Status] != 1 {
		record.ConfirmTime = ""
	}
	return record
}

func ApplyBuyIn(roomId, userId, score, applyType int) (*ApplyScore, error) {
	// 往数据库中插入一条数据
	applyData, err := view.InsertScoreApply(roomId, userId, score, applyType)
	if err != nil {
		return nil, err
	}

	apply := buildApplyScore(applyData)
	addApply(apply.Id, apply)

	return apply, nil
}

func ConfirmBuyIn(applyId, status int) (*ApplyScore, error) {
	apply, err := GetApply(applyId)
	if err != nil {
		return nil, err
	}
	if apply.Status != 0 {
		return nil, errors.New("apply status error")
	}
	// 更新数据库
	newApply, err := view.UpdateScoreApply(applyId, status)
	if err != nil {
		return nil, err
	}

	// 更新本地缓存
	apply.Status = status
	apply.ConfirmTime = newApply.UpdatedTime.String()
	return apply, nil
}

func GetApplyScoreAll(roomId int) ([]ApplyScore, error) {
	records, err := view.GetScoreRecords(roomId, 0)
	if err != nil {
		return nil, err
	}

	applies := []ApplyScore{}
	for _, record := range records {
		applies = append(applies, *buildApplyScore(&record))
	}
	return applies, nil
}
