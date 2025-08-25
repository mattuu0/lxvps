package services

import (
	"auth/models"
)

type OauthProvider struct {
	ProviderCode string `json:"ProviderCode"`
	ProviderName string `json:"ProviderName"`
	ClientID     string `json:"ClientID"`
	ClientSecret string `json:"ClientSecret"`
	CallbackURL  string `json:"CallbackURL"`
	IsEnabled    int    `json:"IsEnabled"` // JSONの数値に合わせてint型
}

// Oauth プロバイダ一覧を取得
func GetOauthProviders() []OauthProvider {
	// 返すデータ
	returnProviders := []OauthProvider{}

	// データベースから取得
	providers := models.GetOauthProviders()

	for _, provider := range providers {
		// データを返す
		returnProviders = append(returnProviders, OauthProvider{
			ProviderCode: string(provider.ProviderCode),
			ProviderName: provider.ProviderName,
			ClientID:     provider.ClientID,
			ClientSecret: provider.ClientSecret,
			CallbackURL:  provider.CallbackURL,
			IsEnabled:    provider.IsEnabled,
		})
	}

	return returnProviders
}

// 更新する
func UpdateOauthProviders(providers []OauthProvider) error {
	for _, provider := range providers {
		// プロバイダを取得
		getProvider, err := models.GetProvider(models.ProviderCode(provider.ProviderCode))

		// エラー処理
		if err != nil {
			return err
		}

		// データを更新する
		getProvider.CallbackURL = provider.CallbackURL
		getProvider.ClientID = provider.ClientID
		getProvider.ClientSecret = provider.ClientSecret
		getProvider.IsEnabled = provider.IsEnabled

		// プロバイダを更新
		err = models.UpdateOauthProvider(*getProvider)

		// エラー処理
		if err != nil {
			return err
		}
	}

	return nil
}

// Basic プロバイダ更新
type UpdateBasicProviderArgs struct {
	
}

func UpdateBasicProvider() {

}