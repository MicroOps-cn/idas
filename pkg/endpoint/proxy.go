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
	"io"
	gohttp "net/http"
	"sort"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/go-kit/kit/endpoint"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

type ProxyResponse struct {
	Header gohttp.Header `json:"-"`
	Body   io.ReadCloser `json:"-"`
	Code   int           `json:"code"`
	Error  error         `json:"error"`
}

type ProxyConfig struct {
	Token                 string
	ExternalURL           string
	ClientId              string
	RedirectURICookieName string
}

func MakeGetProxyConfigEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r, ok := request.(*restful.Request)
		if !ok {
			return nil, errors.NewServerError(500, "system error: request type exception")
		}
		host, _, _ := strings.Cut(r.Request.Host, ":")
		proxyConfig, err := s.GetProxyConfig(ctx, host)
		if err != nil {
			return nil, err
		}
		token, err := s.CreateToken(ctx, models.TokenTypeAppProxyLogin, proxyConfig)
		if err != nil {
			return nil, err
		}

		return &ProxyConfig{Token: token.Id, ClientId: proxyConfig.AppId}, nil
	}
}

func MakeProxyRequestEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp := &ProxyResponse{Code: 500, Error: errors.NewServerError(500, "system error")}
		r, ok := request.(*gohttp.Request)
		if !ok {
			return resp, nil
		}
		user, ok := ctx.Value(global.MetaUser).(*models.User)
		if !ok || user == nil {
			resp.Error = errors.NewServerError(401, "system error: no authorization")
			resp.Code = 401
			return resp, nil
		}
		var proxyConfig *models.AppProxyConfig

		proxyConfig, ok = ctx.Value(global.MetaProxyConfig).(*models.AppProxyConfig)
		if !ok {
			resp.Error = errors.NewServerError(403, "system error: forbidden")
			resp.Code = 403
			return resp, nil
		}
		sort.Sort(proxyConfig.Urls)
		var proxyURLConfig *models.AppProxyUrl
		for _, proxyURL := range proxyConfig.Urls {
			if strings.HasPrefix(r.URL.Path, proxyURL.Url) {
				proxyURLConfig = proxyURL
				break
			}
		}
		if proxyURLConfig == nil {
			return nil, errors.StatusNotFound(r.URL.Path)
		}
		if len(proxyURLConfig.Upstream) == 0 {
			proxyURLConfig.Upstream = proxyConfig.Upstream
		}
		if len(proxyConfig.URLRoles) > 0 {
			var roleMatched bool
			for _, role := range proxyConfig.URLRoles {
				if role.AppProxyURLId == proxyURLConfig.Id && role.AppRoleId == user.RoleId {
					roleMatched = true
					break
				}
			}
			if !roleMatched {
				return nil, errors.StatusForbidden(r.URL.Path)
			}
		}

		oriResp, err := s.SendProxyRequest(ctx, r, proxyConfig, proxyURLConfig)
		if err != nil {
			resp.Error = err
			resp.Code = 500
			return resp, nil
		}
		resp.Header = oriResp.Header.Clone()
		if proxyConfig.HstsOffload {
			resp.Header.Del("Strict-Transport-Security")
		}
		resp.Body = oriResp.Body
		resp.Code = oriResp.StatusCode
		return resp, nil
	}
}
