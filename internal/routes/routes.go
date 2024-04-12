package routes

import (
	"ToDoProject/internal/routes/auth"
	"ToDoProject/internal/routes/user"
	"ToDoProject/internal/routes/user/passwordRecovery"
	"ToDoProject/internal/routes/user/todo"
	"github.com/labstack/echo/v4"
)

func Routes(g *echo.Group) {
	auth.RouterAuth{}.Init(g.Group("/auth"))
	user.RouterUser{}.Init(g.Group("/user"))
	passwordRecovery.RouterPassRecovery{}.Init(g.Group("/password/recovery"))
	todo.RouterTodo{}.Init(g.Group("/user/todo"))
}
