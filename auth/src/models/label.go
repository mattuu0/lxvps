package models

type Label struct {
	ID    uint   `gorm:"primarykey"`      // ラベルのプライマリキー
	Name  string `gorm:"unique;not null"` // ラベル名（ユニークかつNULL不可）
	Color string `gorm:"default:#000000"` // ラベルの色

	CreatedAt int64 `gorm:"autoCreateTime"`	// ラベルの作成日時

	// これも同じ中間テーブル "user_labels" を指定します。
	Users []*User `gorm:"many2many:user_labels;constraint:OnDelete:CASCADE"`
}

func GetLabels() ([]Label, error) {
	var labels []Label

	// 取得する
	err := dbconn.Find(&labels).Error
	return labels, err
}

func GetLabel(name string) (*Label, error) {
	var label Label

	// 取得する
	err := dbconn.Where(&Label{Name: name}).First(&label).Error
	return &label, err
}

// ID からラベルを取得する
func GetLabelByID(id uint) (*Label, error) {
	var label Label

	// 取得する
	err := dbconn.Where(&Label{ID: id}).First(&label).Error
	return &label, err
}

func CreateLabel(label *Label) error {
	return dbconn.Create(label).Error
}

func UpdateLabel(label *Label) error {
	return dbconn.Save(label).Error
}

func DeleteLabel(label *Label) error {
	return dbconn.Delete(label).Error
}

// ユーザーにラベルを追加する
func (usr *User) AddLabel(labelName string) error {
	// ラベルを取得する
	label, err := GetLabel(labelName)
	if err != nil {
		return err
	}

	// 追加する
	return dbconn.Model(usr).Association("Labels").Append(label)
}

// ユーザーのラベルを全て削除する
func (usr *User) RemoveAllLabels() error {
	// 削除する
	return dbconn.Model(usr).Association("Labels").Clear()
}

// ユーザーからラベルを削除する
func (usr *User) RemoveLabel(name string) error {
	// ラベルを取得する
	label, err := GetLabel(name)
	if err != nil {
		return err
	}

	// 削除する
	return dbconn.Model(usr).Association("Labels").Delete(label)
}

// ラベルのリストを返す
func (usr *User) GetLabels() ([]Label, error) {
	// ユーザーのラベルを取得する
	var labels []Label
	err := dbconn.Model(usr).Association("Labels").Find(&labels)
	return labels, err
}

//ラベルの名前リストを返す
func (usr *User) GetLabelNames() ([]string, error) {
	// ユーザーのラベルを取得する
	var labels []Label
	err := dbconn.Model(usr).Association("Labels").Find(&labels)

	// ラベルの名前を返す
	labelNames := []string{}
	for _, label := range labels {
		// ラベルの名前を返す
		labelNames = append(labelNames, label.Name)
	}

	return labelNames, err
}