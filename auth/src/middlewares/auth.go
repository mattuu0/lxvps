package middlewares

import (
	"auth/logger"
	"auth/models"
	"auth/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// 認証ミドルウェア
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// ヘッダからトークンを取得
		token := ctx.Request().Header.Get("Authorization")
		if token == "" {
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		// トークンを検証
		session, err := services.GetSession(token)
		if err != nil {
			logger.PrintErr(err)
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		// ユーザーを取得する
		user,result := models.GetUser(session.UserID)

		// エラー処理
		if result.Error != nil {
			logger.PrintErr(result.Error)
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		// ユーザーがBANされている時
		if user.IsBanned == 1 {
			return ctx.JSON(http.StatusForbidden, echo.Map{"error": "Your account has been banned"})
		}

		// セッションを設定
		ctx.Set("session", session)

		// 認証処理
		return next(ctx)
	}
}