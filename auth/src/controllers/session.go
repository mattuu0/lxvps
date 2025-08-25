package controllers

import (
	"auth/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetSessions(ctx echo.Context) error {
	// サービスを呼び出す
	sessions, err := services.GetAllSessions();

	// エラー処理
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, sessions)
}

// セッションを削除
func DeleteSession(ctx echo.Context) error {
	// ヘッダからセッションIDを取得
	sessionID := ctx.Request().Header.Get("sessionid")

	// サービスを呼び出す
	err := services.DeleteSession(sessionID)

	// エラー処理
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})
}