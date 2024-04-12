package todo

import (
	"ToDoProject/internal/middlewares"
	"ToDoProject/internal/models"
	"github.com/labstack/echo/v4"
)

type RouterTodo struct{}

func (ctrl RouterTodo) Init(g *echo.Group) {
	g.POST("/create", ctrl.createTask, middlewares.Authorized(models.JWTAccess))
}
