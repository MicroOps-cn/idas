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

	"github.com/go-kit/kit/endpoint"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/utils/version"
)

func MakeGetGlobalConfigEndpoint(_ service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		globalConfig := config.Get().GetGlobal()
		resp := &GlobalConfig{
			Title:            globalConfig.Title,
			SubTitle:         globalConfig.SubTitle,
			Logo:             globalConfig.Logo,
			Copyright:        globalConfig.Copyright,
			DefaultLoginType: LoginType(LoginType_value[globalConfig.DefaultLoginType]),
			Version:          version.Version,
		}

		oauth2 := globalConfig.Oauth2
		if !globalConfig.DisableLoginForm {
			resp.LoginType = append(resp.LoginType, &GlobalLoginType{Type: LoginType_normal})
			if !config.GetRuntimeConfig().GetSecurity().ForceEnableMfa {
				resp.LoginType = append(resp.LoginType, &GlobalLoginType{Type: LoginType_email})
			}
		}
		for _, options := range oauth2 {
			resp.LoginType = append(resp.LoginType, &GlobalLoginType{
				Id:        options.Id,
				Type:      LoginType_oauth2,
				Name:      options.Name,
				Icon:      options.Icon,
				AutoLogin: options.AutoLogin,
			})
		}
		return resp, nil
	}
}
