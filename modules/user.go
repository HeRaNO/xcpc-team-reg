package modules

import (
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/util"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Fetch from RDB
	// return info
	util.SuccessResponse(w, r, "hello")
}

func ModifyUserInfo(w http.ResponseWriter, r *http.Request) {
	// validate info
	// write to RDB
}
