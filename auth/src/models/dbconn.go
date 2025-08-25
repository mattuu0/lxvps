package models

import (
	"auth/logger"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbconn *gorm.DB = nil
)

func OpenDB() (*gorm.DB,error) {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// エラー処理
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigreteTable(db *gorm.DB) error {
	logger.Println("マイグレーションを実行しています...")

	// マイグレーション
	err := db.AutoMigrate(&User{})

	// エラー処理
	if err != nil {
		logger.PrintErr("User テーブルのマイグレーションに失敗しました", err)
		return err
	}
	
	// マイグレーション
	err = db.AutoMigrate(&Provider{})

	// エラー処理
	if err != nil {
		logger.PrintErr("Provider テーブルのマイグレーションに失敗しました", err)
		return err
	}

	// マイグレーション
	err = db.AutoMigrate(&Session{})

	// エラー処理
	if err != nil {
		logger.PrintErr("Session テーブルのマイグレーションに失敗しました", err)
		return err
	}

	// マイグレーション
	err = db.AutoMigrate(&Label{})

	// エラー処理
	if err != nil {
		logger.PrintErr("Label テーブルのマイグレーションに失敗しました", err)
		return err
	}

	// マイグレーション
	err = db.AutoMigrate(&AdminUser{})

	// エラー処理
	if err != nil {
		logger.PrintErr("AdminUser テーブルのマイグレーションに失敗しました", err)
		return err
	}

	logger.Println("マイグレーションを実行しました")

	return nil
}

func Init() error {
	logger.Println("データベース接続を確立しています...")

	// データベース接続
	db, err := OpenDB()
	if err != nil {
		logger.PrintErr("データベース接続に失敗しました", err)
		return err
	}

	// データベース接続確認
	err = MigreteTable(db)

	// エラー処理
	if err != nil {
		return err
	}

	// グローバル変数に格納
	dbconn = db

	// プロバイダを初期化する
	InitProviders()

	return nil
}

func GetDB() *gorm.DB {
	return dbconn
}