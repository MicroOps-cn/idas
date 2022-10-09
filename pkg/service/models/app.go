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
	Proxy       *AppProxy         `gorm:"foreignKey:AppId" json:"proxy,omitempty"`
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
	AppId  string `json:"appId" gorm:"type:char(36);not null;index:idx_app_user,unique"`
	App    *App   `json:"app,omitempty"`
	UserId string `json:"userId" gorm:"type:char(36);not null;index:idx_app_user,unique"`
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

type AppProxyUrls []*AppProxyUrl

//func (m *AppProxyUrl) UnmarshalJSONPB(unmarshaler *jsonpb.Unmarshaler, bytes []byte) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//var _ jsonpb.JSONPBUnmarshaler = &AppProxyUrl{}
