/*
 Copyright © 2022 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package models

import (
	"database/sql"
	"encoding/json"
	g "github.com/MicroOps-cn/fuck/generator"
	"github.com/MicroOps-cn/idas/config"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/MicroOps-cn/idas/pkg/global"
)

type TokenType string

const (
	TokenTypeToken         TokenType = "token"
	TokenTypeRefreshToken  TokenType = "refresh_token"
	TokenTypeCode          TokenType = "code"
	TokenTypeResetPassword TokenType = "reset_password"
	TokenTypeLoginSession  TokenType = "session"
	TokenTypeActive        TokenType = "active"
	TokenTypeAppProxyLogin TokenType = "app_proxy_login"
	TokenTypeTotpSecret    TokenType = "totp_secret"
	TokenTypeLoginCode     TokenType = "login_code"
	TokenTypeEnableMFA     TokenType = "enable_mfa"
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
		expirationTime := config.GetRuntimeConfig().GetLoginSessionInactivityTime()
		if expirationTime == 0 {
			expirationTime = 30 * 24
		}
		return time.Now().UTC().Add(time.Duration(expirationTime) * time.Hour)
	case TokenTypeActive:
		return time.Now().UTC().Add(global.ActiveExpiration)
	case TokenTypeLoginCode:
		return time.Now().UTC().Add(time.Minute * 3)
	default:
		return time.Now().UTC().Add(time.Minute * 10)
	}
}

type HasId interface {
	GetId() string
}

type Token struct {
	Id         string       `json:"id" gorm:"primary_key;type:char(36)" valid:"required"`
	CreateTime time.Time    `json:"createTime,omitempty" gorm:"not null;type:datetime;omitempty"`
	Data       sql.RawBytes `json:"-" gorm:"not null"`
	RelationId string       `json:"relationId" gorm:"type:char(36)"`
	Expiry     time.Time    `json:"expiry,omitempty" gorm:"not null;type:datetime;omitempty"`
	Type       TokenType    `json:"type" gorm:"type:varchar(20)"`
	LastSeen   time.Time    `json:"lastSeen"`
	ParentId   string       `json:"parentId" gorm:"type:char(36)"`
}

func (s *Token) GetRelationId() string {
	return s.RelationId
}

func NewToken(tokenType TokenType, data interface{}) (*Token, error) {
	token := &Token{Id: g.NewId(string(tokenType)), CreateTime: time.Now(), Type: tokenType, LastSeen: time.Now(), Expiry: tokenType.GetExpiry()}
	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if obj, ok := data.(HasId); ok {
		token.RelationId = obj.GetId()
		token.Id = g.NewId(token.RelationId)
	}
	token.Data = rawData
	return token, nil
}

func getValElem(valOf reflect.Value) reflect.Value {
	if valOf.Kind() == reflect.Ptr || valOf.Kind() == reflect.Pointer {
		return getValElem(valOf.Elem())
	}
	return valOf
}

func (s *Token) To(r interface{}) error {
	return json.Unmarshal(s.Data, r)
}

func (s *Token) BeforeCreate(db *gorm.DB) error {
	if s.Id == "" {
		id := g.NewId(db.Statement.Table)
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

type Counter struct {
	Id         string     `json:"id" gorm:"primary_key;type:char(36)" valid:"required"`
	Seed       string     `json:"seek" gorm:"type:varchar(128)"`
	Count      int64      `json:"count"`
	ExpireTime *time.Time `json:"expireTime"`
}
