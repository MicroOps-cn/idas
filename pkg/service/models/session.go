package models

import (
	"database/sql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

func (s *Session) BeforeCreate(db *gorm.DB) error {
	if s.Id == "" {
		id := NewId()
		if len(id) != 32 {
			return errors.New("生成ID失败: " + id)
		}
		db.Statement.SetColumn("Id", id)
	}
	if s.CreateTime.IsZero() {
		db.Statement.SetColumn("CreateTime", time.Now().UTC())
	}
	return nil
}

type Session struct {
	Id         string       `json:"id" gorm:"primary_key;type:char(32)" valid:"required"`
	CreateTime time.Time    `json:"createTime,omitempty" gorm:"not null;type:datetime;omitempty"`
	UserId     string       `json:"userId" gorm:"not null;type:char(32)"`
	Key        string       `json:"-" gorm:"not null;type:char(32)"`
	Data       sql.RawBytes `json:"-" gorm:"not null"`
	Expiry     time.Time    `json:"expiry,omitempty" gorm:"not null;type:datetime;omitempty"`
	LastSeen   time.Time    `json:"lastSeen" gorm:"not null;type:datetime;omitempty"`
}
