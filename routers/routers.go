package routers

import (
	"net/http"

	"github.com/HeRaNO/xcpc-team-reg/middleware"
	"github.com/HeRaNO/xcpc-team-reg/modules"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func InitRouters() http.Handler {
	// Register routers...
	r := chi.NewRouter()
	r.Use(chi_middleware.RequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/", modules.SayHello)
		r.Get("/getIDSchool", modules.GetIDSchool)
		r.Get("/getSchoolID", modules.GetSchoolID)
		r.Post("/register", modules.Register)
		r.Post("/verifyEmail", modules.SendValidationEmail)
		r.Get("/getContestInfo", modules.GetContestInfo)

		r.Group(func(r chi.Router) {
			r.Use(middleware.CheckUnauthorized)

			r.Post("/login", modules.Login)
			r.Post("/forgot", modules.ForgotPwd)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticator)

			r.Mount("/admin", adminRouter())

			r.Get("/getUserInfo", modules.GetUserInfo)
			r.Get("/getTeamInfo", modules.GetTeamInfo)

			r.Group(func(r chi.Router) {
				r.Use(middleware.CheckTime)
				r.Post("/modifyUserInfo", modules.ModifyUserInfo)
				r.Post("/modifyTeamInfo", modules.ModifyTeamInfo)

				r.Post("/createTeam", modules.CreateTeam)
				r.Post("/joinTeam", modules.JoinTeam)
				r.Post("/quitTeam", modules.QuitTeam)
			})

			r.Post("/logout", modules.Logout)
		})
	})

	return r
}

func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.CheckAdmin)

	r.Get("/adminHello", modules.SayHelloAdmin)

	r.Get("/export", modules.ExportTeamInfo)
	r.Post("/import", modules.ImportTeamInfo)
	return r
}
