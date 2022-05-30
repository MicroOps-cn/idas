package models

import (
	"database/sql"
	"idas/pkg/global"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (s *Session) BeforeCreate(db *gorm.DB) error {
	if s.Id == "" {
		id := NewId()
		if len(id) != 36 {
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
	Id         string       `json:"id" gorm:"primary_key;type:char(36)" valid:"required"`
	CreateTime time.Time    `json:"createTime,omitempty" gorm:"default:now();not null;type:datetime;omitempty"`
	UserId     string       `json:"userId" gorm:"not null;type:char(36)"`
	Data       sql.RawBytes `json:"-" gorm:"not null"`
	Expiry     time.Time    `json:"expiry,omitempty" gorm:"not null;type:datetime;omitempty"`
	LastSeen   time.Time    `json:"lastSeen" gorm:"default:now();not null;type:datetime;omitempty"`
}

type TokenType string

const (
	TokenTypeToken         = "token"
	TokenTypeRefreshToken  = "refresh_token"
	TokenTypeCode          = "code"
	TokenTypeResetPassword = "reset_password"
)

func (t TokenType) GetExpiry() time.Time {
	switch t {
	case TokenTypeCode:
		return time.Now().UTC().Add(global.AuthCodeExpiration)
	case TokenTypeToken:
		return time.Now().UTC().Add(global.TokenExpiration)
	case TokenTypeRefreshToken:
		return time.Now().UTC().Add(global.RefreshTokenExpiration)
	case TokenTypeResetPassword:
		return time.Now().UTC().Add(global.ResetPasswordExpiration)
	default:
		return time.Now().UTC().Add(time.Minute * 10)
	}
}

type Token struct {
	Id         string       `json:"Id"`
	CreateTime time.Time    `json:"createTime,omitempty" gorm:"default:now();not null;type:datetime;omitempty"`
	Data       sql.RawBytes `json:"-" gorm:"not null"`
	RelationId string       `json:"relation_id"`
	Expiry     time.Time    `json:"expiry,omitempty" gorm:"not null;type:datetime;omitempty"`
	Type       TokenType    `json:"type"`
}
