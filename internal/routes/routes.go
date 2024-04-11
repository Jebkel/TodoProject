package routes

import (
	"ToDoProject/internal/routes/auth"
	"ToDoProject/internal/routes/passwordRecovery"
	"ToDoProject/internal/routes/user"
	"github.com/labstack/echo/v4"
)

func Routes(g *echo.Group) {
	auth.RouterAuth{}.Init(g.Group("/auth"))
	user.RouterUser{}.Init(g.Group("/user"))
	passwordRecovery.RouterPassRecovery{}.Init(g.Group("/password/recovery"))
}
