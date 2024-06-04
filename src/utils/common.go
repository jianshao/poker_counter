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

type RequestParams struct {
	RoomId    int `json:"room_id,omitempty"`
	UserId    int `json:"user_id,omitempty"`
	Score     int `json:"score,omitempty"`
	Status    int `json:"status,omitempty"`
	ApplyId   int `json:"apply_id,omitempty"`
	ApplyType int `json:"apply_type,omitempty"`
}

type RoomInfoResp struct {
	RoomId    int              `json:"room_id"`
	Owner     int              `json:"owner"`
	Status    int              `json:"status"`
	StartTime string           `json:"start_time"`
	Players   []PlayerInfoResp `json:"players"`
	// RoomName  string `json:"room_name"`
	// RoomType  string `json:"room_type"`
}

type PlayerInfoResp struct {
	PlayerId   int    `json:"player_id"`
	Status     int    `json:"status"`
	PlayerName string `json:"player_name"`
	// PlayerType string `json:"player_type"`
	CurrScore  int    `json:"curr_score"`
	FinalScore int    `json:"final_score"`
	JoinTime   string `json:"join_time"`
	ExitTime   string `json:"exit_time"`
}

type ApplyScoreListResp struct {
	ApplyList []ApplyScoreResp `json:"applies"`
	Count     int              `json:"count"`
}

type ApplyScoreResp struct {
	ApplyId     int    `json:"apply_id"`
	PlayerId    int    `json:"player_id"`
	RoomId      int    `json:"room_id"`
	Score       int    `json:"score"`
	Status      int    `json:"status"`
	ApplyType   int    `json:"apply_type"`
	ApplyTime   string `json:"apply_time"`
	ConfirmTime string `json:"confirm_time"`
}

type HandleFunc func(c *gin.Context)

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
