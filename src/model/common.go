package model

import "github.com/jianshao/poker_counter/src/model/schedule"

func Init() error {
	schedule.Init()
	return nil
}
