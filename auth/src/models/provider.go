package models

import (
	"auth/logger"
	"os"
)

type ProviderCode string

const (
	Google    ProviderCode = "google"
	Github    ProviderCode = "github"
	Discord   ProviderCode = "discord"
	Microsoft ProviderCode = "microsoftonline"
	Basic     ProviderCode = "basic"
)

var (
	// Oauth プロバイダ一覧
	OauthProvides  = []ProviderCode{
		Google,
		Github,
		Discord,
		Microsoft,
	}
)

type Provider struct {
    ProviderName string       `gorm:"primaryKey"` // 認証プロバイダ名
    ClientID     string       // 認証プロバイダのクライアントID
    ClientSecret string       // 認証プロバイダのクライアントシークレット
    CallbackURL  string       // 認証プロバイダのコールバックURL
    ProviderCode ProviderCode `gorm:"type:varchar(255);uniqueIndex"` // 認証プロバイダのコード
    IsEnabled    int          `gorm:"default:0"` // 認証プロバイダの有効状態
    Users        []User       `gorm:"foreignKey:ProvCode;references:ProviderCode"` // プロバイダが持つユーザー
}

// プロバイダを取得
func GetProvider(providerCode ProviderCode) (*Provider, error) {
	var provider Provider

	// 取得する
	err := dbconn.First(&provider, &Provider{ProviderCode: providerCode}).Error
	return &provider, err
}

// Oauthプロバイダを更新
func UpdateOauthProvider(provider Provider) error {
	// データを保存する
	return dbconn.Save(&provider).Error
}

func CreateProvider(provider *Provider) error {
	return dbconn.Create(provider).Error
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// プロバイダを初期化する
func InitProviders() {
	// Google が有効かどうか
	isGoogleEnabled := os.Getenv("GoogleClientID") != "" && os.Getenv("GoogleClientSecret") != ""

	// Google
	err := CreateProvider(&Provider{
		ProviderName: "Google",
		ClientID:     os.Getenv("GoogleClientID"),
		ClientSecret: os.Getenv("GoogleClientSecret"),
		CallbackURL:  os.Getenv("GoogleCallback"),
		ProviderCode: Google,
		IsEnabled:    BoolToInt(isGoogleEnabled),
		Users:        []User{},
	})

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
	}

	// GitHub が有効かどうか
	isGithubEnabled := os.Getenv("GithubClientID") != "" && os.Getenv("GithubClientSecret") != ""

	// GitHub
	err = CreateProvider(&Provider{
		ProviderName: "GitHub",
		ClientID:     os.Getenv("GithubClientID"),
		ClientSecret: os.Getenv("GithubClientSecret"),
		CallbackURL:  os.Getenv("GithubCallback"),
		ProviderCode: Github,
		IsEnabled:    BoolToInt(isGithubEnabled),
		Users:        []User{},
	})

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
	}

	// Discord が有効かどうか
	isDiscordEnabled := os.Getenv("DiscordClientID") != "" && os.Getenv("DiscordClientSecret") != ""

	// Discord
	err = CreateProvider(&Provider{
		ProviderName: "Discord",
		ClientID:     os.Getenv("DiscordClientID"),
		ClientSecret: os.Getenv("DiscordClientSecret"),
		CallbackURL:  os.Getenv("DiscordCallback"),
		ProviderCode: Discord,
		IsEnabled:    BoolToInt(isDiscordEnabled),
		Users:        []User{},
	})

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
	}

	// Microsoft が有効かどうか
	isMicrosoftEnabled := os.Getenv("MicrosoftClientID") != "" && os.Getenv("MicrosoftClientSecret") != ""

	// Microsoft
	err = CreateProvider(&Provider{
		ProviderName: "Microsoft",
		ClientID:     os.Getenv("MicrosoftClientID"),
		ClientSecret: os.Getenv("MicrosoftClientSecret"),
		CallbackURL:  os.Getenv("MicrosoftCallback"),
		ProviderCode: Microsoft,
		IsEnabled:    BoolToInt(isMicrosoftEnabled),
		Users:        []User{},
	})

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
	}

	// basic
	err = CreateProvider(&Provider{
		ProviderName: "Basic",
		ClientID:     "",
		ClientSecret: "",
		CallbackURL:  "",
		ProviderCode: Basic,
		IsEnabled:    0,
		Users:        []User{},
	})

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
	}

	// 完了
	logger.Println("Providers initialized")
}

// Oauth のプロバイダ取得
func GetOauthProviders() []Provider {
	// 返すデータ
	returnProviders := []Provider{}

	for _, providerName := range OauthProvides {
		// プロバイダを取得する
		provider,err := GetProvider(providerName)

		// エラー処理
		if err != nil {
			logger.PrintErr(err)
			continue
		}

		// リストの追加
		returnProviders = append(returnProviders, *provider)
	}
	
	return returnProviders
}