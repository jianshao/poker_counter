package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jianshao/poker_counter/prisma/db"
	"github.com/jianshao/poker_counter/src/model/user"
	"github.com/jianshao/poker_counter/src/utils"
)

type UserReq struct {
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	OpenId string `json:"open_id,omitempty"`
	Code   string `json:"code,omitempty"`
}

func userLoginCtrl(c *gin.Context) {
	// 从请求体中读取 code
	var params UserReq
	if err := c.BindJSON(&params); err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 1, "Invalid request body")
		return
	}

	// 将用户信息载入，即为活跃状态
	userInfo := user.UserLogin(params.Id)

	// 返回用户信息给客户端
	utils.BuildResponseOk(c, buildPlayerInfoResp(userInfo))
}

// 根据code获取openid，用于后续登录/注册
func getOpenIdCtrl(c *gin.Context) {
	// 从请求体中读取 code
	var params UserReq
	if err := c.BindJSON(&params); err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 1, "Invalid request body")
		return
	}

	openid, _, err := utils.GetWechatOpenidAndSessionKey(params.Code)
	if err == nil {
		utils.BuildResponseOk(c, buildPlayerInfoResp(&user.PlayerInfo{
			OpenId: openid,
		}))
	} else {
		// 微信接口出错
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
	}
}

func userCheckCtrl(c *gin.Context) {
	// 从请求体中读取 code
	var params UserReq
	if err := c.BindJSON(&params); err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 1, "Invalid request body")
		return
	}

	user, err := user.UserCheck(params.OpenId)
	if err == nil {
		utils.BuildResponseOk(c, buildPlayerInfoResp(user))
	} else {
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
	}
}

func userRegisterCtrl(c *gin.Context) {
	// 从请求体中读取 code
	var params UserReq
	if err := c.BindJSON(&params); err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 1, "Invalid request body")
		return
	}

	if params.OpenId == "" || params.Name == "" {
		utils.BuildResponse(c, http.StatusOK, nil, 1, "Invalid request body")
		return
	}

	user, err := user.UserRegister(params.Name, params.OpenId)
	if err == nil {
		utils.BuildResponseOk(c, buildPlayerInfoResp(user))
	} else {
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
	}
}

// 暂时没有更新
func userUpdateCtrl(c *gin.Context) {
	// 从请求体中读取 code
	var params UserReq
	if err := c.BindJSON(&params); err != nil || params.Name == "" {
		utils.BuildResponse(c, http.StatusOK, nil, 1, "Invalid request body")
		return
	}

	client := utils.GetPrismaClient()
	if client == nil {
		utils.BuildResponse(c, http.StatusOK, nil, 2, "failed to get prisma client")
		return
	}

	client.User.UpsertOne(
		db.User.ID.Equals(params.Id),
	).Update(
		db.User.Name.Set(params.Name),
	).Exec(context.Background())

	utils.BuildResponseOk(c, nil)
}
