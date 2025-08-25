package controllers

import (
	"auth/logger"
	"auth/models"
	"auth/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetToken(ctx echo.Context) error {
	// セッションを取得
	session, ok := ctx.Get("session").(*models.Session)

	// エラー処理
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	// ユーザーIDを取得
	userID := session.UserID

	// トークンを取得
	token, err := services.GetAccessToken(userID)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success", "token": token})
}