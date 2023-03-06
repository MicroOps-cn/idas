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
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/MicroOps-cn/idas/pkg/global"
)

type TokenType string

const (
	TokenTypeParent        TokenType = "parent"
	TokenTypeToken         TokenType = "token"
	TokenTypeRefreshToken  TokenType = "refresh_token"
	TokenTypeCode          TokenType = "code"
	TokenTypeResetPassword TokenType = "reset_password"
	TokenTypeLoginSession  TokenType = "session"
	TokenTypeActive        TokenType = "active"
	TokenTypeAppProxyLogin TokenType = "app_proxy_login"
	TokenTypeTotpSecret    TokenType = "totp_secret"
	TokenTypeLoginCode     TokenType = "login_code"
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
	Childrens  []*Token     `json:"-" gorm:"-"`
}

func (s *Token) GetRelationId() string {
	var ids []string
	if len(s.RelationId) != 0 {
		ids = append(ids, s.RelationId)
	}
	if len(s.Childrens) > 0 {
		for _, children := range s.Childrens {
			childRelId := children.GetRelationId()
			if len(childRelId) > 0 {
				ids = append(ids, childRelId)
			}
		}
	}
	return strings.Join(ids, ",")
}

func NewToken(tokenType TokenType, data ...interface{}) (*Token, error) {
	token := &Token{Id: NewId(), CreateTime: time.Now(), Type: tokenType, LastSeen: time.Now(), Expiry: tokenType.GetExpiry()}
	if len(data) > 1 {
		token.Type = TokenTypeParent
		for _, d := range data {
			rawData, err := json.Marshal(d)
			if err != nil {
				return nil, err
			}
			childToken := &Token{Id: NewId(), ParentId: token.Id, CreateTime: time.Now(), LastSeen: time.Now(), Data: rawData, Type: tokenType, Expiry: tokenType.GetExpiry()}
			if obj, ok := d.(HasId); ok {
				childToken.RelationId = obj.GetId()
			}
			token.Childrens = append(token.Childrens, childToken)
		}
	} else if len(data) == 1 {
		rawData, err := json.Marshal(data[0])
		if err != nil {
			return nil, err
		}
		if obj, ok := data[0].(HasId); ok {
			token.RelationId = obj.GetId()
		}
		token.Data = rawData
	}
	return token, nil
}

func getValElem(valOf reflect.Value) reflect.Value {
	if valOf.Kind() == reflect.Ptr || valOf.Kind() == reflect.Pointer {
		return getValElem(valOf.Elem())
	}
	return valOf
}

func (s *Token) To(r interface{}) error {
	elem := getValElem(reflect.ValueOf(r))
	switch elem.Kind() {
	case reflect.Array, reflect.Slice:
		buf := bytes.NewBuffer([]byte("["))
		if s.Type != TokenTypeParent {
			buf.Write(s.Data)
		} else {
			for idx, children := range s.Childrens {
				if idx != 0 {
					buf.WriteRune(',')
				}
				buf.Write(children.Data)
			}
		}
		buf.WriteRune(']')
		return json.Unmarshal(buf.Bytes(), r)
	default:
		if s.Type == TokenTypeParent {
			return fmt.Errorf("the receiver type does not match, it should be array instead of struct")
		}
		return json.Unmarshal(s.Data, r)

	}
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
