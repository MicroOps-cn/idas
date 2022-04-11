package models

type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
)

type GrantMode int8

const (
	GrantModeManual GrantMode = 0 // 手动授权
	GrantModeFull   GrantMode = 1 // 全员均可登陆
)

type GroupStatus uint8

const (
	GroupStatusUnknown GroupStatus = iota
	GroupUserStatusNormal
	GroupStatusDisable
)

type App struct {
	Model
	Name        string      `gorm:"type:varchar(50);" json:"name"`
	Description string      `gorm:"type:varchar(200);" json:"description"`
	Avatar      string      `gorm:"type:varchar(200)" json:"avatar"`
	GrantType   GrantType   `gorm:"type:varchar(20);" json:"grantType"`
	GrantMode   GrantMode   `gorm:"type:varchar(20);" json:"grantMode"`
	Storage     string      `gorm:"-" json:"storage"`
	User        []*User     `gorm:"many2many:t_app_user" json:"user,omitempty"`
	Status      GroupStatus `gorm:"not null;default:0" json:"status"`
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
