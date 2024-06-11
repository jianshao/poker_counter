package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BuildRouterPath(ver, path string) string {
	return "api/" + ver + "/" + path
}

func BuildResponseOk(c *gin.Context, data interface{}) {
	BuildResponse(c, http.StatusOK, data, 0, "success")
}

func BuildResponse(c *gin.Context, status int, data interface{}, code int, message string) {
	c.JSON(status, ApiResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func GetCurrTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Init() {

}

func Close() {
	closePrisma()
	closeRedis()
}
