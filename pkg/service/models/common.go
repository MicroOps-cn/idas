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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type File struct {
	Model
	Name     string `gorm:"type:varchar(128);" `
	Path     string `gorm:"type:varchar(256);" `
	MimiType string `gorm:"type:varchar(50);" `
	Size     int64
}

type Secret string

func (s Secret) MarshalYAML() (interface{}, error) {
	return "<secret>", nil
}

func (s Secret) MarshalJSON() ([]byte, error) {
	return []byte(`"<secret>"`), nil
}

type JSON json.RawMessage

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func (JSON) GormDataType() string {
	return "json"
}

// MarshalJSON returns m as the JSON encoding of m.
func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

type Event struct {
	Model
	UserId   string        `json:"userId" gorm:"type:varchar(36);column:user_id"`
	Username string        `json:"username" gorm:"type:varchar(40);column:username"`
	Action   string        `json:"action" gorm:"type:varchar(255);column:action"`
	ClientIP string        `json:"clientIp" gorm:"type:varchar(50);column:client_ip"`
	Location string        `json:"location" gorm:"type:varchar(50);column:location"`
	Status   bool          `json:"status" gorm:"column:status"`
	Took     time.Duration `json:"took" gorm:"column:took"`
	Message  string        `json:"message"  gorm:"type:varchar(255);column:message"`
}

func (m *Event) BeforeCreate(db *gorm.DB) error {
	if len(m.Message) > 255 {
		db.Statement.SetColumn("error", m.Message[:255])
	}
	if len(m.Action) > 255 {
		db.Statement.SetColumn("error", m.Action[:255])
	}
	return m.Model.BeforeCreate(db)
}

func (m *Event) BeforeSave(db *gorm.DB) error {
	if len(m.Message) > 255 {
		db.Statement.SetColumn("error", m.Message[:255])
	}
	if len(m.Action) > 255 {
		db.Statement.SetColumn("error", m.Action[:255])
	}
	return m.Model.BeforeSave(db)
}

type EventLog struct {
	Model
	EventId string        `json:"eventId" gorm:"type:varchar(36);column:event_id"`
	Log     CompressField `json:"log"`
}
