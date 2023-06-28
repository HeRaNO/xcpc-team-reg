package middleware

import (
	"net/http"
	"time"

	"github.com/HeRaNO/xcpc-team-reg/internal"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/logger/accesslog"
	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/cookie"
)

func InitMw(h *server.Hertz) {
	h.Use(accesslog.New())
	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://" + internal.Domain, "https://" + internal.Domain},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           5 * time.Minute,
	}))

	storeSession := cookie.NewStore(internal.SessionSecret)
	storeSession.Options(sessions.Options{
		Path:     "/",
		Domain:   internal.Domain,
		MaxAge:   internal.LoginExpireTime,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	h.Use(sessions.New(internal.SessionName, storeSession))
}
