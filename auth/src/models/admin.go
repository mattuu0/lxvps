package models

import (
	"errors"

	"gorm.io/gorm"
)

type AdminUser struct {
	UserID       string	`gorm:"type:varchar(255);primaryKey"`	// 管理者ユーザーID
	Username     string `gorm:"type:varchar(255);uniqueIndex"`
	PasswordHash string `gorm:"default:''"`
	IsSystem     int    `gorm:"default:0"`
	CreatedAt    int64  `gorm:"autoCreateTime"`
	UpdatedAt    int64  `gorm:"autoUpdateTime"`
}

func CreateAdminUser(user *AdminUser) error {
	return dbconn.Create(user).Error
}

func GetAdminUser(username string) (*AdminUser, GetResult) {
	var user AdminUser

	// 取得する
	err := dbconn.Where(&AdminUser{Username: username}).First(&user).Error

	return &user, GetResult{
		Error:    err,
		IsExists: !errors.Is(err,gorm.ErrRecordNotFound),
	}
}

func GetAdminUserByUserID(userID string) (*AdminUser, GetResult) {
	var user AdminUser

	// 取得する
	err := dbconn.Where(&AdminUser{UserID: userID}).First(&user).Error

	return &user, GetResult{
		Error:    err,
		IsExists: !errors.Is(err,gorm.ErrRecordNotFound),
	}
}

func GetAllAdminUsers() ([]AdminUser, error) {
	var users []AdminUser

	// 取得する
	err := dbconn.Find(&users).Error
	return users, err
}

func GetSystemAdminUser() (*AdminUser, error) {
	var user AdminUser

	// 取得する
	err := dbconn.Where(&AdminUser{IsSystem: 1}).First(&user).Error
	return &user, err
}