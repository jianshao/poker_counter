package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jianshao/poker_counter/src/model/room"
	"github.com/jianshao/poker_counter/src/utils"
)

type roomRequestParams struct {
	RoomId int `json:"room_id,omitempty"`
	UserId int `json:"user_id,omitempty"`
}

func buildRoomParams(c *gin.Context) (*roomRequestParams, error) {
	var request roomRequestParams
	if err := c.BindJSON(&request); err != nil {
		return nil, err
	}
	return &request, nil
}

// 1. owner create room
func createRoomCtrl(c *gin.Context) {
	params, err := buildRoomParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}

	if params.UserId == 0 {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 2, "用户信息不正确")
		return
	}

	roomInfo, err := room.CreateRoom(params.UserId)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 3, err.Error())
	} else {
		utils.BuildResponseOk(c, buildRoomInfoResp(roomInfo))
	}
}

func checkRoomCtrl(c *gin.Context) {
	roomIdStr := c.DefaultQuery("room_id", "")
	if roomId, err := strconv.Atoi(roomIdStr); err == nil {
		roomInfo := room.CheckRoom(roomId)
		if roomInfo != nil {
			utils.BuildResponseOk(c, buildRoomInfoResp(roomInfo))
		} else {
			utils.BuildResponse(c, http.StatusOK, nil, 2, "room not exist")
		}
	} else {
		utils.BuildResponse(c, http.StatusOK, nil, 1, "room id error")
	}
}

// owner close room
func closeRoomCtrl(c *gin.Context) {
	params, err := buildRoomParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}
	room.CloseRoom(params.RoomId, params.UserId)

	utils.BuildResponseOk(c, nil)
}

type boolResp struct {
	success bool `json:"success"`
}

func buildBoolResp(success bool) boolResp {
	return boolResp{
		success: success,
	}
}

// 用户进入房间，会保存用户所在的房间，下次登录时自动跳转到该房间
func entryRoomCtrl(c *gin.Context) {
	params, err := buildRoomParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}
	res, err := room.EntryRoom(params.RoomId, params.UserId)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
	} else {
		utils.BuildResponseOk(c, buildBoolResp(res))
	}
}

func leaveRoomCtrl(c *gin.Context) {
	params, err := buildRoomParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}
	res, err := room.LeaveRoom(params.RoomId, params.UserId)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
	} else {
		utils.BuildResponseOk(c, buildBoolResp(res))
	}
}

// user join game
func joinGameCtrl(c *gin.Context) {
	params, err := buildRoomParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}
	res, err := room.JoinGame(params.RoomId, params.UserId)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
	} else {
		utils.BuildResponseOk(c, buildBoolResp(res))
	}
}

// user leave game
func quitGameCtrl(c *gin.Context) {
	params, err := buildRoomParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}
	res, err := room.QuitGame(params.RoomId, params.UserId)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
	} else {
		utils.BuildResponseOk(c, buildBoolResp(res))
	}
}

func getRoomInfoCtrl(c *gin.Context) {
	roomIdStr := c.DefaultQuery("room_id", "")
	if roomId, err := strconv.Atoi(roomIdStr); err == nil {
		roomInfo, err := room.GetRoomInfo(roomId)
		if err == nil {
			utils.BuildResponseOk(c, buildRoomInfoResp(roomInfo))
		} else {
			utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
		}
	} else {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "room id error")
	}
}
