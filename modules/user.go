package modules

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/model"
	"github.com/HeRaNO/xcpc-team-reg/util"
	jsoniter "github.com/json-iterator/go"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Fetch from RDB
	// return info

	uid := getUserIDFromReq(r)
	if uid <= 0 {
		util.ErrorResponse(w, r, "user status error", config.ERR_UNAUTHORIZED)
		return
	}

	info, err := model.GetUserInfoByID(r.Context(), uid)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}
	if info == nil {
		util.ErrorResponse(w, r, "no user record but why???", config.ERR_INTERNAL)
		return
	}

	util.SuccessResponse(w, r, *info)
}

func ModifyUserInfo(w http.ResponseWriter, r *http.Request) {
	// validate info
	// write to RDB

	defer r.Body.Close()
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	usrinfo := model.UserInfoModify{}
	usrinfo.Name = template.HTMLEscapeString(usrinfo.Name)
	usrinfo.StuID = template.HTMLEscapeString(usrinfo.StuID)
	err = jsoniter.Unmarshal(bd, &usrinfo)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	if !config.StuIDMap[len(usrinfo.StuID)] { // Just a naive check
		util.ErrorResponse(w, r, "student id's length is wrong", config.ERR_WRONGINFO)
		return
	}

	uid := getUserIDFromReq(r)
	if uid <= 0 {
		util.ErrorResponse(w, r, "user status error", config.ERR_UNAUTHORIZED)
		return
	}

	err = model.ModifyUserInfoByID(r.Context(), uid, &usrinfo)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	util.SuccessResponse(w, r, "ok")
}
