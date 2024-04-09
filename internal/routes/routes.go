package routes

import (
	"ToDoProject/internal/routes/auth"
	"ToDoProject/internal/routes/user"
	"github.com/labstack/echo/v4"
)

func Routes(g *echo.Group) {
	auth.RouterAuth{}.Init(g.Group("/auth"))
	user.RouterUser{}.Init(g.Group("/user"))
}
