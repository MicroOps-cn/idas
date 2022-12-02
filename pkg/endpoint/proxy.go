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

package endpoint

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/errors"
)

type ProxyResponse struct {
	Header http.Header
	Body   io.ReadCloser
	Code   int
	Error  error
}

func MakeProxyRequestEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp := &ProxyResponse{Code: 500, Error: fmt.Errorf("system error")}
		r, ok := request.(http.Request)
		if !ok {
			return resp, nil
		}
		users, ok := ctx.Value(global.MetaUser).([]*models.User)
		if !ok || len(users) == 0 {
			resp.Error = fmt.Errorf("system error: no authorization")
			resp.Code = 401
			return resp, nil
		}
		err := errors.NewMultipleError()
		var (
			proxyConfig *models.AppProxyConfig
			e           error
		)
		for _, user := range users {
			if proxyConfig, e = s.GetProxyConfig(ctx, user, r.Host, r.Method, r.URL.EscapedPath()); err != nil {
				_ = err.Append(e)
			} else if proxyConfig != nil {
				break
			}
		}
		if err.HasError() {
			resp.Error = err
			return resp, nil
		} else if proxyConfig == nil {
			resp.Error = fmt.Errorf("not found")
			resp.Code = 404
			return resp, nil
		}

		return &ProxyResponse{Code: 200}, err
	}
}
