package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jianshao/poker_counter/src/model/records"
	"github.com/jianshao/poker_counter/src/model/room"
	"github.com/jianshao/poker_counter/src/model/user"
)

func Init(r *gin.Engine) {
	buildRouters(r)
}

type ApplyScoreResp struct {
	Id          int    `json:"id"`
	PlayerId    int    `json:"player_id"`
	RoomId      int    `json:"room_id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	Status      int    `json:"status"`
	ApplyType   int    `json:"apply_type"`
	ApplyTime   string `json:"apply_time"`
	ConfirmTime string `json:"confirm_time"`
}

type ApplyScoreListResp struct {
	ApplyList []ApplyScoreResp `json:"applies"`
	Count     int              `json:"count"`
}

type PlayerInfoResp struct {
	Id         int              `json:"id"`
	Name       string           `json:"name"`
	OpenId     string           `json:"open_id"`
	Status     int              `json:"status"`
	CurrRoomId int              `json:"curr_room_id"`
	CurrScore  int              `json:"curr_score"`
	FinalScore int              `json:"final_score"`
	JoinTime   string           `json:"join_time"`
	ExitTime   string           `json:"exit_time"`
	Applies    []ApplyScoreResp `json:"applies"`
}

type RoomInfoResp struct {
	Id        int              `json:"room_id"`
	Owner     int              `json:"owner"`
	Status    int              `json:"status"`
	StartTime string           `json:"start_time"`
	Players   []PlayerInfoResp `json:"players"`
}

func buildApplyScoreResp(applyScore *records.ApplyScore) ApplyScoreResp {
	return ApplyScoreResp{
		Id:          applyScore.Id,
		PlayerId:    applyScore.UserId,
		Name:        applyScore.Name,
		RoomId:      applyScore.RoomId,
		Score:       applyScore.Score,
		Status:      applyScore.Status,
		ApplyType:   applyScore.ApplyType,
		ApplyTime:   applyScore.ApplyTime,
		ConfirmTime: applyScore.ConfirmTime,
	}
}

func buildPlayerInfoResp(userInfo *user.PlayerInfo, roomId int) PlayerInfoResp {
	// 只需要用户本身的静态数据
	if roomId == 0 {
		return PlayerInfoResp{
			Id:         userInfo.Id,
			Name:       userInfo.Name,
			OpenId:     userInfo.OpenId,
			CurrRoomId: userInfo.CurrRoomId,
		}
	}
	room, ok := userInfo.Rooms[roomId]
	if !ok {
		room = &user.UserRoomInfo{
			ApplyList: map[int]int{},
		}
	}
	applies := []ApplyScoreResp{}
	for _, applyId := range room.ApplyList {
		apply, _ := records.GetApply(applyId)
		applies = append(applies, buildApplyScoreResp(apply))
	}
	return PlayerInfoResp{
		Id:         userInfo.Id,
		Name:       userInfo.Name,
		OpenId:     userInfo.OpenId,
		CurrRoomId: userInfo.CurrRoomId,
		Status:     room.Status,
		CurrScore:  room.CurrScore,
		FinalScore: room.FinalScore,
		JoinTime:   room.JoinTime,
		ExitTime:   room.ExitTime,
		Applies:    applies,
	}
}

func buildRoomInfoResp(roomInfo *room.RoomInfo) RoomInfoResp {
	players := []PlayerInfoResp{}
	for _, playerId := range roomInfo.Players {
		players = append(players, buildPlayerInfoResp(user.GetUser(playerId), roomInfo.RoomId))
	}
	return RoomInfoResp{
		Id:      roomInfo.RoomId,
		Owner:   roomInfo.Owner,
		Status:  roomInfo.Status,
		Players: players,
	}
}

func buildApplyListResp(applyList []records.ApplyScore) ApplyScoreListResp {
	applyListResp := []ApplyScoreResp{}
	for _, apply := range applyList {
		applyListResp = append(applyListResp, buildApplyScoreResp(&apply))
	}
	return ApplyScoreListResp{
		ApplyList: applyListResp,
		Count:     len(applyListResp),
	}
}
