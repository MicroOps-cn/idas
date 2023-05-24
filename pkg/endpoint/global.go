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
	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

func MakeGetLoginTypeEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var resp []GlobalLoginType
		oauth2 := config.Get().GetGlobal().Oauth2
		if !config.Get().GetGlobal().DisableLoginForm {
			resp = append(resp, GlobalLoginType{Type: LoginType_normal}, GlobalLoginType{Type: LoginType_email})
		}
		for _, options := range oauth2 {
			resp = append(resp, GlobalLoginType{
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
