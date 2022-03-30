package models

import (
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"time"
)

type UserStatus int8

const (
	UserStatusUnknown UserStatus = iota
	UserStatusNormal
	UserStatusDisable
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	Model
	Username    string       `gorm:"type:varchar(20);" json:"username"`
	Salt        sql.RawBytes `gorm:"type:tinyblob;" json:"-" `
	Password    sql.RawBytes `gorm:"type:tinyblob;" json:"password,omitempty"`
	Email       string       `gorm:"type:varchar(50);" json:"email"`
	PhoneNumber string       `json:"phoneNumber"`
	FullName    string       `gorm:"type:varchar(20);" json:"fullName"`
	Avatar      string       `json:"avatar"`
	Status      UserStatus   `gorm:"not null;default:0" json:"status"`
	LoginTime   time.Time    `json:"loginTime"`
	Role        UserRole     `json:"role"`
	Storage     string       `gorm:"-" json:"storage"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type plain User
	u.Password = nil
	return json.Marshal(plain(u))
}

func (u User) GenSecret(password ...string) []byte {
	sha := sha1.New()
	sha.Write([]byte(u.Salt))
	if len(password) > 0 {
		sha.Write([]byte(password[0]))
	} else {
		sha.Write(u.Password)
	}
	return sha.Sum(nil)
}
