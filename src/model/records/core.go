package records

import (
	"errors"

	"github.com/jianshao/poker_counter/src/view"
)

type ApplyScore struct {
	Id          int
	RoomId      int
	UserId      int
	Name        string
	Score       int
	Status      int
	ApplyType   int
	ApplyTime   string
	ConfirmTime string
}

var (
	gAppliesMap = map[int]*ApplyScore{}
)

func Init() {

}

func GetApply(applyId int) (*ApplyScore, error) {
	if apply, ok := gAppliesMap[applyId]; ok {
		return apply, nil
	}
	// TODO:
	return loadApply(applyId)
}

func addApply(applyId int, apply *ApplyScore) error {
	if _, ok := gAppliesMap[applyId]; !ok {
		gAppliesMap[applyId] = apply
		return nil
	}
	return errors.New("apply already existed")
}

func loadApply(applyId int) (*ApplyScore, error) {
	apply, err := view.GetScoreRecordById(applyId)
	if err != nil {
		return nil, err
	}
	return buildApplyScore(apply), nil
}
