package middleware

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/util"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) < 8 || strings.ToUpper(token[0:6]) != "BEARER" {
			util.ErrorResponse(w, r, "wrong authentication method", config.ERR_UNAUTHORIZED)
			return
		}
		fgpCookie, err := r.Cookie(config.JWTFingerPrintName)
		if err != nil {
			util.ErrorResponse(w, r, err.Error(), config.ERR_UNAUTHORIZED)
			return
		}

		tokenFromUser := []byte(token[7:])
		secretFromUser := fgpCookie.Value
		if secretFromUser == "" {
			util.ErrorResponse(w, r, "wrong finger print status", config.ERR_UNAUTHORIZED)
			return
		}

		t, err := jwt.Parse(tokenFromUser)
		if err != nil {
			util.ErrorResponse(w, r, err.Error(), config.ERR_UNAUTHORIZED)
			return
		}
		err = jwt.Validate(t)
		if err != nil {
			util.ErrorResponse(w, r, err.Error(), config.ERR_UNAUTHORIZED)
			return
		}

		tokenItShouldBe, err := jwt.Sign(t, jwa.HS256, config.JWTSecret)
		if err != nil {
			util.ErrorResponse(w, r, err.Error(), config.ERR_UNAUTHORIZED)
			return
		}

		if !bytes.Equal(tokenFromUser, tokenItShouldBe) {
			util.ErrorResponse(w, r, "authentication failed", config.ERR_UNAUTHORIZED)
			return
		}

		id, ok := t.Get(config.JWTIDName)
		if !ok {
			util.ErrorResponse(w, r, "payload error", config.ERR_UNAUTHORIZED)
			return
		}
		admin, ok := t.Get(config.JWTAdminName)
		if !ok {
			util.ErrorResponse(w, r, "payload error", config.ERR_UNAUTHORIZED)
			return
		}
		secretFromJWT, ok := t.Get(config.JWTFingerPrintName)
		if !ok {
			util.ErrorResponse(w, r, "payload error", config.ERR_UNAUTHORIZED)
			return
		}

		uid, err := strconv.ParseInt(id.(string), 10, 64)
		if err != nil {
			util.ErrorResponse(w, r, err.Error(), config.ERR_UNAUTHORIZED)
			return
		}

		isAdmin := false
		if admin.(string) == "1" {
			isAdmin = true
		} else if admin.(string) != "0" {
			util.ErrorResponse(w, r, "payload error", config.ERR_UNAUTHORIZED)
			return
		}

		if util.SHA256([]byte(secretFromUser)) != secretFromJWT.(string) {
			util.ErrorResponse(w, r, "authentication failed", config.ERR_UNAUTHORIZED)
			return
		}

		ctxUserInfo := config.CtxUserInfo{
			ID:    uid,
			Admin: isAdmin,
		}
		ctx := context.WithValue(r.Context(), config.CtxUserInfoName, &ctxUserInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CheckUnauthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(config.JWTFingerPrintName)
		if err == nil {
			util.ErrorResponse(w, r, "login status hasn't expired", config.ERR_UNAUTHORIZED)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := r.Context().Value(config.CtxUserInfoName).(*config.CtxUserInfo)
		if !userInfo.Admin {
			util.ErrorResponse(w, r, "permission denied", config.ERR_UNAUTHORIZED)
			return
		}
		next.ServeHTTP(w, r)
	})
}
