package passwordRecovery

import (
	"github.com/labstack/echo/v4"
)

// RouterPassRecovery : RouterPassRecovery struct
type RouterPassRecovery struct{}

func (ctrl RouterPassRecovery) Init(g *echo.Group) {
	g.POST("send", ctrl.SendPasswordResetCode)
	g.POST("reset", ctrl.ResetPassword)
}
