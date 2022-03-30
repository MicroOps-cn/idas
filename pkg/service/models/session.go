package models

import (
	"database/sql"
	"time"
)

type Session struct {
	Key    string `gorm:"primary_key;type:char(32)"`
	Data   sql.RawBytes
	Expiry time.Time `gorm:"type:datetime;omitempty" json:"update_time,omitempty"`
}
