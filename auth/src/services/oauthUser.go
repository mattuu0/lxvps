package services

import (
	"auth/models"
	"auth/utils"
	"errors"
)

type OauthUserArgs struct {
	Name           string // ユーザー名
	Email          string // メールアドレス
	ProviderCode   string // 認証プロバイダコード
	ProviderUserID string // 認証プロバイダユーザーID
	RemoteIP       string // IPアドレス
	UserAgent      string // User-Agent
	AvaterURL      string // アバターURL
}

// Oauthユーザーを作成する
func LoginOauthUser(args OauthUserArgs) (string, error) {
	// UUID を生成
	uid := utils.GenID()

	// 現在時刻を取得
	now := utils.NowTime()

	// メールアドレスがない時
	if args.Email == "" {
		return "", errors.New("メールアドレスの取得に失敗しました")
	}

	// ユーザーを取得する
	user, result := models.GetUserByEmail(args.Email)

	// 存在する時
	if result.IsExists {
		// ユーザーが取得できた時

		// プロバイダが同じかどうか
		if user.ProvUID != args.ProviderUserID {
			return "", errors.New("同一プロバイダのユーザーが見つかりません")
		}

		// セッションを追加する
		token, err := NewSession(SessionArgs{
			UserID:    user.UserID,
			RemoteIP:  args.RemoteIP,
			UserAgent: args.UserAgent,
		})

		return token, err
	}

	// 存在しない時

	// ユーザーを作成する
	err := models.CreateUser(&models.User{
		UserID:       uid,
		Name:         args.Name,
		Email:        args.Email,
		PasswordHash: "",
		ProvUID:      args.ProviderUserID,
		CreatedAt:    now,
	}, models.ProviderCode(args.ProviderCode))

	// エラー処理
	if err != nil {
		return "", err
	}

	// 画像を保存する (10mb まで)
	if args.AvaterURL != "" {
		// 画像を保存
		err = ProcessImageFromURL(IconDir + "/" + uid + ".png", args.AvaterURL, MaxImageSize, 10)

		// エラー処理
		if err != nil {
			return "", err
		}
	}

	// トークンを生成
	token, err := NewSession(SessionArgs{
		UserID:    uid,
		RemoteIP:  args.RemoteIP,
		UserAgent: args.UserAgent,
	})

	// エラー処理
	if err != nil {
		return "", err
	}

	return token, nil
}
