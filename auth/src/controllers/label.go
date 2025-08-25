package controllers

import (
	"auth/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ラベルを作成する関数
func CreateLabel(ctx echo.Context) error {
	// リクエストボディを取得
	args := services.CreateLabelArgs{}

	// bind する
	if err := ctx.Bind(&args); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// ラベルを作成する
	if err := services.CreateLabel(args); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func GetLabels(ctx echo.Context) error {
	// ラベルを取得する
	labels, err := services.GetLabels()

	// エラー処理
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, labels)
}

func DeleteLabel(ctx echo.Context) error {
	// リクエストボディを取得
	args := services.DeleteLabelArgs{}

	// bind する
	if err := ctx.Bind(&args); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// ラベルを削除する
	if err := services.DeleteLabel(args); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})
}


// ラベルを更新する
func UpdateLabel(ctx echo.Context) error {
	// リクエストボディを取得
	args := services.LabelUpdateArgs{}

	// bind する
	if err := ctx.Bind(&args); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// ラベルを更新する
	if err := services.UpdateLabel(args); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})
}