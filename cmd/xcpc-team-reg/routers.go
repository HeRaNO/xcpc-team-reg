package main

import (
	"github.com/HeRaNO/xcpc-team-reg/internal/handlers"
	"github.com/HeRaNO/xcpc-team-reg/internal/middleware"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func registerRouter(h *server.Hertz) {
	middleware.InitMw(h)

	h.GET("/", handlers.Ping)
	h.GET("/getIDSchool", handlers.GetIDSchool)
	h.GET("/getContestInfo", handlers.GetContestInfo)
	h.POST("/sendEmail", handlers.SendValidationEmail)

	unAuthGroup := h.Group("", middleware.CheckUnauthorized())
	unAuthGroup.POST("/register", handlers.Register)
	unAuthGroup.POST("/forgot", handlers.ForgotPwd)
	unAuthGroup.POST("/login", handlers.Login)

	authGroup := h.Group("", middleware.Authenticator())
	authGroup.GET("/getUserInfo", handlers.GetUserInfo)
	authGroup.GET("/getTeamInfo", handlers.GetTeamInfo)
	authGroup.POST("/logout", handlers.Logout)

	beforeEndAuthGroup := authGroup.Group("", middleware.CheckBeforeEnd())
	beforeEndAuthGroup.POST("/modifyUserInfo", handlers.ModifyUserInfo)

	regTimeAuthGroup := beforeEndAuthGroup.Group("", middleware.CheckAfterBegin())
	regTimeAuthGroup.POST("/modifyTeamInfo", handlers.ModifyTeamInfo)
	regTimeAuthGroup.POST("/createTeam", handlers.CreateTeam)
	regTimeAuthGroup.POST("/joinTeam", handlers.JoinTeam)
	regTimeAuthGroup.POST("/quitTeam", handlers.QuitTeam)
}
