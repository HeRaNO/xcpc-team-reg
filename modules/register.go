package modules

import (
	"io/ioutil"
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/model"
	"github.com/HeRaNO/xcpc-team-reg/util"
	jsoniter "github.com/json-iterator/go"
)

func Register(w http.ResponseWriter, r *http.Request) {
	// validate student id's length
	// validate email token
	// insert into RDB
	// go to login
	defer r.Body.Close()
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
	}

	usrinfo := model.UserRegister{}
	err = jsoniter.Unmarshal(bd, &usrinfo)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
	}

	uid, _ := model.GetUserIDByEmail(r.Context(), &usrinfo.Email)
	if uid != -1 {
		util.ErrorResponse(w, r, "email has already registered", config.ERR_WRONGINFO)
		return
	}

	if !config.StuIDMap[len(usrinfo.StuID)] { // Just a naive check
		util.ErrorResponse(w, r, "student id's length is wrong", config.ERR_WRONGINFO)
		return
	}

	err = model.ValidateEmailToken(r.Context(), &usrinfo.Email, &usrinfo.EmailToken, &usrinfo.Action)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_WRONGINFO)
		return
	}

	err = model.CreateNewUser(r.Context(), usrinfo, 0)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	util.SuccessResponse(w, r, "ok")
}

func ForgotPwd(w http.ResponseWriter, r *http.Request) {
	// ...
}
