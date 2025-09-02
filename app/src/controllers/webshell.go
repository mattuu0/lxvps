package controllers

import (
	"app/logger"
	"app/lxdctl"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  8192,
		WriteBufferSize: 8192,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func ConnectShellWs(ctx echo.Context) error {
	// WebSocket をアップグレード
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// 別スレッドで実行する
	err = lxdctl.InitLxd("test-instance", ws)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})
}