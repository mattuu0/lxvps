package services

import (
	"auth/models"
	"time"
)

type CreateLabelArgs struct {
	Name  string `json:"name"`  // ラベル名
	Color string `json:"color"` // ラベル色
}

func CreateLabel(args CreateLabelArgs) error {
	// ラベルを作成する
	return models.CreateLabel(&models.Label{
		Name:  args.Name,
		Color: args.Color,
	})
}

type Label struct {
	// TypeScriptの string に合わせます。
	// これがユニークなIDとして使われることを想定し、JSONキーは "id" です。
	ID string `json:"id"`

	// ラベル名です。JSONキーは "name" です。
	Name string `json:"name"`

	// ラベルの色情報などです。JSONキーは "color" です。
	Color string `json:"color"`

	// 作成日時を表す文字列です。JSONキーは "createdAt" です。
	// time.Time 型に変換したい場合は、JSONデコード後に別途処理が必要です。
	CreatedAt string `json:"createdAt"`
}

func GetLabels() ([]Label, error) {
	// ラベル一覧を取得
	labels, err := models.GetLabels()

	// エラー処理
	if err != nil {
		return []Label{}, err
	}

	// 返す用のラベル
	returnLabels := []Label{}
	for _, val := range labels {
		// 返す用のラベルに追加
		returnLabels = append(returnLabels, Label{
			ID:        val.Name,
			Name:      val.Name,
			Color:     val.Color,
			CreatedAt: FormatUnixTimestampToString(val.CreatedAt, time.RFC3339),
		})
	}

	// 返す
	return returnLabels, nil
}

type DeleteLabelArgs struct {
	ID string `json:"id"`
}

func DeleteLabel(args DeleteLabelArgs) error {
	// ラベルを取得する
	label, err := models.GetLabel(args.ID)
	if err != nil {
		return err
	}

	// ラベルを削除する
	return models.DeleteLabel(label)
}

type LabelUpdateArgs struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

func UpdateLabel(args LabelUpdateArgs) error {
	// ラベルを取得する
	label, err := models.GetLabel(args.ID)
	if err != nil {
		return err
	}

	// ラベルを更新する
	label.Name = args.Name
	label.Color = args.Color
	return models.UpdateLabel(label)
}
