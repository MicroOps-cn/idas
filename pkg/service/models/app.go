package models

type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
)

type GrantMode int8

const (
	GrantModeManual GrantMode = 0 // 全员均可登陆
	GrantModeFull   GrantMode = 1 // 手动授权
)

type App struct {
	Model
	Name        string    `gorm:"type:varchar(50);"`
	Description string    `gorm:"type:varchar(50);"`
	Avatar      string    `gorm:"type:varchar(200)"`
	GrantType   GrantType `gorm:"type:varchar(20);" json:"grantType"`
	GrantMode   GrantMode `gorm:"type:varchar(20);" json:"grantMode"`
}

type AppRole struct {
	Model
	Name        string `gorm:"type:varchar(50);"`
	Description string `gorm:"type:varchar(50);"`
}

type AppUser struct {
	Model
	AppId  string `json:"appId" gorm:"type:char(32)" json:"appId"`
	App    *App   `json:"app,omitempty"`
	UserId string `json:"userId" gorm:"type:char(32)" json:"userId"`
	User   *User  `json:"user,omitempty"`
}
