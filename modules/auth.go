package modules

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/model"
	"github.com/HeRaNO/xcpc-team-reg/util"
	jsoniter "github.com/json-iterator/go"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func generateJWTToken(nowTime time.Time, uid int64, isAdmin bool, userFingerPrint *string) (string, error) {
	t := jwt.New()
	t.Set(jwt.IssuedAtKey, nowTime.Unix())
	t.Set(jwt.ExpirationKey, nowTime.Add(config.LOGIN_EXPIRETIME).Unix())
	uidToken := strconv.FormatInt(uid, 10)
	t.Set(config.JWTIDName, uidToken)
	admin := "0"
	if isAdmin {
		admin = "1"
	}
	t.Set(config.JWTAdminName, admin)
	fgp := util.SHA256([]byte(*userFingerPrint))
	t.Set(config.JWTFingerPrintName, fgp)
	tokenByte, err := jwt.Sign(t, jwa.HS256, config.JWTSecret)
	if err != nil {
		return "", err
	}
	return string(tokenByte), nil
}

// Must insure it can get user_id from request
func getUserIDFromReq(r *http.Request) int64 {
	userInfo := r.Context().Value(config.CtxUserInfoName).(*config.CtxUserInfo)
	return userInfo.ID
}

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
	}

	info := model.UserLogin{}
	err = jsoniter.Unmarshal(bd, &info)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
	}

	uid, err := model.GetUserIDByEmail(r.Context(), &info.Email)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}
	if uid == -1 {
		util.ErrorResponse(w, r, "no such user", config.ERR_WRONGINFO)
		return
	}

	err = model.ValidateAuthInfo(r.Context(), uid, &info.Email, &info.PwdToken)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_WRONGINFO)
		return
	}

	isAdmin, err := model.GetAdminByUserID(r.Context(), uid)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	fingerPrintToken, err := util.GenToken(config.FingerPrintTokenLength)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	nowTime := time.Now()
	secretToken := fmt.Sprintf("%d_%d_%s", nowTime.Unix(), uid, fingerPrintToken)

	jwtToken, err := generateJWTToken(nowTime, uid, isAdmin, &secretToken)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	fgpCookie := http.Cookie{
		Name:     config.JWTFingerPrintName,
		Value:    secretToken,
		Path:     "/",
		Domain:   config.Domain,
		Expires:  nowTime.Add(config.LOGIN_EXPIRETIME),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &fgpCookie)
	util.SuccessResponse(w, r, jwtToken)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	uid := getUserIDFromReq(r)
	if uid <= 0 {
		util.ErrorResponse(w, r, "user status error", config.ERR_UNAUTHORIZED)
		return
	}

	fgpExpireCookie := http.Cookie{
		Name:     config.JWTFingerPrintName,
		Path:     "/",
		Domain:   config.Domain,
		Expires:  time.Now(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &fgpExpireCookie)
	util.SuccessResponse(w, r, "ok")
}
