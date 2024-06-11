package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jianshao/poker_counter/src/utils"
)

func buildRouters(r *gin.Engine) {

	// user
	r.POST(utils.BuildRouterPath("v1", "user/update"), userUpdateCtrl)
	r.POST(utils.BuildRouterPath("v1", "openid"), getOpenIdCtrl)
	r.POST(utils.BuildRouterPath("v1", "user/check"), userCheckCtrl)
	r.POST(utils.BuildRouterPath("v1", "user/register"), userRegisterCtrl)
	r.POST(utils.BuildRouterPath("v1", "user/login"), userLoginCtrl)

	// room
	r.POST(utils.BuildRouterPath("v1", "room/create"), createRoomCtrl)
	r.POST(utils.BuildRouterPath("v1", "room/user/entry"), entryRoomCtrl)
	r.POST(utils.BuildRouterPath("v1", "room/user/game/join"), joinGameCtrl)
	r.POST(utils.BuildRouterPath("v1", "room/user/game/quit"), quitGameCtrl)
	r.POST(utils.BuildRouterPath("v1", "room/user/leave"), leaveRoomCtrl)
	r.POST(utils.BuildRouterPath("v1", "room/close"), closeRoomCtrl)
	r.GET(utils.BuildRouterPath("v1", "room/info"), getRoomInfoCtrl)
	r.GET(utils.BuildRouterPath("v1", "room/check"), checkRoomCtrl)

	// records
	r.POST(utils.BuildRouterPath("v1", "room/score/apply"), applyBuyInCtrl)
	r.POST(utils.BuildRouterPath("v1", "room/score/confirm"), confirmBuyInCtrl)
	r.GET(utils.BuildRouterPath("v1", "room/score/all"), getApplyScoreAllCtrl)
}
