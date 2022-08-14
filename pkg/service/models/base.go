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
