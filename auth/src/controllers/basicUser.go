package controllers

import (
	"auth/logger"
	"auth/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BasicUserArgs struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// 一般ユーザーを作成
func CreateBasicUser(ctx echo.Context) error {
	// リクエストボディを取得
	args := BasicUserArgs{}

	// バインド
	err := ctx.Bind(&args)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// ユーザーを作成する
	token, result := services.CreateBasicUser(services.CreateBasicUserArgs{
		Name:        args.Name,
		Email:       args.Email,
		Password:    args.Password,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr(result.Error)
		return ctx.JSON(result.Code, echo.Map{"error": result.Error.Error()})
	}

	return ctx.JSON(result.Code, echo.Map{"token": token})
}

type LoginBasicUserArgs struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// 一般ユーザーをログイン
func LoginBasicUser(ctx echo.Context) error {
	// リクエストボディを取得
	args := LoginBasicUserArgs{}

	// バインド
	err := ctx.Bind(&args)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// ユーザーをログインする
	token, result := services.LoginBasicUser(services.LoginBasicUserArgs{
		Email:    args.Email,
		Password: args.Password,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr(result.Error)
		return ctx.JSON(result.Code, echo.Map{"error": result.Error.Error()})
	}

	return ctx.JSON(result.Code, echo.Map{"token": token})
}
