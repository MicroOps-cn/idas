package models

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func NewId() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

type Model struct {
	Id         string    `json:"id" gorm:"primary_key;type:char(32)" valid:"required"`
	CreateTime time.Time `json:"createTime,omitempty" gorm:"type:datetime;omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty" gorm:"type:datetime;omitempty"`
	IsDelete   bool      `json:"isDelete" gorm:"not null;default:0"`
}

func (model *Model) BeforeCreate(db *gorm.DB) error {
	if model.Id == "" {
		id := NewId()
		if len(id) != 32 {
			return errors.New("生成ID失败: " + id)
		}
		db.Statement.SetColumn("Id", id)
	}
	if model.UpdateTime.IsZero() {
		db.Statement.SetColumn("UpdateTime", time.Now().UTC())
	}
	if model.CreateTime.IsZero() {
		db.Statement.SetColumn("CreateTime", time.Now().UTC())
	}
	return nil
}

func (model *Model) BeforeSave(db *gorm.DB) error {
	db.Statement.SetColumn("UpdateTime", time.Now().UTC())
	return nil
}
