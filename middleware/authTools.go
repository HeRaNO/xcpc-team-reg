package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/model"
	"github.com/HeRaNO/xcpc-team-reg/util"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("jwt")
		if err != nil {
			util.ErrorResponse(w, r, err.Error(), config.ERR_UNAUTHORIZED)
			return
		}

		tokenFromUser := []byte(token.Value)
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

		id, ok := t.Get("id")
		if !ok {
			util.ErrorResponse(w, r, "payload error", config.ERR_UNAUTHORIZED)
			return
		}

		uid, err := strconv.ParseInt(id.(string), 10, 64)
		if err != nil {
			util.ErrorResponse(w, r, err.Error(), config.ERR_UNAUTHORIZED)
			return
		}

		secret, err := model.GetUserJWTSecret(r.Context(), uid)
		if err != nil {
			util.ErrorResponse(w, r, err.Error(), config.ERR_UNAUTHORIZED)
			return
		}
		if secret == "" {
			util.ErrorResponse(w, r, "login status expired", config.ERR_UNAUTHORIZED)
			return
		}

		tokenItShouldBe, err := jwt.Sign(t, jwa.HS256, []byte(secret))
		if err != nil {
			util.ErrorResponse(w, r, err.Error(), config.ERR_UNAUTHORIZED)
			return
		}

		if !bytes.Equal(tokenFromUser, tokenItShouldBe) {
			util.ErrorResponse(w, r, "authentication failed", config.ERR_UNAUTHORIZED)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckUnauthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("jwt")
		if err == nil {
			util.ErrorResponse(w, r, "login status hasn't expired", config.ERR_UNAUTHORIZED)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _ := r.Cookie("jwt")
		tokenFromUser := []byte(token.Value)
		t, _ := jwt.Parse(tokenFromUser)
		id, _ := t.Get("id")
		uid, _ := strconv.ParseInt(id.(string), 10, 64)
		isAdmin, err := model.GetAdminByUserID(r.Context(), uid)
		if err != nil {
			util.ErrorResponse(w, r, fmt.Sprintf("user_id: %d ", uid)+err.Error(), config.ERR_UNAUTHORIZED)
			return
		}
		if !isAdmin {
			util.ErrorResponse(w, r, fmt.Sprintf("user_id: %d is not an admin", uid), config.ERR_UNAUTHORIZED)
			return
		}

		next.ServeHTTP(w, r)
	})
}
