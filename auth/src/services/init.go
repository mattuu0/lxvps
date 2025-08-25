package services

import (
	"auth/logger"
	"auth/models"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	// トークンシークレット
	TokenSecret = "secret"
)

func Init() {
	// 環境変数からトークンシークレットを取得
	TokenSecret = os.Getenv("TOKEN_SECRET")

	// 環境変数から秘密鍵を取得
	certString := os.Getenv("JWT_PRIVATE_KEY")

	logger.Println("TOKEN_SECRET: " + TokenSecret)
	logger.Println("JWT_PRIVATE_KEY: " + certString)

	// 秘密鍵を初期化
	initJwt(certString)

	// 画像一覧を取得
	filepath.Walk(IconDir, func(path string, info fs.FileInfo, err error) error {
		// ユーザーIDに変換
		userid := filepath.Base(GetFileNameWithoutExtension(path))

		// ユーザーを取得する
		_, result := models.GetUser(userid)

		// ユーザーが存在しない場合
		if !result.IsExists {
			targetPath := IconDir + "/" + userid + ".png"

			// 存在しない場合削除
			err := os.Remove(targetPath)

			// エラー処理
			if err != nil {
				logger.PrintErr(err)
			}

			logger.Println("delete icon: " + targetPath)
		}

		return nil
	})

	// adminのenv が設定されている場合
	adminEmail := os.Getenv("AdminEmail")
	adminPassword := os.Getenv("AdminPassword")

	// admin が有効か
	if adminEmail != "" && adminPassword != "" {
		logger.Println("管理者ユーザーを作成しています")
		
		// 管理者用のユーザーを作成する
		err := CreateAdminUser(CreateAdminUserArgs{
			Username: adminEmail,
			Password: adminPassword,
		})

		// エラー処理
		if err != nil {
			logger.PrintErr(err)
		}
	}
}

// GetFileNameWithoutExtension は指定されたパスからファイル名の拡張子なしの部分を取得します。
func GetFileNameWithoutExtension(filePath string) string {
	// 1. パスからファイル名（拡張子含む）を取得
	fileName := filepath.Base(filePath)

	// 2. ファイル名から拡張子を取得
	ext := filepath.Ext(fileName)

	// 3. ファイル名から拡張子部分を取り除く
	// strings.TrimSuffix を使うと、指定したsuffixがなければそのままの文字列が返されるため安全です。
	return strings.TrimSuffix(fileName, ext)
}
