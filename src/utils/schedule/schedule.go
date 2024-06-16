package schedule

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jianshao/poker_counter/src/utils/logs"
)

const (
	SCHEDULE_TYPE_FIXED    = 0
	SCHEDULE_TYPE_INTERVAL = 1

	SCHEDULE_STATUS_INIT = 0
	SCHEDULE_STATUS_RUN  = 1
	SCHEDULE_STATUS_STOP = 2
)

// 定时任务组件，可以指定固定时间或每隔一段时间执行
type ScheduleHandler func() error

type Schedule struct {
	Id           int64
	Name         string
	Desc         string
	Status       int // 状态，0-停止，1-运行
	Handler      ScheduleHandler
	Type         int       // 执行类型，0-固定时间（默认值），1-每隔一段时间
	FirstProTime time.Time // 首次执行时间
	Interval     int64     // 间隔时间，单位秒
	NextProTime  time.Time // 下次执行时间
}

type scheduleNode struct {
	pre  *scheduleNode
	next *scheduleNode
	data *Schedule
}

type ScheduleMgr struct {
	SchedulesMap  map[int64]*scheduleNode // map用于快速查找
	SchedulesList *scheduleNode           // 有序列表用于按照时间顺序执行
	StopRunning   chan bool               // 停止运行信号
	lock          sync.Mutex              // 锁保障并发安全
}

var (
	gManager *ScheduleMgr = nil
)

func getSchedule(id int64) (*scheduleNode, error) {
	sche, ok := gManager.SchedulesMap[id]
	if ok {
		return sche, nil
	} else {
		return nil, errors.New("not exist")
	}
}

// 将nodesche插入到合适位置，保证列表有序
func addNode2List(node *scheduleNode, list *scheduleNode) *scheduleNode {
	// 如果列表为空，则直接插入
	if list == nil {
		return node
	}

	head := list
	// 如果插入的节点时间在列表中最小的时间之前，则插入到列表头部
	if head.data.NextProTime.Compare(node.data.NextProTime) == 1 {
		node.next = head
		head.pre = node
		return node
	}

	// 正常插入
	currNode := head
	for currNode.next != nil {
		if node.data.NextProTime.Compare(currNode.next.data.NextProTime) == -1 {
			break
		}
		currNode = currNode.next
	}
	if currNode.next != nil {
		node.next = currNode.next
		node.pre = currNode
		currNode.next.pre = node
		currNode.next = node
	} else {
		currNode.next = node
		node.pre = currNode
	}
	return head
}

func delNodeFromList(node *scheduleNode, list *scheduleNode) *scheduleNode {
	head := list
	if head == node {
		if node.next == nil {
			return nil
		}
		node.next.pre = nil
		return node.next
	}

	if node.pre != nil {
		node.pre.next = node.next
	}
	if node.next != nil {
		node.next.pre = node.pre
	}
	return head
}

// 将sche插入到列表中
func addNode(sche *Schedule) error {
	if gManager == nil {
		return errors.New("schedule not init")
	}

	// 相同ID的任务已经存在，不能重复添加
	_, ok := gManager.SchedulesMap[sche.Id]
	if ok {
		return fmt.Errorf("schedule %d already exist", sche.Id)
	}

	node := &scheduleNode{data: sche}
	gManager.lock.Lock()
	gManager.SchedulesMap[sche.Id] = node
	gManager.SchedulesList = addNode2List(node, gManager.SchedulesList)
	gManager.lock.Unlock()
	return nil
}

func removeNode(id int64) error {
	node, err := getSchedule(id)
	if err == nil {
		gManager.lock.Lock()
		// 从map中删除
		delete(gManager.SchedulesMap, id)
		// 从列表中删除
		gManager.SchedulesList = delNodeFromList(node, gManager.SchedulesList)
		gManager.lock.Unlock()
	}
	return err
}

func getFirstNode() *scheduleNode {
	if gManager == nil || gManager.SchedulesList == nil {
		return nil
	}
	return gManager.SchedulesList
}

func afterProc(node *scheduleNode) {
	// 将当前任务节点后移到合适位置，移动后仍然保证列表有序
	gManager.lock.Lock()
	gManager.SchedulesList = delNodeFromList(node, gManager.SchedulesList)
	// 如果是固定时间，则不需要更新
	if node.data.Type == SCHEDULE_TYPE_INTERVAL {
		node.data.NextProTime = node.data.NextProTime.Add(time.Second * time.Duration(node.data.Interval))
		gManager.SchedulesList = addNode2List(node, gManager.SchedulesList)
	} else {
		node.data.Status = SCHEDULE_STATUS_STOP
	}
	gManager.lock.Unlock()
}

func doProc(schedule *Schedule) {
	// 检查状态是否可执行
	if schedule.Status != SCHEDULE_STATUS_RUN {
		logs.Info(nil, fmt.Sprintf("schedule %d not run, skipped", schedule.Id))
	} else {
		err := schedule.Handler()
		if err != nil {
			logs.Info(nil, fmt.Sprintf("schedule %d run failed: %s", schedule.Id, err.Error()))
		} else {
			logs.Info(nil, fmt.Sprintf("schedule %d run success ", schedule.Id))
		}
	}
}

func printList(list *scheduleNode) {
	if list == nil {
		logs.Info(nil, "print list: schedule list is empty")
		return
	}

	logs.Info(nil, "schedule list:")
	for true {
		logs.Info(nil, fmt.Sprintf("schedule %d", list.data.Id))
		if list.next == nil {
			break
		}
		list = list.next
	}
}

// 处理所有到期的任务
func processSchedules() {
	for true {
		// printList(gManager.SchedulesList)
		// 执行之前做检查
		// 从列表中取出第一个，检查是否到时间
		node := getFirstNode()
		if node == nil {
			logs.Info(nil, "schedule list is empty")
			return
		}

		// 没到时间，退出
		sche := node.data
		if time.Now().Compare(sche.NextProTime) == -1 {
			logs.Info(nil, fmt.Sprintf("schedule %d time not ready, return.", sche.Id))
			return
		}

		// 执行定时任务
		doProc(sche)

		// 执行完成，更新任务列表
		afterProc(node)
	}
}

func Init() error {
	// 只能初始化一次
	if gManager == nil {
		gManager = &ScheduleMgr{
			SchedulesMap:  make(map[int64]*scheduleNode),
			SchedulesList: nil,
			StopRunning:   make(chan bool),
			lock:          sync.Mutex{},
		}
		logs.Info(nil, "schedule init")
	}
	return nil
}

func Run() {
	// 启动后，每隔一秒检查一次
	logs.Info(nil, "schedule run")
	for true {
		select {
		case <-gManager.StopRunning:
			logs.Info(nil, "schedule stop")
			return
		case <-time.After(time.Second * 5):
			// check schedule
			logs.Info(nil, "schedule check")
			processSchedules()
			// default:
			// 	time.Sleep(time.Second)
			// 	logs.Info(nil, "schedule sleep")
		}
	}
}

func Destroy() {
	if gManager != nil {
		gManager.StopRunning <- true
		close(gManager.StopRunning)
		gManager = nil
	}
	logs.Info(nil, "schedule destroy")
}

func AddSchedule(schedule *Schedule) (int64, error) {
	if schedule.Handler == nil {
		return 0, errors.New("handler is nil")
	}
	if schedule.Name == "" {
		return 0, errors.New("name is empty")
	}
	if schedule.FirstProTime.IsZero() {
		return 0, errors.New("first pro time is empty")
	}
	if schedule.Type == SCHEDULE_TYPE_INTERVAL && schedule.Interval <= 0 {
		return 0, errors.New("interval is invalid")
	}

	schedule.Id = time.Now().UnixMicro()
	schedule.NextProTime = schedule.FirstProTime
	err := addNode(schedule)
	if err != nil {
		logs.Info(nil, fmt.Sprintf("add schedule %v failed: %s", schedule, err.Error()))
	} else {
		logs.Info(nil, fmt.Sprintf("add schedule %v", schedule))
	}
	return schedule.Id, nil
}

func RemoveSchedule(id int64) error {
	logs.Info(nil, fmt.Sprintf("remove schedule %d", id))
	return removeNode(id)
}

// 任务状态机
func modify(id int64, status int) error {
	sche, err := getSchedule(id)
	if err != nil {
		return fmt.Errorf("schedule %d not exist", id)
	}

	switch status {
	case SCHEDULE_STATUS_INIT:
		err = errors.New("can not modify to init")
	case SCHEDULE_STATUS_RUN:
		if sche.data.Status == SCHEDULE_STATUS_RUN {
			return errors.New("already started")
		}
		sche.data.Status = SCHEDULE_STATUS_RUN
	case SCHEDULE_STATUS_STOP:
		if sche.data.Status == SCHEDULE_STATUS_RUN {
			sche.data.Status = SCHEDULE_STATUS_STOP
		} else {
			err = errors.New("schedule not started")
		}
	}
	return err
}

func StartSchedule(id int64) error {
	return modify(id, SCHEDULE_STATUS_RUN)
}

func StopSchedule(id int64) error {
	return modify(id, SCHEDULE_STATUS_STOP)
}

func GetSchedule(id int64) (*Schedule, error) {
	sche, err := getSchedule(id)
	if err == nil {
		return sche.data, nil
	} else {
		return nil, err
	}
}
