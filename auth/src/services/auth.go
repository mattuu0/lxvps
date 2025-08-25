package services

import "auth/models"

func Logout(session *models.Session) error {
	// ユーザー取得
	user, result := models.GetUser(session.UserID)

	// エラー処理
	if result.Error != nil {
		return result.Error
	}

	// セッションを削除
	return user.DeleteSession(session.SessionID)
}
