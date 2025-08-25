package models

type Session struct {
    SessionID string `gorm:"primaryKey"` // セッションID
    UserID    string // ユーザーID
    UserAgent string // ユーザーエージェント
    RemoteIP  string // リモートIP
    CreatedAt int64  `gorm:"autoCreateTime"` // セッション作成日
}

// セッションを追加
func (usr *User) NewSession(session *Session) error {
	// ユーザのセッションに追加
	err := dbconn.Model(usr).Association("Sessions").Append(session)

	// エラー処理
	if err != nil {
		return err
	}

	return nil
}

// セッションを削除
func (usr *User) DeleteSession(sessionid string) error {
	// ユーザのセッションから削除
	err := dbconn.Model(usr).Association("Sessions").Unscoped().Delete(&Session{SessionID: sessionid})
	
	// エラー処理
	if err != nil {
		return err
	}
	
	return nil
}

// セッション取得
func GetSession(sessionid string) (*Session, error) {
	var session Session

	// 取得する
	err := dbconn.Where(&Session{SessionID: sessionid}).First(&session).Error
	return &session, err
}

// 全てのセッションを取得
func GetAllSessions() ([]Session,error) {
	var Sessions []Session

	// 取得する
	err := dbconn.Find(&Sessions).Error
	return Sessions, err
}

// セッションを削除
func DeleteSession(sessionid string) error {
	// 取得する
	err := dbconn.Where(&Session{SessionID: sessionid}).Unscoped().Delete(&Session{}).Error
	return err
}