package auth

import (
	"ToDoProject/internal/middlewares"
	"ToDoProject/internal/models"
	"github.com/labstack/echo/v4"
)

// RouterAuth : RouterAuth struct
type RouterAuth struct{}

func (ctrl RouterAuth) Init(g *echo.Group) {
	g.POST("/register", ctrl.Register)
	g.POST("/login", ctrl.Login)
	g.POST("/logout", ctrl.Logout, middlewares.Authorized(models.JWTAccess))
}
