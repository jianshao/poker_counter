package room

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jianshao/poker_counter/src/utils"
)

type ApplyScore struct {
	applyId    int
	userId     int
	score      int
	status     int
	applyType  int
	applyTime  string
	acceptTime string
}

type PlayerInfo struct {
	id         int
	name       string
	status     int
	currScore  int
	finalScore int
	joinTime   string
	exitTime   string
	applyList  []ApplyScore
}

type RoomInfo struct {
	roomId  int
	owner   int
	status  int
	players map[int]*PlayerInfo
}

var (
	roomMap = map[int]*RoomInfo{}
	roomId  = 1
	applyId = 1
)

func register(r *gin.Engine) {
	r.POST(utils.BuildRouterPath("v1", "room/create"), createRoom)
	r.POST(utils.BuildRouterPath("v1", "room/join"), joinRoom)
	r.POST(utils.BuildRouterPath("v1", "room/close"), closeRoom)
	r.POST(utils.BuildRouterPath("v1", "room/score/apply"), applyBuyIn)
	r.POST(utils.BuildRouterPath("v1", "room/score/confirm"), confirmBuyIn)
	r.POST(utils.BuildRouterPath("v1", "room/leave"), leaveRoom)
	r.GET(utils.BuildRouterPath("v1", "room/info"), getRoomInfo)
	r.GET(utils.BuildRouterPath("v1", "room/score/list"), getApplyScoreList)
	r.GET(utils.BuildRouterPath("v1", "room/score/all"), getApplyScoreAll)
	r.GET(utils.BuildRouterPath("v1", "room/check"), checkRoom)
}

func Init(r *gin.Engine) {
	register(r)
}

func buildParams(c *gin.Context) (*RequestParams, error) {
	var params RequestParams
	if err := c.BindJSON(&params); err != nil {
		return nil, err
	}
	return &params, nil
}

func buildPlayerInfoResp(userInfo *PlayerInfo) PlayerInfoResp {
	return PlayerInfoResp{
		PlayerId:   userInfo.id,
		PlayerName: userInfo.name,
		Status:     userInfo.status,
		CurrScore:  userInfo.currScore,
		FinalScore: userInfo.finalScore,
		JoinTime:   userInfo.joinTime,
		ExitTime:   userInfo.exitTime,
	}
}

func buildPlayerInfoRespList(players map[int]*PlayerInfo) []PlayerInfoResp {
	playersResp := []PlayerInfoResp{}
	for _, player := range players {
		playersResp = append(playersResp, buildPlayerInfoResp(player))
	}
	return playersResp
}

func buildRoomInfoResp(roomInfo *RoomInfo) RoomInfoResp {
	return RoomInfoResp{
		RoomId:  roomInfo.roomId,
		Owner:   roomInfo.owner,
		Status:  roomInfo.status,
		Players: buildPlayerInfoRespList(roomInfo.players),
	}
}

func buildApplyScoreResp(applyScore *ApplyScore, roomId int) ApplyScoreResp {
	return ApplyScoreResp{
		ApplyId:     applyScore.applyId,
		PlayerId:    applyScore.userId,
		RoomId:      roomId,
		Score:       applyScore.score,
		Status:      applyScore.status,
		ApplyType:   applyScore.applyType,
		ApplyTime:   applyScore.applyTime,
		ConfirmTime: applyScore.acceptTime,
	}
}

func buildApplyRespList(applyList []ApplyScore, roomId int) []ApplyScoreResp {
	var applyListResp []ApplyScoreResp
	for _, apply := range applyList {
		applyListResp = append(applyListResp, buildApplyScoreResp(&apply, roomId))
	}
	return applyListResp
}

func generateRoomId() int {
	return 3918
}

func generateApplyId() int {
	applyId += 1
	return applyId
}

func checkRoom(c *gin.Context) {
	roomIdStr := c.DefaultQuery("room_id", "")
	if roomId, err := strconv.Atoi(roomIdStr); err == nil {
		if roomInfo := getActiveRoom(roomId); roomInfo != nil {
			utils.BuildResponseOk(c, nil)
		} else {
			utils.BuildResponse(c, http.StatusOK, nil, 1, "room not exists")
		}
	} else {
		utils.BuildResponse(c, http.StatusOK, nil, 2, "room id error")
	}
}

// 1. owner create room
func createRoom(c *gin.Context) {
	params, err := buildParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}

	roomId := generateRoomId()
	defaultRoomInfo := &RoomInfo{
		roomId:  roomId,
		owner:   params.UserId,
		status:  1,
		players: make(map[int]*PlayerInfo),
	}
	// 先查一下，如果房间已经存在并状态是正在使用中，则不能创建
	if roomInfo := getActiveRoom(roomId); roomInfo != nil && roomInfo.status == 1 {
		utils.BuildResponse(c, http.StatusOK, nil, 1, "room id already exists")
	} else {
		roomMap[roomId] = defaultRoomInfo
		utils.BuildResponseOk(c, buildRoomInfoResp(defaultRoomInfo))
	}
}

// owner close room
func closeRoom(c *gin.Context) {
	params, err := buildParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}

	if roomInfo, ok := roomMap[params.RoomId]; ok {
		if roomInfo.owner != params.UserId {
			utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "only owner can close room")
			return
		}
		roomInfo.status = 0
	}
	utils.BuildResponseOk(c, nil)
}

func getActiveRoom(roomId int) *RoomInfo {
	if roomInfo, ok := roomMap[roomId]; ok && roomInfo.status == 1 {
		return roomInfo
	}
	return nil
}

func getUserFromRoom(userId int, room *RoomInfo) *PlayerInfo {
	if player, ok := room.players[userId]; ok {
		return player
	}
	return nil
}

func getActiveUserFromRoom(userId int, room *RoomInfo) *PlayerInfo {
	if player := getUserFromRoom(userId, room); player != nil && player.status == 1 {
		return player
	}
	return nil
}

// user join room
func joinRoom(c *gin.Context) {
	params, err := buildParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}

	roomId := params.RoomId
	userId := params.UserId
	roomInfo := getActiveRoom(roomId)
	if roomInfo == nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "room not exists")
		return
	}

	if userInfo := getActiveUserFromRoom(userId, roomInfo); userInfo != nil {
		if userInfo.status != 1 {
			userInfo.status = 1
		}
	} else {
		roomInfo.players[userId] = &PlayerInfo{
			id:         userId,
			name:       fmt.Sprintf("玩家%d", userId),
			currScore:  0,
			finalScore: 0,
			status:     1,
			joinTime:   utils.GetCurrTime(),
			exitTime:   "",
			applyList:  []ApplyScore{},
		}
	}
	utils.BuildResponseOk(c, nil)
}

// user apply buy in
func applyBuyIn(c *gin.Context) {
	params, err := buildParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}

	if roomInfo := getActiveRoom(params.RoomId); roomInfo != nil {
		if userInfo := getActiveUserFromRoom(params.UserId, roomInfo); userInfo != nil {
			newApply := ApplyScore{
				applyId:   generateApplyId(),
				userId:    params.UserId,
				score:     params.Score,
				status:    0,
				applyTime: utils.GetCurrTime(),
				applyType: params.ApplyType,
			}
			userInfo.applyList = append(userInfo.applyList, newApply)
			utils.BuildResponseOk(c, newApply)
		} else {
			utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "user is not playing")
		}
	} else {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "room not exists")
	}
}

// owner accept buy in
func confirmBuyIn(c *gin.Context) {
	params, err := buildParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}

	if roomInfo := getActiveRoom(params.RoomId); roomInfo != nil {
		if roomInfo.owner != params.UserId {
			utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "only owner can accept buy in")
			return
		}

		if userInfo := getActiveUserFromRoom(params.UserId, roomInfo); userInfo != nil {
			for i, apply := range userInfo.applyList {
				if apply.applyId == params.ApplyId && apply.status == 0 {
					// 申请类型：1-申请买入，2-申请卸码
					if apply.applyType == 1 {
						userInfo.currScore += apply.score
					} else if apply.applyType == 2 {
						userInfo.finalScore += apply.score
						// 同时将玩家状态修改为"退出"
						userInfo.status = 0
						userInfo.exitTime = utils.GetCurrTime()
					}
					userInfo.applyList[i].status = params.Status
					userInfo.applyList[i].acceptTime = utils.GetCurrTime()
					utils.BuildResponseOk(c, userInfo.applyList[i])
					return
				}
			}
			utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "apply not exists")
		} else {
			utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "user is not playing")
		}
	} else {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "room not exists")
	}
}

// user leave room
func leaveRoom(c *gin.Context) {
	params, err := buildParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}

	if roomInfo := getActiveRoom(params.RoomId); roomInfo != nil {
		if userInfo := getUserFromRoom(params.UserId, roomInfo); userInfo != nil {
			userInfo.status = 0
			userInfo.exitTime = utils.GetCurrTime()
		}
		utils.BuildResponseOk(c, nil)
	} else {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "room not exists")
	}
}

func getRoomInfo(c *gin.Context) {
	roomIdStr := c.DefaultQuery("room_id", "")
	if roomId, err := strconv.Atoi(roomIdStr); err == nil {
		if roomInfo := getActiveRoom(roomId); roomInfo != nil {
			utils.BuildResponseOk(c, buildRoomInfoResp(roomInfo))
		} else {
			utils.BuildResponse(c, http.StatusOK, nil, 1, "room not exists")
		}
	} else {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "room id error")
	}
}

func getApplyScoreList(c *gin.Context) {
	params, err := buildParams(c)
	if err != nil {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, err.Error())
		return
	}

	if roomInfo := getActiveRoom(params.RoomId); roomInfo != nil {
		if userInfo := getUserFromRoom(params.UserId, roomInfo); userInfo != nil {
			applyList := buildApplyRespList(userInfo.applyList, params.RoomId)
			utils.BuildResponseOk(c, ApplyScoreListResp{ApplyList: applyList, Count: len(applyList)})
		} else {
			utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "user not exists")
		}
	} else {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "room not exists")
	}
}

func getApplyScoreAll(c *gin.Context) {
	roomIdStr := c.DefaultQuery("room_id", "")
	if roomId, err := strconv.Atoi(roomIdStr); err == nil {
		if roomInfo := getActiveRoom(roomId); roomInfo != nil {
			applyList := []ApplyScoreResp{}
			for _, player := range roomInfo.players {
				if player.status != 1 {
					continue
				}

				applyList = append(applyList, buildApplyRespList(player.applyList, roomId)...)
			}
			result := []ApplyScoreResp{}
			for _, apply := range applyList {
				if apply.Status == 0 {
					result = append(result, apply)
				}
			}
			utils.BuildResponseOk(c, ApplyScoreListResp{ApplyList: result, Count: len(result)})
		} else {
			utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "room not exists")
		}
	} else {
		utils.BuildResponse(c, http.StatusBadRequest, nil, 1, "room id error")
	}

}
