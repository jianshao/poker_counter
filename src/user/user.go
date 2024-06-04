package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"private/backend/gamesRoom/prisma/db"
	"private/backend/gamesRoom/src/config"
	"private/backend/gamesRoom/src/utils"

	"github.com/gin-gonic/gin"
)

func register(r *gin.Engine) {
	r.POST(utils.BuildRouterPath("v1", "user/login"), userLogin)

	r.POST(utils.BuildRouterPath("v1", "user/update"), userUpdate)
}

func Init(router *gin.Engine) {
	register(router)
}

func userLogin(c *gin.Context) {
	// 从请求体中读取 code
	var params UserReq
	if err := c.BindJSON(&params); err != nil || params.Code == "" {
		utils.BuildResponse(c, http.StatusOK, nil, 1, "Invalid request body")
		return
	}

	// 调用微信接口获取 openid 和 session_key
	openid, sessionKey, err := getWechatOpenidAndSessionKey(params.Code)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 2, err.Error())
		return
	}

	// 根据 openid 获取用户信息（这里需要调用你的业务逻辑代码）
	userInfo, isNewUser, err := getUserInfoByOpenid(openid)
	if err != nil {
		utils.BuildResponse(c, http.StatusOK, nil, 3, err.Error())
		return
	}

	// 返回用户信息给客户端
	utils.BuildResponseOk(c, UserResp{
		UserBase: UserBase{
			Id:         userInfo.ID,
			Name:       userInfo.Name,
			Avatar:     userInfo.Avatar,
			OpenId:     openid,
			SessionKey: sessionKey,
		},
		IsNewUser: isNewUser,
		CreatedAt: userInfo.CreatedTime.String(),
		UpdatedAt: userInfo.UpdatedTime.String(),
	})
}

func getWechatOpenidAndSessionKey(code string) (openid, sessionKey string, err error) {

	base_url := "https://api.weixin.qq.com/sns/jscode2session?grant_type=authorization_code"
	url := fmt.Sprintf("%s&appid=%s&secret=%s&js_code=%s", base_url, config.APP_ID, config.APP_SECRET, code)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result struct {
		Openid     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", "", err
	}

	if result.ErrCode != 0 {
		return "", "", fmt.Errorf("error from wechat: %d - %s", result.ErrCode, result.ErrMsg)
	}

	return result.Openid, result.SessionKey, nil
}

func getUserInfoByOpenid(openid string) (*db.UserModel, bool, error) {
	// 这里应该是调用你的业务逻辑代码，例如查询数据库
	// 以下为示例数据
	client := utils.GetPrismaClient()
	if client == nil {
		return nil, false, fmt.Errorf("failed to get prisma client")
	}

	// 如果不存在就创建一个
	isNewUser := false
	user, err := client.User.FindUnique(
		db.User.Openid.Equals(openid),
	).Exec(context.Background())
	if err != nil {
		if err == db.ErrNotFound {
			// user, err = client.User.CreateOne(
			// 	db.User.Name.Set(""),
			// 	db.User.Avatar.Set(""),
			// 	db.User.Openid.Set(openid),
			// ).Exec(context.Background())
			user = &db.UserModel{}
			user.Openid = openid
			isNewUser = true
		}
	}
	isNewUser = true

	return user, isNewUser, err
}

func userUpdate(c *gin.Context) {
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
