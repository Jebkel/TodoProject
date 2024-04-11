package user

import (
	"ToDoProject/internal/middlewares"
	"ToDoProject/internal/models"
	"github.com/labstack/echo/v4"
)

// RouterUser : RouterUser struct
type RouterUser struct{}

func (ctrl RouterUser) Init(g *echo.Group) {
	g.POST("/me", ctrl.Me, middlewares.Authorized(models.JWTAccess))
	g.POST("/refreshToken", ctrl.RefreshJWTToken, middlewares.Authorized(models.JWTRefresh))
	g.POST("/password/update", ctrl.UpdatePassword, middlewares.Authorized(models.JWTAccess))
}
