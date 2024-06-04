package main

import (
	"log"
	"os"
	"private/backend/gamesRoom/src/room"
	"private/backend/gamesRoom/src/user"
	"private/backend/gamesRoom/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	file, _ := os.OpenFile("logs/gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	router.Use(gin.LoggerWithWriter(file))

	utils.Init()
	room.Init(router)
	user.Init(router)

	router.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func main() {
	Init()
}
