package services

import (
	"auth/logger"
	"auth/models"
	"mime/multipart"
	"os"
	"strconv"
	"time"
)

type GetUserInfo struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
}

func GetInfo(userid string) (GetUserInfo, error) {
	// ユーザー取得
	user, result := models.GetUser(userid)

	// エラー処理
	if result.Error != nil {
		return GetUserInfo{}, result.Error
	}

	return GetUserInfo{
		UserID:   user.UserID,
		Name:     user.Name,
	}, nil
}

type UserInfo struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	ProvCode string `json:"prov_code"`
	ProvUid  string `json:"prov_uid"`
}

func GetMe(userid string) (UserInfo, error) {
	// ユーザー取得
	user, result := models.GetUser(userid)

	// エラー処理
	if result.Error != nil {
		return UserInfo{}, result.Error
	}

	return UserInfo{
		UserID:   user.UserID,
		Name:     user.Name,
		Email:    user.Email,
		ProvCode: string(user.ProvCode),
		ProvUid:  user.ProvUID,
	}, nil
}

// ここからユーザーの更新
type UpdateUserData struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Avatar string   `json:"avatar"`
	Labels []string `json:"labels"` // JSONの文字列配列はGoのスライス ([]string) で受けます
}

// ユーザーを更新する関数
func UpdateUser(args UpdateUserData) error {
	// ユーザーを取得
	user, result := models.GetUser(args.ID)

	// エラー処理
	if result.Error != nil {
		return result.Error
	}

	// ユーザーを更新する
	user.Name = args.Name

	// ラベルを削除する
	err := user.RemoveAllLabels()

	// エラー処理
	if err != nil {
		return err
	}

	// ラベルを回す
	for _, labelName := range args.Labels {
		// ラベルを追加
		err = user.AddLabel(labelName)

		// エラー処理
		if err != nil {
			return err
		}
	}

	// ユーザーを更新
	err = models.UpdateUser(user)

	// エラー処理
	if err != nil {
		return err
	}

	// 10mbまでの画像を保存
	err = ProcessAndSaveImage(IconDir + "/" + args.ID + ".png", args.Avatar, MaxImageSize)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
	}

	return nil
}

// ここまで

// ここからユーザー削除
func DeleteUser(userid string) error {
	// user を取得する
	user, result := models.GetUser(userid)

	// エラー処理
	if result.Error != nil {
		return result.Error
	}

	// 画像ファイルを削除する
	err := os.Remove(IconDir + "/" + user.UserID + ".png")

	// エラー処理
	if err != nil {
		return err
	}

	// ユーザーを削除する
	return models.DeleteUser(userid)
}

// ここまで

// ここからユーザー一覧取得
type User struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Provider   string   `json:"provider"`
	ProviderID string   `json:"providerId"`
	Avatar     string   `json:"avatar"`
	Labels     []string `json:"labels"`
	CreatedAt  string   `json:"createdAt"` // 日時型にする場合は time.Time を使用し、適切なフォーマットでパース・フォーマットする必要があります
	Banned     bool     `json:"banned"`
}

func GetUsers() ([]User, error) {
	// ユーザーを取得
	users, err := models.GetAllUsers()

	// エラー処理
	if err != nil {
		return []User{}, err
	}

	userResponse := []User{}
	for _, user := range users {
		// ラベルを取得
		labels, err := user.GetLabelNames()

		// エラー処理
		if err != nil {
			return []User{}, err
		}

		// ユーザーを返す
		userResponse = append(userResponse, User{
			ID:         user.UserID,
			Name:       user.Name,
			Email:      user.Email,
			Provider:   string(user.ProvCode),
			ProviderID: user.ProvUID,
			Avatar:     "/auth/assets/" + user.UserID + ".png?uptime=" + strconv.FormatInt(user.UpdatedAt, 10), //TODO : 本番環境ではパスを変更できるようにする
			Labels:     labels,
			CreatedAt:  FormatUnixTimestampToString(user.CreatedAt, time.RFC3339),
			Banned:     user.IsBanned == 1,
		})
	}

	return userResponse, nil
}

// ここまで

func FormatUnixTimestampToString(timestamp int64, layout string) string {
	// Unixタイムスタンプ (秒) を time.Time に変換
	// time.Unix(seconds, nanoseconds) を使用
	t := time.Unix(timestamp, 0) // ナノ秒は0とします

	// time.Time を指定されたレイアウトで文字列にフォーマット
	return t.Format(layout)
}

// ここから BAN の処理
type BanArgs struct {
	IsBanned bool   //BANするかどうか
	UserID   string //ユーザーID
}

func ToggleBan(args BanArgs) error {
	// ユーザーを取得する
	user,result := models.GetUser(args.UserID)

	// エラー処理
	if result.Error != nil {
		return result.Error
	}

	// BAN を切り替え
	if args.IsBanned {
		// BANする
		user.IsBanned = 1
	} else {
		// BAN解除
		user.IsBanned = 0
	}

	// ユーザーを更新する
	return models.UpdateUser(user)
}

// ここまで

// ここからユーザーのアイコンを更新
type UpdateIconArgs struct {
	UserID string
	ImgFile multipart.File
}

func UpdateIcon(args UpdateIconArgs) error {
	// ユーザーを取得する
	user, result := models.GetUser(args.UserID)

	// エラー処理
	if result.Error != nil {
		return result.Error
	}

	// 画像をリサイズして保存する
	_, err := SaveResizedImage(args.ImgFile, user.UserID, IconWidth, IconHeight, IconDir)

	// エラー処理
	if err != nil {
		return err
	}

	// 更新日時を更新する
	user.UpdatedAt = time.Now().Unix()

	// ユーザーを更新する
	return models.UpdateUser(user)
}

// ここまで

// ここから
// アイコンを取得する
func GetIcon(userid string) (string, error) {
	// ユーザーを取得する
	user, result := models.GetUser(userid)

	// エラー処理
	if result.Error != nil {
		return "", result.Error
	}

	return "/auth/assets/" + user.UserID + ".png?uptime=" + strconv.FormatInt(user.UpdatedAt, 10), nil
}
// ここまで