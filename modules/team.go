package modules

import (
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/model"
	"github.com/HeRaNO/xcpc-team-reg/util"
)

func GetTeamInfo(w http.ResponseWriter, r *http.Request) {
	// Fetch from RDB
	// return info
	uid := getUserIDFromReq(r)
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

	util.SuccessResponse(w, r, *info)
}

func ModifyTeamInfo(w http.ResponseWriter, r *http.Request) {
	// validate info
	// write to RDB
}
