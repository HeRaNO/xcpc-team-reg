package modules

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/model"
	"github.com/HeRaNO/xcpc-team-reg/util"
	jsoniter "github.com/json-iterator/go"
)

// User create a team
func CreateTeam(w http.ResponseWriter, r *http.Request) {
	// Fetch user team status
	// If user join in a team -> failed, must quit the original team
	// insert record to RDB
	// user.team_id <- team_id

	uid := getUserIDFromReq(r)
	if uid <= 0 {
		util.ErrorResponse(w, r, "user status error", config.ERR_UNAUTHORIZED)
		return
	}

	tid, err := model.GetTeamIDByUserID(r.Context(), uid)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}
	if tid != 0 {
		util.ErrorResponse(w, r, "user has already joined in a team", config.ERR_WRONGINFO)
		return
	}

	defer r.Body.Close()
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	teamInfo := model.TeamInfoModify{}
	err = jsoniter.Unmarshal(bd, &teamInfo)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	teamName := template.HTMLEscapeString(teamInfo.TeamName)
	nameLength := utf8.RuneCountInString(teamName)
	if nameLength < 1 {
		util.ErrorResponse(w, r, "team name cannot be empty", config.ERR_WRONGINFO)
		return
	}
	if nameLength > config.MaxTeamNameLength {
		util.ErrorResponse(w, r, "team name too long", config.ERR_WRONGINFO)
		return
	}

	used, err := model.IsTeamNameUsed(r.Context(), &teamName)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}
	if used {
		util.ErrorResponse(w, r, "team name is used", config.ERR_WRONGINFO)
		return
	}

	tid, inviteToken, err := model.CreateNewTeam(r.Context(), &teamName, uid)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	teamInfoResp := model.JoinTeamRequest{
		TeamID:      fmt.Sprintf("%d", tid),
		InviteToken: inviteToken,
	}

	util.SuccessResponse(w, r, teamInfoResp)
}

// User join in a team
func JoinTeam(w http.ResponseWriter, r *http.Request) {
	// Fetch user team status
	// if user.team_id != 0 -> failed, has joined a team
	// Fetch team status from RDB
	// failed -> failed, not exist
	// If mem_cnt >= 3 -> failed, cannot join
	// user.team_id <- team_id, team.mem_cnt++

	uid := getUserIDFromReq(r)
	if uid <= 0 {
		util.ErrorResponse(w, r, "user status error", config.ERR_UNAUTHORIZED)
		return
	}

	tid, err := model.GetTeamIDByUserID(r.Context(), uid)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}
	if tid != 0 {
		util.ErrorResponse(w, r, "user has already joined in a team", config.ERR_WRONGINFO)
		return
	}

	defer r.Body.Close()
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	teamInfo := model.JoinTeamRequest{}
	err = jsoniter.Unmarshal(bd, &teamInfo)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	joinTeamID, err := strconv.ParseInt(teamInfo.TeamID, 10, 64)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_WRONGINFO)
		return
	}
	joinTeamToken := teamInfo.InviteToken

	isValidToken, err := model.ValidateTeamInviteToken(r.Context(), joinTeamID, &joinTeamToken)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}
	if !isValidToken {
		util.ErrorResponse(w, r, "invalid invite token or team id", config.ERR_WRONGINFO)
		return
	}

	joined, err := model.UserJoinTeam(r.Context(), uid, joinTeamID)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}
	if !joined {
		util.ErrorResponse(w, r, "cannot join in the team", config.ERR_WRONGINFO)
		return
	}

	util.SuccessResponse(w, r, "ok")
}

// User quit a team
func QuitTeam(w http.ResponseWriter, r *http.Request) {
	// Fetch user team status
	// if user.team_id = 0 -> failed, hasn't join a team
	// Fetch team status from RDB
	// failed -> failed, not exist
	// user.team_id <- 0
	// team.mem_cnt = 1 -> delete the team record
	// else team.mem_cnt--

	uid := getUserIDFromReq(r)
	if uid <= 0 {
		util.ErrorResponse(w, r, "user status error", config.ERR_UNAUTHORIZED)
		return
	}

	tid, err := model.GetTeamIDByUserID(r.Context(), uid)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}
	if tid == 0 {
		util.ErrorResponse(w, r, "user hasn't joined in a team", config.ERR_WRONGINFO)
		return
	}

	err = model.UserQuitTeam(r.Context(), uid, tid)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	util.SuccessResponse(w, r, "ok")
}
