package models

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"idas/pkg/global"
)

type TokenType string

const (
	TokenTypeToken         = "token"
	TokenTypeRefreshToken  = "refresh_token"
	TokenTypeCode          = "code"
	TokenTypeResetPassword = "reset_password"
	TokenTypeLoginSession  = "session"
	TokenTypeActive        = "active"
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
	case TokenTypeLoginSession:
		return time.Now().UTC().Add(global.LoginSessionExpiration)
	case TokenTypeActive:
		return time.Now().UTC().Add(global.ActiveExpiration)
	default:
		return time.Now().UTC().Add(time.Minute * 10)
	}
}

type HasId interface {
	GetId() string
}

type Token struct {
	Id         string       `json:"id" gorm:"primary_key;type:char(36)" valid:"required"`
	CreateTime time.Time    `json:"createTime,omitempty" gorm:"default:now();not null;type:datetime;omitempty"`
	Data       sql.RawBytes `json:"-" gorm:"not null"`
	RelationId string       `json:"relationId" gorm:"type:varchar(1024)"`
	Expiry     time.Time    `json:"expiry,omitempty" gorm:"not null;type:datetime;omitempty"`
	Type       TokenType    `json:"type" gorm:"type:varchar(20)"`
	LastSeen   time.Time    `json:"lastSeen"`
}

func NewToken(tokenType TokenType, data ...interface{}) (*Token, error) {
	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	token := &Token{Data: rawData, Type: tokenType, Expiry: tokenType.GetExpiry()}
	var relationIds []string
	for _, d := range data {
		if x, ok := d.(HasId); ok {
			relationIds = append(relationIds, x.GetId())
		}
	}
	token.RelationId = strings.Join(relationIds, ",")
	return token, nil
}

func (s *Token) BeforeCreate(db *gorm.DB) error {
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
	if s.LastSeen.IsZero() {
		db.Statement.SetColumn("LastSeen", time.Now().UTC())
	}
	return nil
}
