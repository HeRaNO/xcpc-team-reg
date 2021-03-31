package modules

import (
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

func generateJWTToken(uid int64, jwtSecret *string) (string, error) {
	t := jwt.New()
	t.Set(jwt.IssuedAtKey, time.Now().Unix())
	t.Set(jwt.ExpirationKey, time.Now().Add(config.LOGIN_EXPIRETIME).Unix())
	idToken := strconv.FormatInt(uid, 10)
	t.Set("id", idToken)
	tokenByte, err := jwt.Sign(t, jwa.HS256, []byte(*jwtSecret))
	if err != nil {
		return "", err
	}
	return string(tokenByte), nil
}

// Must insure it can get user_id from request
func getUserIDFromReq(r *http.Request) int64 {
	token, _ := r.Cookie("jwt")
	jwtToken := token.Value
	t, _ := jwt.Parse([]byte(jwtToken))
	id, _ := t.Get("id")
	uid, _ := strconv.ParseInt(id.(string), 10, 64)
	return uid
}

// Must insure it can get jwt from cookies
func getJWTFromReq(r *http.Request) string {
	token, _ := r.Cookie("jwt")
	jwtToken := token.Value
	return jwtToken
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

	secretToken := util.GenToken(20)
	err = model.SetUserJWTSecret(r.Context(), uid, &secretToken)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	jwtToken, err := generateJWTToken(uid, &secretToken)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	jwtCookie := http.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		Expires:  time.Now().Add(config.LOGIN_EXPIRETIME),
		HttpOnly: true,
	}

	http.SetCookie(w, &jwtCookie)
	util.SuccessResponse(w, r, "ok")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	uid := getUserIDFromReq(r)
	err := model.DelUserJWTSecret(r.Context(), uid)
	if err != nil {
		util.ErrorResponse(w, r, err.Error(), config.ERR_INTERNAL)
		return
	}

	jwtCookie := http.Cookie{
		Name:     "jwt",
		Value:    getJWTFromReq(r),
		HttpOnly: true,
		MaxAge:   -1,
	}

	http.SetCookie(w, &jwtCookie)
	util.SuccessResponse(w, r, "ok")
}
