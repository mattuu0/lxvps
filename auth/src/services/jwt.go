package services

import (
	"auth/logger"
	"auth/models"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// 秘密鍵
	JwtPrivateKey ed25519.PrivateKey

	// 公開鍵
	JwtPublicKey ed25519.PublicKey
)

func initJwt(certString string) {
	// PEMブロックの解析
	// PEM形式の秘密鍵を解析
	block, _ := pem.Decode([]byte(certString))
	if block == nil {
		logger.PrintErr("PEMデータの解析に失敗しました")
		return
	}

	// 秘密鍵のパース
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logger.PrintErr("秘密鍵のパースに失敗しました: %v", err)
		return
	}

	// Ed25519秘密鍵の型アサーション
	edPrivateKey, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		logger.PrintErr("キーはEd25519秘密鍵ではありません")
		return
	}

	// 秘密鍵を保存
	JwtPrivateKey = edPrivateKey

	// 秘密鍵の情報を表示
	logger.Println("Ed25519秘密鍵の読み込みに成功しました\n")
	logger.Println("秘密鍵の長さ: %d バイト\n", len(edPrivateKey))

	// 対応する公開鍵を取得 (秘密鍵の後半32バイトが公開鍵)
	publicKey := edPrivateKey.Public().(ed25519.PublicKey)
	logger.Println("対応する公開鍵の長さ: %d バイト\n", len(publicKey))
	logger.Println("対応する公開鍵の値: %v\n", publicKey)

	// 公開鍵を保存
	JwtPublicKey = publicKey
}

const (
	tokenExpiry = time.Minute * 10
)

type AccessTokenClaim struct {
	UserID   string   // ユーザーID
	Labels   []string // ラベル
	ProvCode models.ProviderCode
	ProvUid  string
}

func AccessTokenJwt(args AccessTokenClaim) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		// 有効期限
		"exp": time.Now().Add(tokenExpiry).Unix(),
		// ラベル
		"labels": args.Labels,
		// ユーザーID
		"userID": args.UserID,
		// プロバイダ
		"provCode": args.ProvCode,
		// プロバイダUID
		"provUid": args.ProvUid,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(JwtPrivateKey)

	return tokenString, err
}
