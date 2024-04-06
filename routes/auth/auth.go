package auth

import (
	"ToDoProject/app/middlewares"
	"ToDoProject/database/models"
	"github.com/labstack/echo/v4"
)

// RouterAuth AuthRouter : AuthRouter struct
type RouterAuth struct{}

func (ctrl RouterAuth) Init(g *echo.Group) {
	g.POST("/register", ctrl.Register)
	g.POST("/login", ctrl.Login)
	g.POST("/logout", ctrl.Logout, middlewares.Authorized(models.JWTAccess))
	g.POST("/me", ctrl.Me, middlewares.Authorized(models.JWTAccess))
	g.POST("/refreshToken", ctrl.RefreshJWTToken, middlewares.Authorized(models.JWTRefresh))
}
