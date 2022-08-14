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
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type User struct {
	Model
	Username    string              `gorm:"type:varchar(20);unique" json:"username"`
	Salt        sql.RawBytes        `gorm:"type:tinyblob;" json:"-" `
	Password    sql.RawBytes        `gorm:"type:tinyblob;" json:"password,omitempty"`
	Email       string              `gorm:"type:varchar(50);" json:"email" valid:"email,optional"`
	PhoneNumber string              `gorm:"type:varchar(50);" json:"phoneNumber" valid:"numeric,optional"`
	FullName    string              `gorm:"type:varchar(50);" json:"fullName"`
	Avatar      string              `gorm:"type:varchar(128);" json:"avatar"`
	Status      UserMeta_UserStatus `gorm:"not null;default:0" json:"status"`
	LoginTime   *time.Time          `json:"loginTime,omitempty"`
	RoleId      string              `gorm:"->;-:migration" json:"roleId,omitempty"`
	Role        string              `gorm:"->;-:migration" json:"role,omitempty"`
	App         []*App              `gorm:"many2many:app_user" json:"app,omitempty"`
	Storage     string              `gorm:"-" json:"storage"`
}

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

type UserKey struct {
	Model
	Name    string `gorm:"type:varchar(50)" json:"name"`
	User    *User  `gorm:"-" json:"-"`
	UserId  string `gorm:"type:char(36);" json:"userId"`
	Key     string `gorm:"type:varchar(50);" json:"key"`
	Secret  string `gorm:"type:varchar(50);" json:"secret"`
	Private string `gorm:"-" json:"key2,omitempty"`
}
