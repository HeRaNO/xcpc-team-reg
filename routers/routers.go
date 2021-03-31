package routers

import (
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/modules"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouters() http.Handler {
	// Register routers...
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", modules.SayHello)
		r.Get("/getSchool", modules.GetSchool)

		r.Post("/login", modules.Login)
		r.Post("/register", modules.Register)
		r.Post("/verifyEmail", modules.SendValidationEmail)
		r.Post("/forgot", modules.ForgotPwd)

		r.Group(func(r chi.Router) {
			r.Use(modules.Authenticator)

			r.Mount("/admin", adminRouter())

			r.Get("/getUserInfo", modules.GetUserInfo)
			r.Post("/getTeamInfo", modules.GetTeamInfo)

			r.Post("/modifyUserInfo", modules.ModifyUserInfo)
			r.Post("/modifyTeamInfo", modules.ModifyTeamInfo)

			r.Post("/createTeam", modules.CreateTeam)
			r.Post("/joinTeam", modules.JoinTeam)
			r.Post("/quitTeam", modules.QuitTeam)
			r.Post("/logout", modules.Logout)
		})
	})

	return r
}

func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(modules.CheckAdmin)

	r.Post("/export", modules.ExportTeamInfo)
	r.Post("/createContest", modules.CreateContest)
	r.Post("/import", modules.ImportTeamInfo)

	return r
}
