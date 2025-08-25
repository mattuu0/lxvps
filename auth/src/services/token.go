package services

import "auth/models"

func GetAccessToken(userID string) (string, error) {
	// ユーザーを取得
	user, result := models.GetUser(userID)

	// エラー処理
	if result.Error != nil {
		return "", result.Error
	}

	// ラベルの名前を返す
	labels := []string{}

	for _, label := range user.Labels {
		// ラベルの名前を返す
		labels = append(labels, label.Name)
	}

	// トークンを生成
	token, err := AccessTokenJwt(AccessTokenClaim{UserID: userID, Labels: labels, ProvCode: user.ProvCode, ProvUid: user.ProvUID})

	return token, err
}
