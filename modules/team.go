package modules

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/model"
	"github.com/HeRaNO/xcpc-team-reg/util"
	jsoniter "github.com/json-iterator/go"
)

func GetTeamInfo(w http.ResponseWriter, r *http.Request) {
	// Fetch from RDB
	// return info
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
		util.ErrorResponse(w, r, "user hasn't joined a team", config.ERR_WRONGINFO)
		return
	}

	info, err := model.GetTeamInfoByTeamID(r.Context(), tid)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	info.TeamName = template.HTMLEscapeString(info.TeamName)

	util.SuccessResponse(w, r, *info)
}

func ModifyTeamInfo(w http.ResponseWriter, r *http.Request) {
	// validate info
	// write to RDB

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
		util.ErrorResponse(w, r, "user hasn't joined a team", config.ERR_WRONGINFO)
		return
	}

	defer r.Body.Close()
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	teamInfo := model.TeamInfoModify{}
	teamInfo.TeamName = template.HTMLEscapeString(teamInfo.TeamName)
	err = jsoniter.Unmarshal(bd, &teamInfo)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	teamName := teamInfo.TeamName
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

	err = model.ModifyTeamInfoByTeamID(r.Context(), tid, &teamName)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	util.SuccessResponse(w, r, "ok")
}
