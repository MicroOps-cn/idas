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
	Name        string      `gorm:"type:varchar(50);not null" json:"name"`
	Description string      `gorm:"type:varchar(200);" json:"description"`
	Avatar      string      `gorm:"type:varchar(200)" json:"avatar"`
	GrantType   GrantType   `gorm:"type:varchar(20);" json:"grantType"`
	GrantMode   GrantMode   `gorm:"type:varchar(20);" json:"grantMode"`
	Status      GroupStatus `gorm:"not null;default:0" json:"status"`
	User        []*User     `gorm:"many2many:t_app_user" json:"user,omitempty"`
	Role        []*AppRole  `gorm:"-" json:"role,omitempty"`
	Storage     string      `gorm:"-" json:"storage"`
}

type AppRole struct {
	Model
	Name string `gorm:"type:varchar(50);" json:"name"`
	//Description string  `gorm:"type:varchar(50);" json:"description"`
	User      []*User `gorm:"-" json:"user,omitempty"`
	IsDefault bool    `json:"isDefault" gorm:"not null;default:0"`
}

type AppUser struct {
	Model
	AppId  string `json:"appId" gorm:"type:char(32);not null" json:"appId"`
	App    *App   `json:"app,omitempty"`
	UserId string `json:"userId" gorm:"type:char(32);not null" json:"userId"`
	User   *User  `json:"user,omitempty"`
	Role   string `json:"role" gorm:"default: '';not null"`
}
