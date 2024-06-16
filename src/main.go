package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jianshao/poker_counter/src/controller"
	"github.com/jianshao/poker_counter/src/utils"
	"github.com/jianshao/poker_counter/src/utils/logs"
	"github.com/jianshao/poker_counter/src/utils/schedule"
	"github.com/joho/godotenv"
)

func handler() error {
	fmt.Println("handler")
	return nil
}

func handler2() error {
	fmt.Println("2")
	return nil
}

func AddSchedule(name string, handler func() error, firstProTime time.Time, interval int, scheduleType int) {
	id, err := schedule.AddSchedule(&schedule.Schedule{
		Name:         name,
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

func Init(router *gin.Engine) {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	file, _ := os.OpenFile("./logs/gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	router.Use(gin.LoggerWithWriter(file))

	controller.Init(router)
	logs.Init()
	go func() {
		schedule.Init()
		AddSchedule("handler1", handler, time.Now().Add(time.Second*10), 0, schedule.SCHEDULE_TYPE_FIXED)

		time.Sleep(time.Second)
		AddSchedule("handler2", handler2, time.Now().Add(time.Second*10), 10, schedule.SCHEDULE_TYPE_INTERVAL)

		schedule.Run()
	}()
}

func close(router *gin.Engine) {
	utils.Close()
	schedule.Destroy()
}

func main() {
	router := gin.Default()
	env := os.Getenv("ENVIRONMENT")
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	Init(router)

	// 创建一个通道来接收信号
	// 监听中断信号，例如在 Unix 系统中的 SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT, syscall.SIGKILL, syscall.SIGSTOP)

	go func(router *gin.Engine) {
		// 阻塞直到接收到信号
		<-c
		// 执行退出前的清理工作
		fmt.Println("收到退出信号，正在退出...")
		close(router)

		// 退出程序
		os.Exit(1)
	}(router)

	router.Run(":8989") // 监听并在 0.0.0.0:8080 上启动服务
}
