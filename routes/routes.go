package routes

import (
	"ToDoProject/routes/auth"
	"github.com/labstack/echo/v4"
)

func Routes(g *echo.Group) {
	auth.RouterAuth{}.Init(g.Group("/auth"))

}
