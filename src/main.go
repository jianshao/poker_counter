package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jianshao/poker_counter/src/room"
	"github.com/jianshao/poker_counter/src/user"
	"github.com/jianshao/poker_counter/src/utils"
	"github.com/joho/godotenv"
)

func Init() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router := gin.Default()
	env := os.Getenv("ENVIRONMENT")
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	file, _ := os.OpenFile("./logs/gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	router.Use(gin.LoggerWithWriter(file))

	utils.Init()
	room.Init(router)
	user.Init(router)

	router.Run(":8989") // 监听并在 0.0.0.0:8080 上启动服务
}

func main() {
	Init()
}
