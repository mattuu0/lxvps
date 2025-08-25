package services

import (
	"auth/models"
	"auth/utils"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type SessionArgs struct {
	UserID    string // ユーザーID
	RemoteIP  string // リモートIP
	UserAgent string // ユーザーエージェント
}

func GenSessionToken(SessionID string) (string, error) {
	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"SessionID": SessionID,
	})

	// トークンに署名
	signedToken, err := token.SignedString([]byte(TokenSecret))

	// エラー処理
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// トークンを検証
func ValidateSessionToken(tokenString string) (string, error) {
	// トークンを検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(TokenSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))

	// エラー処理
	if err != nil {
		return "", err
	}

	// トークンを検証
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["SessionID"].(string), nil
	}

	return "", nil
}

// セッションを作成してトークンを返す
func NewSession(args SessionArgs) (string, error) {
	// ユーザーIDを取得
	user, result := models.GetUser(args.UserID)

	// エラー処理
	if result.Error != nil {
		return "", result.Error
	}

	// BANされている時
	if user.IsBanned == 1 {
		return "", errors.New("Your account has been banned")
	}

	// セッションIDを生成
	SessionID := utils.GenID()

	// セッションを作成
	session := models.Session{
		SessionID: SessionID,
		UserID:    args.UserID,
		RemoteIP:  args.RemoteIP,
		UserAgent: args.UserAgent,
	}

	// セッションを追加
	if err := user.NewSession(&session); err != nil {
		return "", err
	}

	// トークンを生成
	token, err := GenSessionToken(SessionID)

	return token, err
}

func GetSession(tokenString string) (*models.Session, error) {
	// トークンを検証
	SessionID, err := ValidateSessionToken(tokenString)

	// エラー処理
	if err != nil {
		return nil, err
	}

	// セッションを取得
	return models.GetSession(SessionID)
}

// ここからセッション一覧
// Session represents a user session.
type Session struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	IPAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent"`
	CreatedAt int64  `json:"createdAt"`
	ExpiresAt int64  `json:"expiresAt"` // 最終アクティブから有効期限に変更
	IsActive  bool   `json:"isActive"`
}

func GetAllSessions() ([]Session, error) {
	// サービスを呼び出す
	sessions, err := models.GetAllSessions()

	// エラー処理
	if err != nil {
		return nil, err
	}

	// 返すデータ
	returnSessions := make([]Session, len(sessions))

	// セッションを回す
	for i, session := range sessions {
		returnSessions[i] = Session{
			ID:        session.SessionID,
			UserID:    session.UserID,
			IPAddress: session.RemoteIP,
			UserAgent: session.UserAgent,
			CreatedAt: session.CreatedAt * 1000,
			ExpiresAt: session.CreatedAt * 1000 + 1000*60*60*24*30,
			IsActive:  true,
		}
	}

	return returnSessions, nil
}

// ここまで

// ここからセッション削除
func DeleteSession(SessionID string) error {
	// サービスを呼び出す
	return models.DeleteSession(SessionID)
}