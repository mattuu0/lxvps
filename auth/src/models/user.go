package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	UserID       string       `gorm:"type:varchar(255);primaryKey"`                             // ユーザーID
	Name         string       `gorm:"type:varchar(255)"`                                        // ユーザー名
	Email        string       `gorm:"type:varchar(255);uniqueIndex:idx_users_email,length:255"` // メールアドレス
	ProvCode     ProviderCode `gorm:"type:varchar(255);index:idx_prov_code,length:255"`         // 認証プロバイダコード
	ProvUID      string       `gorm:"type:varchar(255);index:idx_prov_uid,length:255"`          // 認証プロバイダUID
	PasswordHash string       `gorm:"default:''"`                                               // ハッシュ化されたパスワード
	CreatedAt    int64        `gorm:"autoCreateTime"`                                           // ユーザー作成日
	Sessions     []Session    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`            // ユーザーが持つセッション
	IsBanned     int          `gorm:"default:0"`                                                // ユーザーの禁止状態
	IsSystem     int          `gorm:"default:0"`                                                // システムユーザーかどうか
	Labels       []Label      `gorm:"many2many:user_labels;constraint:OnDelete:CASCADE"`        // ユーザーのラベル
	UpdatedAt    int64        `gorm:"autoUpdateTime"`                                           // ユーザー更新日
}

func CreateUser(user *User, ProviderCode ProviderCode) error {
	// プロバイダを取得する
	provider, err := GetProvider(ProviderCode)

	// エラー処理
	if err != nil {
		return err
	}

	// provider にユーザーを追加する
	user.ProvCode = ProviderCode

	// ユーザを作成する
	err = dbconn.Create(user).Error

	// エラー処理
	if err != nil {
		return err
	}

	// プロバイダにユーザーを追加する
	err = dbconn.Model(provider).Association("Users").Append(user)

	// エラー処理
	if err != nil {
		return err
	}

	return nil
}

type GetResult struct {
	Error error	//エラー型
	IsExists bool	//存在するか
}

// ユーザーを取得
func GetUser(userID string) (*User, GetResult) {
	var user User

	// 取得する
	err := dbconn.Where(&User{UserID: userID}).Preload("Labels").First(&user).Error


	return &user, GetResult{
		Error:    err,
		IsExists: !errors.Is(err,gorm.ErrRecordNotFound),
	}
}

// メールアドレスからユーザーを取得
func GetUserByEmail(email string) (*User, GetResult) {
	var user User

	// 取得する
	err := dbconn.Where(&User{Email: email}).First(&user).Error

	return &user, GetResult{
		Error:    err,
		IsExists: !errors.Is(err,gorm.ErrRecordNotFound),
	}
}

// 全てのユーザーを取得
func GetAllUsers() ([]User, error) {
	var users []User

	// 取得する
	err := dbconn.Find(&users).Error
	return users, err
}

// ユーザーを更新する
func UpdateUser(user *User) error {
	// 更新する
	return dbconn.Save(user).Error
}

// ユーザーを削除する
func DeleteUser(userid string) error {
	// ユーザーを取得
	user,result := GetUser(userid)

	// エラー処理
	if result.Error != nil {
		return result.Error
	}

	return dbconn.Unscoped().Delete(user).Error
}

//ユーザーを検索する
func SearchUserByName(name string) ([]User, error) {
	var users []User
	err := dbconn.Where("name LIKE ?", "%"+name+"%").Find(&users).Error
	return users, err
}

func SearchUserByEmail(email string) ([]User, error) {
	var users []User
	err := dbconn.Where("email LIKE ?", "%"+email+"%").Find(&users).Error
	return users, err
}