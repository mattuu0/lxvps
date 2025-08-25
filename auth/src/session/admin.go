package session

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func GetAdminSession(ctx echo.Context) (*sessions.Session, error) {
	return AdminSession.Get(ctx.Request(), AdminSessionName)
}