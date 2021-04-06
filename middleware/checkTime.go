package middleware

import (
	"net/http"
	"time"

	"github.com/HeRaNO/xcpc-team-reg/config"
	"github.com/HeRaNO/xcpc-team-reg/util"
)

func CheckTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nowTime := time.Now()
		if nowTime.Before(config.ContestStartTime) {
			util.ErrorResponse(w, r, "registration has not started yet", config.ERR_OUTOFTIME)
			return
		}
		if nowTime.After(config.ContestEndTime) {
			util.ErrorResponse(w, r, "registration has ended", config.ERR_OUTOFTIME)
			return
		}

		next.ServeHTTP(w, r)
	})
}
