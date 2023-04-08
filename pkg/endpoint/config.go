/*
 Copyright © 2023 MicroOps-cn.

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

package endpoint

import (
	"context"
	"fmt"

	"github.com/MicroOps-cn/fuck/conv"
	"github.com/go-kit/kit/endpoint"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/service"
)

func MakeGetSecurityConfigEndpoint(_ service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := SimpleResponseWrapper[*config.RuntimeSecurityConfig]{}
		resp.Data = config.GetRuntimeConfig().Security
		return resp, nil
	}
}

func MakePatchSecurityConfigEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchSecurityConfigRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		config.SetRuntimeConfig(func(c *config.RuntimeConfig) {
			dst := map[string]interface{}{}
			if resp.Error = conv.JSON(req, &dst); resp.Error != nil {
				return
			}
			if resp.Error = svc.PatchSystemConfig(ctx, "security", dst); resp.Error != nil {
				return
			}
			if c.Security == nil {
				c.Security = &config.RuntimeSecurityConfig{}
			}
			if req.AccountInactiveLock != nil {
				c.Security.AccountInactiveLock = *req.AccountInactiveLock
			}
			if req.ForceEnableMfa != nil {
				c.Security.ForceEnableMfa = *req.ForceEnableMfa
			}
			if req.PasswordComplexity != nil {
				c.Security.PasswordComplexity = *req.PasswordComplexity
			}
			if req.PasswordMinLength != nil {
				c.Security.PasswordMinLength = *req.PasswordMinLength
			}
			if req.PasswordExpireTime != nil {
				c.Security.PasswordExpireTime = *req.PasswordExpireTime
			}
			if req.PasswordFailedLockThreshold != nil {
				c.Security.PasswordFailedLockThreshold = *req.PasswordFailedLockThreshold
			}
			if req.PasswordFailedLockDuration != nil {
				c.Security.PasswordFailedLockDuration = *req.PasswordFailedLockDuration
			}
			if req.PasswordHistory != nil {
				c.Security.PasswordHistory = *req.PasswordHistory
			}
		})
		fmt.Println(resp)
		return resp, nil
	}
}
