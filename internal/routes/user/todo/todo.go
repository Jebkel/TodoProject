package todo

import (
	"github.com/labstack/echo/v4"
)

type RouterTodo struct{}

func (ctrl RouterTodo) Init(g *echo.Group) {
	g.POST("/create", ctrl.createTask)
	g.POST("/update", ctrl.updateTask)
	g.POST("/switchStatus", ctrl.switchStatusTask)
	g.POST("/delete", ctrl.deleteTask)
	g.POST("/all", ctrl.getTasks)
}
