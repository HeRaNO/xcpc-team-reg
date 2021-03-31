package modules

import (
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/util"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Fetch from RDB
	// return info

	uid := getUserIDFromReq(r)
	util.SuccessResponse(w, r, map[string]int64{"uid": uid})
}

func ModifyUserInfo(w http.ResponseWriter, r *http.Request) {
	// validate info
	// write to RDB
}
