package middlewares

import (
	"auth/logger"
	"auth/models"
	"auth/session"
	"net/http"

	"github.com/labstack/echo/v4"
)

// 認証ミドルウェア
func RequireAdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// セッションを取得
		session,err := session.GetAdminSession(ctx)

		// エラー処理
		if err != nil {
			logger.PrintErr(err)
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		// adminID を取得
		adminID := session.Values["adminID"]

		// nil の時
		if adminID == nil {
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		// admin ユーザー取得
		auser,result := models.GetAdminUserByUserID(adminID.(string))

		// エラー処理
		if result.Error != nil {
			logger.PrintErr(err)
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		// データを設定
		ctx.Set("auser", auser)

		// 認証処理
		return next(ctx)
	}
}