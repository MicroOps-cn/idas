/*
 Copyright © 2023 MicroOps-cn.

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
	"compress/zlib"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
)

type FieldsConfig []*FieldConfig

func (c *FieldsConfig) GormDataType() string {
	return "blob"
}

// Scan implements the Scanner interface.
func (c *FieldsConfig) Scan(value any) error {
	var val []byte
	switch vt := value.(type) {
	case []uint8:
		val = vt
	case string:
		val = []byte(vt)
	default:
		return fmt.Errorf("failed to resolve field, type exception: %T", value)
	}
	if len(val) > 0 {
		if val[0] == 0x78 {
			r, err := zlib.NewReader(bytes.NewBuffer(val))
			if err != nil {
				return err
			}
			val, err = io.ReadAll(r)
			if err != nil && err != io.ErrUnexpectedEOF {
				return err
			}
		}
	}
	return json.Unmarshal(val, c)
}

// Value implements the driver Valuer interface.
func (c FieldsConfig) Value() (driver.Value, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	w := zlib.NewWriter(buf)
	_, err = w.Write(marshal)
	if err != nil {
		return nil, err
	}
	if err = w.Flush(); err != nil {
		return nil, err
	}
	return buf.String(), err
}

type PageConfig struct {
	Model
	Name        string       `gorm:"type:varchar(50);unique" json:"name"`
	Description string       `gorm:"type:varchar(250)" json:"description"`
	Fields      FieldsConfig `gorm:"type:blob" json:"fields,omitempty"`
	Icon        string       `gorm:"type:varchar(128)" json:"icon"`
	IsDisable   bool         `json:"isDisable"`
}

type PageData struct {
	Model
	PageId string `json:"pageId" valid:"required"`
	Data   *JSON  `json:"data"`
}
