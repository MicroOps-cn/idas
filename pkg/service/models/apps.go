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

type Apps []*App

//type AppRole struct {
//	Model
//	AppId     string  `json:"appId" gorm:"type:char(36);not null"`
//	Name      string  `gorm:"type:varchar(50);" json:"name"`
//	Config    string  `json:"config" json:"config"`
//	User      []*User `gorm:"-" json:"user,omitempty"`
//	IsDefault bool    `json:"isDefault" gorm:"not null;default:0"`
//}

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

func (a AppProxyUrls) Len() int {
	return len(a)
}

func (a AppProxyUrls) Less(i, j int) bool {
	return a[i].Index < a[j].Index
}

func (a AppProxyUrls) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type AppProxyConfig struct {
	*AppProxyUrl
	Domain   string `json:"domain" gorm:"type:varchar(50);"`
	Upstream string `json:"upstream" gorm:"type:varchar(50);"`
}
