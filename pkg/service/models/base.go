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
	"encoding/binary"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const midUint64 = uint64(1) << 63

func NewId(seed ...string) string {
	ts := uint64(time.Now().UnixMicro())
	if ts < midUint64 {
		var hash uint64
		for _, s := range seed {
			seedBytes := []byte(s)
			for i := 0; i < len(seedBytes); i += 8 {
				if len(seedBytes[i:]) <= 8 {
					tmp := make([]byte, 8)
					copy(tmp, seedBytes[i:])
					hash += binary.BigEndian.Uint64(tmp)
					break
				}
				hash += binary.BigEndian.Uint64(seedBytes[i : i+8])
			}
		}
		if hash > midUint64 {
			hash = hash - midUint64
		}
		ts += hash
	}

	var id = uuid.NewV4()
	binary.BigEndian.PutUint64(id[:8], ts)
	return id.String()
}

func (m *Model) BeforeCreate(db *gorm.DB) error {
	if m.Id == "" {
		id := NewId()
		if len(id) != 36 {
			return errors.New("Failed to generate ID: " + id)
		}
		db.Statement.SetColumn("Id", id)
	}
	if m.UpdateTime.IsZero() {
		db.Statement.SetColumn("UpdateTime", time.Now().UTC())
	}
	if m.CreateTime.IsZero() {
		db.Statement.SetColumn("CreateTime", time.Now().UTC())
	}
	return nil
}

func (m *Model) BeforeSave(db *gorm.DB) error {
	db.Statement.SetColumn("UpdateTime", time.Now().UTC())
	return nil
}
