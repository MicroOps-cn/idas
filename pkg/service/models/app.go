package models

type App struct {
	Model
	Name        string            `gorm:"type:varchar(50);not null;unique" json:"name"`
	Description string            `gorm:"type:varchar(200);" json:"description"`
	Avatar      string            `gorm:"type:varchar(128)" json:"avatar"`
	GrantType   AppMeta_GrantType `gorm:"type:TINYINT(3);not null;default:0"  json:"grantType"`
	GrantMode   AppMeta_GrantMode `gorm:"type:TINYINT(3)not null;default:0" json:"grantMode"`
	Status      AppMeta_Status    `gorm:"type:TINYINT(3)not null;default:0" json:"status"`
	User        []*User           `gorm:"many2many:app_user" json:"user,omitempty"`
	Role        AppRoles          `gorm:"foreignKey:AppId" json:"role,omitempty"`
	Storage     string            `gorm:"-" json:"storage"`
}

type AppRole struct {
	Model
	AppId     string  `json:"appId" gorm:"type:char(36);not null"`
	Name      string  `gorm:"type:varchar(50);" json:"name"`
	Config    string  `json:"config" json:"config"`
	User      []*User `gorm:"-" json:"user,omitempty"`
	IsDefault bool    `json:"isDefault" gorm:"not null;default:0"`
}

type AppRoles []*AppRole

func (roles AppRoles) GetRole(name string) *AppRole {
	for _, role := range roles {
		if role.Name == name {
			return role
		}
	}
	return nil
}

type AppUser struct {
	Model
	AppId  string `json:"appId" gorm:"type:char(36);not null"`
	App    *App   `json:"app,omitempty"`
	UserId string `json:"userId" gorm:"type:char(36);not null"`
	User   *User  `json:"user,omitempty"`
	RoleId string `json:"roleId" gorm:"default:'';type:char(36);not null"`
}

type AppAuthCode struct {
	Model
	SessionId string `json:"session_id" gorm:"type:CHAR(36);not null"`
	AppId     string `json:"appId" gorm:"type:CHAR(36);not null"`
	Scope     string `json:"scope" gorm:"type:varchar(128);not null"`
	Storage   string `json:"storage" gorm:"type:varchar(128);not null"`
}
