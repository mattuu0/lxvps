package services

import "app/logger"

func Init() {
	logger.Println("LXDを初期化しています...")
	TestLxd()
}