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
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/MicroOps-cn/idas/pkg/utils/sign"
	"reflect"
	"strings"
	"time"

	"github.com/MicroOps-cn/fuck/crypto"
	uuid "github.com/satori/go.uuid"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/errors"
)

func (u User) MarshalJSON() ([]byte, error) {
	type plain User
	u.Password = nil
	return json.Marshal(plain(u))
}

func (u User) GenSecret(password ...string) []byte {
	sha := sha1.New()
	sha.Write(u.Salt)
	if len(password) > 0 {
		sha.Write([]byte(password[0]))
	} else {
		sha.Write(u.Password)
	}
	return sha.Sum(nil)
}

func (u User) GetAttr(name string) string {
	if len(name) == 0 {
		return ""
	}
	ot := reflect.TypeOf(u)
	for i := 0; i < ot.NumField(); i++ {
		ft := ot.Field(i)
		if cut, _, _ := strings.Cut(ft.Tag.Get("json"), ","); len(cut) != 0 {
			if cut == name {
				val := reflect.ValueOf(u).Field(i).Interface()
				switch v := val.(type) {
				case string:
					return v
				case []byte:
					return string(v)
				default:
					return fmt.Sprint(v)
				}
			}
		}
	}
	return ""
}

func (u User) IsForceMfa() bool {
	return config.GetRuntimeConfig().Security.ForceEnableMfa || (u.ExtendedData != nil && u.ExtendedData.ForceMFA)
}

type Users []*User

func (u Users) Id() (ids []string) {
	for _, user := range u {
		ids = append(ids, user.Id)
	}
	return
}

func (u Users) GetById(id string) *User {
	for _, user := range u {
		if user.Id == id {
			return user
		}
	}
	return nil
}

type UserKey struct {
	Model
	Name    string `gorm:"type:varchar(50)" json:"name"`
	UserId  string `gorm:"type:char(36);" json:"userId"`
	Key     string `gorm:"type:varchar(50);" json:"key"`
	Secret  string `gorm:"type:varchar(50);" json:"secret"`
	Private string `gorm:"-" json:"private,omitempty"`
}

type AppKey struct {
	Model
	Name   string `gorm:"type:varchar(50)" json:"name"`
	AppId  string `gorm:"type:char(36);" json:"appId"`
	Key    string `gorm:"type:varchar(50);" json:"key"`
	Secret string `gorm:"type:varchar(50);" json:"secret"`
}

type UserExt struct {
	UserId             string `json:"userId" gorm:"primary_key;type:char(36)"`
	ForceMFA           bool
	TOTPSecret         sql.RawBytes `json:"-" gorm:"column:totp_secret;type:tinyblob"`
	TOTPSalt           sql.RawBytes `json:"-" gorm:"column:totp_salt;type:tinyblob"`
	EmailAsMFA         bool         `json:"emailAsMFA" gorm:"column:email_as_mfa"`
	SmsAsMFA           bool         `json:"smsAsMFA" gorm:"column:sms_as_mfa"`
	TOTPAsMFA          bool         `json:"totpAsMFA" gorm:"column:totp_as_mfa"`
	PasswordModifyTime time.Time    `json:"passwordModifyTime" gorm:"column:password_modify_time"`
	LoginTime          time.Time    ` json:"loginTime" gorm:"column:login_time"`
	ActivationTime     time.Time    ` json:"activationTime" gorm:"column:activation_time"`
}

func (u *UserExt) SetSecret(secret string) (err error) {
	globalSecret := config.Get().GetGlobal().GetSecret()
	if globalSecret == "" {
		return errors.NewServerError(500, "global secret is not set")
	}
	u.TOTPSalt = uuid.NewV4().Bytes()
	key := sha256.Sum256([]byte(string(u.TOTPSalt) + (globalSecret)))
	u.TOTPSecret, err = crypto.NewAESCipher(key[:]).CBCEncrypt([]byte(secret))
	return err
}

func (u UserExt) GetSecret() (secret string, err error) {
	if len(u.TOTPSecret) == 0 || len(u.TOTPSalt) == 0 {
		return "", nil
	}
	globalSecret := config.Get().GetGlobal().GetSecret()
	if globalSecret == "" {
		return "", errors.NewServerError(500, "global secret is not set")
	}
	key := sha256.Sum256([]byte(string(u.TOTPSalt) + (globalSecret)))
	sec, err := crypto.NewAESCipher(key[:]).CBCDecrypt(u.TOTPSecret)
	return string(sec), err
}

type UserPasswordHistory struct {
	Model
	UserId string `gorm:"type:char(36);" json:"userId"`
	Hash   string `gorm:"type:char(88);" json:"hash"`
}

func (h *UserPasswordHistory) SetPassword(p string) error {
	if len(h.UserId) == 0 {
		return fmt.Errorf("User ID is empty. ")
	}
	h.Hash = sign.SumSha512Hmac(h.UserId, p)
	return nil
}

type WeakPassword struct {
	Id   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Hash string `gorm:"type:char(88);uniqueIndex:idx_hash" json:"hash"`
}
