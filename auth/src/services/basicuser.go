package services

import (
	"auth/models"
	"auth/structs"
	"auth/utils"
	"errors"
	"net/http"
)

type CreateBasicUserArgs struct {
	Name     string // ユーザー名
	Email    string // メールアドレス
	Password string // パスワード
}

// 一般ユーザーを作成する (返却値: トークン, HttpResult)
func CreateBasicUser(args CreateBasicUserArgs) (string, structs.HttpResult) {
	// プロバイダを取得
	provider, err := models.GetProvider(models.Basic)

	// エラー処理
	if err != nil {
		return "", structs.HttpResult{
			Code: http.StatusInternalServerError,
			Message: "failed to get provider",
			Error:   err,
			Success: false,
		}
	}

	if provider.IsEnabled == 0 {
		return "", structs.HttpResult{
			Code: http.StatusUnauthorized,
			Message: "provider is disabled",
			Error:   errors.New("provider is disabled"),
			Success: false,
		}
	}

	// UUID を生成
	uid := utils.GenID()

	// 現在時刻を取得
	now := utils.NowTime()

	// ユーザーを取得する
	_, result := models.GetUserByEmail(args.Email)

	// エラー処理
	if result.IsExists {
		// 存在するとき
		return "", structs.HttpResult{
			Code: http.StatusConflict,
			Message: "user already exists",
			Error:   err,
			Success: false,
		}
	}

	// パスワードをハッシュ化する
	hashed, err := utils.HashPassword(args.Password)

	// エラー処理
	if err != nil {
		return "", structs.HttpResult{
			Code: http.StatusInternalServerError,
			Message: "failed to hash password",
			Error:   err,
			Success: false,
		}
	}

	// ユーザーを作成する
	err = models.CreateUser(&models.User{
		UserID:       uid,
		Name:         args.Name,
		Email:        args.Email,
		ProvCode:     "",
		ProvUID:      "",
		PasswordHash: hashed,
		CreatedAt:    now,
	}, models.Basic)

	// エラー処理
	if err != nil {
		return "", structs.HttpResult{
			Code: http.StatusInternalServerError,
			Message: "failed to create user",
			Error:   err,
			Success: false,
		}
	}

	// セッションを作成する
	token, err := NewSession(SessionArgs{
		UserID:   uid,
		RemoteIP: "",
		UserAgent: "",
	})

	// エラー処理
	if err != nil {
		return "", structs.HttpResult{
			Code: http.StatusInternalServerError,
			Message: "failed to create session",
			Error:   err,
			Success: false,
		}
	}

	return token, structs.HttpResult{
		Code: http.StatusOK,
		Message: "user created",
		Error:   nil,
		Success: true,
	}
}

type LoginBasicUserArgs struct {
	Email    string // メールアドレス
	Password string // パスワード
}

// ログインしてトークンを返す (返却値: トークン, HttpResult)
func LoginBasicUser(args LoginBasicUserArgs) (string,structs.HttpResult) {
	// プロバイダを取得
	provider, err := models.GetProvider(models.Basic)

	// エラー処理
	if err != nil {
		return "",structs.HttpResult{
			Code: http.StatusInternalServerError,
			Message: "failed to get provider",
			Error:   err,
			Success: false,
		}
	}

	if provider.IsEnabled == 0 {
		return "",structs.HttpResult{
			Code: http.StatusUnauthorized,
			Message: "provider is disabled",
			Error:   errors.New("provider is disabled"),
			Success: false,
		}
	}

	// ユーザーを取得する
	user, result := models.GetUserByEmail(args.Email)

	// エラー処理
	if result.Error != nil {
		return "",structs.HttpResult{
			Code: http.StatusInternalServerError,
			Message: "failed to get user",
			Error:   err,
			Success: false,
		}
	}

	// プロバイダをチェックする
	if user.ProvCode != models.Basic {
		// basic 以外の場合はエラーを返す
		return "",structs.HttpResult{
			Code: http.StatusBadRequest,
			Message: "invalid provider",
			Error:   errors.New("invalid provider"),
			Success: false,
		}
	}

	// パスワードをチェックする
	if !utils.CheckPasswordHash(args.Password, user.PasswordHash) {
		// パスワードが一致しない場合はエラーを返す
		return "",structs.HttpResult{
			Code: http.StatusBadRequest,
			Message: "invalid password",
			Error:   errors.New("invalid password"),
			Success: false,
		}
	}

	// セッションを作成する
	token, err := NewSession(SessionArgs{
		UserID:    user.UserID,
		RemoteIP:  "",
		UserAgent: "",
	})

	// エラー処理
	if err != nil {
		return "",structs.HttpResult{
			Code: http.StatusInternalServerError,
			Message: "failed to create session",
			Error:   err,
			Success: false,
		}
	}

	return token,structs.HttpResult{
		Code: http.StatusOK,
		Message: "success",
		Error:   nil,
		Success: true,
	}
}
