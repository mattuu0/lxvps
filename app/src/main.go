package main

import (
	"app/controllers"
	"app/middlewares"
	"app/models"
	"app/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// モデル初期化
	models.Init()

	// ミドルウェア初期化
	middlewares.Init()

	// サービス初期化
	services.Init()

	// コントローラー初期化
	controllers.Init()

	// ルーター
	router := echo.New()

	// router.Use(middleware.Recover())
	router.Use(middleware.Logger())
	// router.Use(middlewares.PocketAuth())

	router.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"result": "hello world",
		})
	}, middlewares.RequireAuth)

	router.Logger.Fatal(router.Start(":8090"))
}