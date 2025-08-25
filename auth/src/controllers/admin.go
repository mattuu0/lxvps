package controllers

import (
	"auth/logger"
	"auth/models"
	"auth/services"
	"auth/session"
	"net/http"

	"github.com/labstack/echo/v4"
)


func GetAdminStatus(ctx echo.Context) error {
	// admin のステータスを取得
	adstatus,err := services.GetAdminStatus()

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
	}

	return ctx.JSON(http.StatusOK, adstatus)
}

// admin ユーザーを作成する
func CreateAdminUser(ctx echo.Context) error {
	// リクエストボディを取得
	args := services.CreateAdminUserArgs{}

	// bind する
	if err := ctx.Bind(&args); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// admin ユーザーを作成する
	if err := services.CreateAdminUser(args); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})
}

// admin でログインする
func LoginAdminUser(ctx echo.Context) error {
	// bind する
	bindData := services.LoginAdminUserArgs{}

	// bind する
	if err := ctx.Bind(&bindData); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// admin でログインする
	adminUserID,err := services.LoginAdminUser(services.LoginAdminUserArgs{
		Username: bindData.Username,
		Password: bindData.Password,
	})

	// エラー処理
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

		// セッションを作成する
	session,err := session.GetAdminSession(ctx)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)

		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// セッションに adminUserID を保存
	session.Values["adminID"] = adminUserID
	// セッションを保存
	err = session.Save(ctx.Request(), ctx.Response())

	// エラー処理
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func GetAdminInfo(ctx echo.Context) error {
	// ユーザーを取得
	auser := ctx.Get("auser").(*models.AdminUser)

	// パスワードハッシュ以外を返す
	return ctx.JSON(http.StatusOK, models.AdminUser{
		UserID:       auser.UserID,
		Username:     auser.Username,
		PasswordHash: "",
		IsSystem:     auser.IsSystem,
		CreatedAt:    auser.CreatedAt,
	})
}

func AdminLogout(ctx echo.Context) error {
	// セッションを取得
	session,err := session.GetAdminSession(ctx)

	// エラー処理
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// セッションの有効期限を削除する
	session.Options.MaxAge = -1

	session.Values["adminID"] = nil
	// セッションを保存
	err = session.Save(ctx.Request(), ctx.Response())

	// エラー処理
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})
}