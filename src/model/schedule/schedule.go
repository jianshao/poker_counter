package schedule

import (
	"fmt"
	"time"

	"github.com/jianshao/poker_counter/src/model/room"
	"github.com/jianshao/poker_counter/src/utils/logs"
	"github.com/jianshao/poker_counter/src/utils/schedule"
)

func AddSchedule(name, desc string, handler func() error, firstProTime time.Time, interval int, scheduleType int) {
	id, err := schedule.AddSchedule(&schedule.Schedule{
		Name:         name,
		Desc:         desc,
		Handler:      handler,
		Type:         scheduleType,
		FirstProTime: firstProTime,
		Interval:     int64(interval),
	})
	if err != nil {
		logs.Error(nil, fmt.Sprintf("add schedule failed: %s", err.Error()))
	} else {
		logs.Info(nil, fmt.Sprintf(" add schedule %d", id))
		err = schedule.StartSchedule(id)
		if err != nil {
			logs.Error(nil, fmt.Sprintf("start schedule failed: %s", err.Error()))
		}
	}
}

func Init() error {
	// 需要使用单独的协程来执行定时任务
	go func() {
		// 初始化任务调度
		schedule.Init()

		// 增加房间定时清理任务：每天凌晨1点执行，清理超过3天未使用的且未关闭的房间
		processTime := time.Date(2024, 6, 18, 11, 50, 0, 0, time.Local)
		desc := "每天凌晨1点执行,清理超过3天未使用的且未关闭的房间"
		AddSchedule("清理房间", desc, room.ClearUnusedRooms, processTime, 24*3600, schedule.SCHEDULE_TYPE_INTERVAL)

		// 运行任务
		schedule.Run()
	}()
	return nil
}
