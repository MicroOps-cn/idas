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

import "strings"

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
	User  Users
	Proxy *AppProxyConfig
}

func (c *ProxySession) GetId() string {
	return strings.Join(c.User.Id(), ",")
}
