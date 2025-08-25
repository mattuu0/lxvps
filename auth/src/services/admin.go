package services

import (
	"auth/models"
	"auth/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminStatus struct {
	HasSystemUser bool //システムユーザーが存在するか
}

func GetAdminStatus() (AdminStatus, error) {
	// データベースから取得する
	systemUser, err := models.GetSystemAdminUser()

	// エラー処理
	if err != nil {
		return AdminStatus{}, err
	}

	// 存在しない時
	if err == gorm.ErrRecordNotFound {
		return AdminStatus{
			HasSystemUser: false,
		}, nil
	}

	// 存在する時
	if systemUser != nil {
		return AdminStatus{
			HasSystemUser: true,
		}, nil
	}

	return AdminStatus{
		HasSystemUser: false,
	}, nil
}

// ここから admin ユーザー作成
type CreateAdminUserArgs struct {
	Username string	`json:"username"` //ユーザー名
	Password string `json:"password"` //パスワード
}

// admin ユーザーを作成する
func CreateAdminUser(args CreateAdminUserArgs) error {
	// ステータスを確認する
	status, err := GetAdminStatus()
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// ステータスを確認する
	if status.HasSystemUser {
		// すでにシステムユーザーが存在する場合
		return errors.New("admin user already exists")
	}

	// システムユーザーが存在しない場合

	// ユーザーを作成する
	// パスワードをハッシュ化する
	hashed,err := bcrypt.GenerateFromPassword([]byte(args.Password), 15)

	// エラー処理
	if err != nil {
		return err
	}

	// systen ユーザーを作成する
	return models.CreateAdminUser(&models.AdminUser{
		UserID:       utils.GenID(),
		Username:     args.Username,
		PasswordHash: string(hashed),
		IsSystem:     1,
	})
}

// ここまで
// ここから admin ログイン
type LoginAdminUserArgs struct {
	Username string	`json:"username"` //ユーザー名
	Password string `json:"password"` //パスワード
}

// admin ログイン
func LoginAdminUser(args LoginAdminUserArgs) (string, error) {
	// admin ユーザーを取得する
	user, result := models.GetAdminUser(args.Username)

	// エラー処理
	if result.Error != nil {
		return "", result.Error
	}

	// パスワードをチェックする
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(args.Password)); err != nil {
		return "", err
	}

	return user.UserID, nil
}