package models

import (
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func NewId() string {
	return uuid.NewV4().String()
}

type Model struct {
	Id         string    `json:"id" gorm:"primary_key;type:char(36)" valid:"required"`
	CreateTime time.Time `json:"createTime,omitempty" gorm:"type:datetime;not null;omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty" gorm:"type:datetime;not null;omitempty"`
	IsDelete   bool      `json:"isDelete" gorm:"not null;default:0"`
}

func (model *Model) GetId() string {
	return model.Id
}

func (model *Model) BeforeCreate(db *gorm.DB) error {
	if model.Id == "" {
		id := NewId()
		if len(id) != 36 {
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
