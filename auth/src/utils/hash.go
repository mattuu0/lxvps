package utils

import "golang.org/x/crypto/bcrypt"

const (
	// bcrypt のコスト
	bcryptCost = 10
)

// パスワードをハッシュ化
func HashPassword(Password string) (string,error) {
	// パスワードハッシュ生成
	passbin, err := bcrypt.GenerateFromPassword([]byte(Password),bcryptCost)

	return string(passbin), err
}

// ハッシュをチェック
func CheckPasswordHash(Password string, Hash string) bool {
	// ハッシュチェック
	return bcrypt.CompareHashAndPassword([]byte(Hash),[]byte(Password)) == nil
}