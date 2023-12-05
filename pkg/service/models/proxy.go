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
	"crypto/sha256"

	"github.com/MicroOps-cn/fuck/crypto"
	uuid "github.com/satori/go.uuid"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/errors"
)

type AppRoleURL struct {
	AppRoleId     string `json:"app_role_id"`
	AppRoleName   string `json:"app_role_name"`
	AppProxyURLId string `json:"app_proxy_url_id"`
}

type AppProxyConfig struct {
	AppProxy
	URLRoles []AppRoleURL
}

func (c *AppProxyConfig) GetId() string {
	return c.AppProxy.GetAppId()
}

type ProxySession struct {
	User  *User
	Proxy *AppProxyConfig
}

func (c *ProxySession) GetId() string {
	return c.User.Id
}

func (c *AppProxy) SetJwtSecret(secret string) (err error) {
	globalSecret := config.Get().GetGlobal().GetSecret()
	if globalSecret == "" {
		return errors.NewServerError(500, "global secret is not set")
	}
	c.JwtSecretSalt = uuid.NewV4().Bytes()
	key := sha256.Sum256([]byte(string(c.JwtSecretSalt) + (globalSecret)))
	c.JwtSecret, err = crypto.NewAESCipher(key[:]).CBCEncrypt([]byte(secret))
	return err
}

func (c *AppProxy) GetJwtSecret() (string, error) {
	if len(c.JwtSecret) == 0 || len(c.JwtSecretSalt) == 0 {
		return "", nil
	}
	globalSecret := config.Get().GetGlobal().GetSecret()
	if globalSecret == "" {
		return "", errors.NewServerError(500, "global secret is not set")
	}
	key := sha256.Sum256([]byte(string(c.JwtSecretSalt) + (globalSecret)))
	sec, err := crypto.NewAESCipher(key[:]).CBCDecrypt(c.JwtSecret)
	return string(sec), err
}
