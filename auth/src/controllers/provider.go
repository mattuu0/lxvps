package controllers

import (
	"auth/logger"
	"auth/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ProviderConfig は各認証プロバイダーの設定を表す構造体です。
type ProviderConfig struct {
	Enabled      bool   `json:"enabled"`
	ClientID     string `json:"clientId"`     // JSONキーは camelCase
	ClientSecret string `json:"clientSecret"` // JSONキーは camelCase
	CallbackURL  string `json:"callbackUrl"`  // JSONキーは camelCase
}

// AuthConfig は認証プロバイダー全体のコンフィグレーションを表す構造体です。
// JSONのトップレベルのキー（"google", "discord" など）に対応するフィールドを持ちます。
type AuthConfig struct {
	Google    ProviderConfig `json:"google"`
	Discord   ProviderConfig `json:"discord"`
	Github    ProviderConfig `json:"github"`
	Microsoft ProviderConfig `json:"microsoft"`
	// 他のプロバイダーがあればここに追加
	// MyCustomProvider ProviderConfig `json:"mycustomprovider"`
}

// プロバイダ一覧を取得する関数
func GetProviders(ctx echo.Context) error {
	// プロバイダ一覧を取得
	return ctx.JSON(http.StatusOK,echo.Map{"message": "success"})
}

func UpdateProviders(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK,echo.Map{"message": "success"})
}

// Oauth プロバイダを取得
func GetOauthProviders(ctx echo.Context) error {
	// サービスから取得
	providers := services.GetOauthProviders()

	return ctx.JSON(http.StatusOK,providers)
}

// basic プロバイダ更新
func BasicUpdate(ctx echo.Context) error {
	// Basic 認証を更新する
	// services.UpdateBasicProvider()
	return nil
}

// プロバイダを更新
func UpdateOauthProviders(ctx echo.Context) error {
	bindData := []services.OauthProvider{}

	// bind する
	if err := ctx.Bind(&bindData); err != nil {
		logger.PrintErr(err)

		return ctx.JSON(http.StatusBadRequest,echo.Map{
			"error" : err.Error(),
		})
	}

	// 更新する
	err := services.UpdateOauthProviders(bindData)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)

		return ctx.JSON(http.StatusInternalServerError,echo.Map{
			"error" : err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK,echo.Map{
		"result" : "success",
	})
}