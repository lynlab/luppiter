package keyvalue

import (
	"errors"
	"time"

	"github.com/lynlab/luppiter"
)

var (
	KeyValueInvalidKey   = errors.New("invalid key")
	KeyValueErrNotExists = errors.New("item not exists")
)

// Item 은 키-밸류 스토리지 값입니다.
type Item struct {
	ID int64

	Namespace string `gorm:"type:varchar(255); unique_index:unique_namespace_key"`
	Key       string `gorm:"type:varchar(255); unique_index:unique_namespace_key"`
	Value     string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Item) TableName() string {
	return "key_value_items"
}

func GetKeyValueItem(namespace, key string) (string, error) {
	var item Item
	luppiter.DB.Where(&Item{Namespace: namespace, Key: key}).First(&item)
	if item.ID == 0 {
		return "", KeyValueErrNotExists
	}
	return item.Value, nil
}

func SetKeyValueItem(namespace, key, value string) {
	luppiter.DB.Save(&Item{Namespace: namespace, Key: key, Value: value})
}

func init() {
	luppiter.DB.AutoMigrate(&Item{})
}
