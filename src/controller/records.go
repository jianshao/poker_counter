package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jianshao/poker_counter/src/model/room"
	"github.com/jianshao/poker_counter/src/utils"
)

type RecordsReq struct {
	RoomId    int `json:"room_id,omitempty"`
	Owner     int `json:"owner,omitempty"`
	UserId    int `json:"user_id,omitempty"`
	ApplyId   int `json:"apply_id,omitempty"`
	Score     int `json:"score,omitempty"`
	ApplyType int `json:"apply_type,omitempty"`
	Status    int `json:"status,omitempty"`
}

func buildRecordParams(c *gin.Context) (*RecordsReq, error) {
	var params RecordsReq
	if err := c.BindJSON(&params); err != nil {
		return nil, err
	}
	return &params, nil
}

// user apply buy in
func applyBuyInCtrl(c *gin.Context) {
	params, err := buildRecordParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}

	apply, err := room.ApplyBuyIn(params.RoomId, params.UserId, params.Score, params.ApplyType)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
	} else {
		utils.BuildResponseOk(c, buildApplyScoreResp(apply))
	}
}

// owner accept buy in
func confirmBuyInCtrl(c *gin.Context) {
	params, err := buildRecordParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}

	apply, err := room.ConfirmBuyIn(params.RoomId, params.Owner, params.ApplyId, params.Status)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
	} else {
		utils.BuildResponseOk(c, buildApplyScoreResp(apply))
	}

}

func getApplyScoreAllCtrl(c *gin.Context) {
	roomIdStr := c.DefaultQuery("room_id", "")
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 1, err.Error())
	}

	applies, err := room.GetAllScoreApplies(roomId)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
	} else {
		utils.BuildResponseOk(c, buildApplyListResp(applies))
	}
}
