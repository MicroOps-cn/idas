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
	"context"
	"database/sql/driver"
	"fmt"
	"strings"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/idas/config"
	"github.com/go-kit/log/level"

	jwtutils "github.com/MicroOps-cn/idas/pkg/utils/jwt"
)

type Apps []*App

func (a Apps) Id() []string {
	ids := make([]string, len(a))
	for idx, app := range a {
		ids[idx] = app.Id
	}
	return ids
}

func (a Apps) GetById(id string) *App {
	for _, app := range a {
		if id == app.Id {
			return app
		}
	}
	return nil
}

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

func (roles AppRoles) GetRoleById(id string) *AppRole {
	for _, role := range roles {
		if role.Id == id {
			return role
		}
	}
	return nil
}

func (roles AppRoles) GetId() (ids []string) {
	for _, role := range roles {
		ids = append(ids, role.Id)
	}
	return
}

type AppUsers []*AppUser

func (s AppUsers) Id() (ids []string) {
	for _, user := range s {
		ids = append(ids, user.Id)
	}
	return
}

func (s AppUsers) GetByUserId(id string) *AppUser {
	for _, user := range s {
		if user.UserId == id {
			return user
		}
	}
	return nil
}

func (s AppUsers) UserId() (ids []string) {
	for _, user := range s {
		ids = append(ids, user.UserId)
	}
	return
}

func (s AppUsers) GetByAppId(id string) *AppUser {
	for _, user := range s {
		if user.AppId == id {
			return user
		}
	}
	return nil
}

type AppUser struct {
	Model
	AppId  string `json:"appId" gorm:"type:char(36);not null;index:idx_app_user,unique"`
	UserId string `json:"userId" gorm:"type:char(36);not null;index:idx_app_user,unique"`
	RoleId string `json:"roleId" gorm:"default:'';type:char(36);not null"`
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

func (a AppProxyUrls) Id() (ids []string) {
	for _, url := range a {
		ids = append(ids, url.Id)
	}
	return
}

type AuthorizedRedirectUrls []string

func (c *AuthorizedRedirectUrls) GormDataType() string {
	return "string"
}

// Scan implements the Scanner interface.
func (c *AuthorizedRedirectUrls) Scan(value any) error {
	switch vt := value.(type) {
	case []uint8:
		*c = strings.Split(string(vt), "\n")
	case string:
		*c = strings.Split(vt, "\n")
	default:
		return fmt.Errorf("failed to resolve field, type exception: %T", value)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (c AuthorizedRedirectUrls) Value() (driver.Value, error) {
	return strings.Join(c, "\n"), nil
}

func NewAuthorizedRedirectUrls(urls []string) AuthorizedRedirectUrls {
	res := make(AuthorizedRedirectUrls, len(urls))
	copy(res, urls)
	return res
}

func (m AppOAuth2) TableName() string {
	return "t_app_oauth2"
}

func (m AppOAuth2) GetJWTIssuer(ctx context.Context) jwtutils.JWTIssuer {
	logger := logs.GetContextLogger(ctx)

	if m.JwtSignatureMethod != AppMeta_default && m.JwtSignatureKey != nil {
		signKey, err := m.JwtSignatureKey.UnsafeString()
		if err != nil {
			level.Error(logger).Log("msg", "failed to decrypt jwt signature key", "err", err)
			return config.Get().GetJwtIssuer()
		}
		issuer, err := jwtutils.NewJWTIssuer(m.AppId, m.JwtSignatureMethod.String(), signKey)
		if err != nil {
			level.Error(logger).Log("msg", "failed to init jwt issuer", "err", err)
			return config.Get().GetJwtIssuer()
		}
		return issuer
	}
	return config.Get().GetJwtIssuer()
}

func (m *App) GetJWTIssuer(ctx context.Context) jwtutils.JWTIssuer {
	if m.OAuth2 == nil {
		return config.Get().GetJwtIssuer()
	}
	if len(m.OAuth2.AppId) == 0 {
		m.OAuth2.AppId = m.Id
	}
	return m.OAuth2.GetJWTIssuer(ctx)
}
