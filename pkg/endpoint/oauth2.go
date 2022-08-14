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
	"net/url"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/logs"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

func MakeOAuthTokensEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*OAuthTokenRequest)
		resp := OAuthTokenResponse{TokenType: "Bearer"}
		if restfulReq := request.(RestfulRequester).GetRestfulRequest(); restfulReq == nil {
			err = fmt.Errorf("invalid_grant")
		} else {
			switch req.GrantType {
			case OAuthGrantType_authorization_code:
				resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.GetOAuthTokenByAuthorizationCode(ctx, req.Code, req.ClientId)
			case OAuthGrantType_password:
				resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.GetOAuthTokenByPassword(ctx, req.Username, req.Password)
			case OAuthGrantType_client_credentials:
				if username, password, ok := restfulReq.Request.BasicAuth(); ok {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.GetOAuthTokenByPassword(ctx, username, password)
				} else {
					err = fmt.Errorf("invalid_request")
				}
			case OAuthGrantType_refresh_token:
				if username, password, ok := restfulReq.Request.BasicAuth(); ok {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByPassword(ctx, req.RefreshToken, username, password)
				} else if len(req.Username) != 0 && len(req.Password) != 0 {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByPassword(ctx, req.RefreshToken, req.Username, req.Password)
				} else {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByAuthorizationCode(ctx, req.RefreshToken, req.ClientId, req.ClientSecret)
				}
			default:
				err = fmt.Errorf("unsupported_grant_type")
			}
		}

		if err != nil {
			resp.Error = err.Error()
			if restfulResp := request.(RestfulRequester).GetRestfulResponse(); restfulResp != nil {
				restfulResp.WriteHeader(400)
			}
		}
		return &resp, nil
	}
}

func MakeOAuthAuthorizeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		logger := logs.GetContextLogger(ctx)
		req := request.(Requester).GetRequestData().(*OAuthAuthorizeRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		var code string

		stdResp := request.(RestfulRequester).GetRestfulResponse()

		if len(req.ClientId) == 0 {
			return nil, errors.ParameterError("client_id")
		}
		users, ok := ctx.Value(global.MetaUser).([]*models.User)
		if !ok || len(users) == 0 {
			level.Warn(logger).Log("msg", "failed to get user from context")
			resp.Error = errors.NotLoginError
			return resp, nil
		}
		sessionId, ok := ctx.Value(global.LoginSession).(string)
		if !ok || len(sessionId) == 0 {
			level.Warn(logger).Log("msg", "failed to get session from context")
			resp.Error = errors.NotLoginError
			return resp, nil
		}
		uri, err := url.Parse(req.RedirectUri)
		if err != nil {
			return nil, errors.ParameterError("redirect_uri")
		}
		query := uri.Query()
		for _, user := range users {
			if code, err = s.GetAuthCodeByClientId(ctx, req.ClientId, user, sessionId, user.Storage); errors.IsNotFount(err) {
				continue
			} else if err != nil {
				return nil, err
			}
			break
		}
		if code == "" {
			return nil, errors.StatusNotFound("Authorize")
		}
		switch req.ResponseType {
		case OAuthAuthorizeRequest_code, OAuthAuthorizeRequest_default:
			query.Add("code", code)
			query.Add("state", req.State)
			uri.RawQuery = query.Encode()
			stdResp.AddHeader("Location", uri.String())
			stdResp.WriteHeader(302)
		case OAuthAuthorizeRequest_token:
			accessToken, refreshToken, expiresIn, err := s.GetOAuthTokenByAuthorizationCode(ctx, code, req.ClientId)
			if err != nil {
				return nil, err
			}
			query.Add("access_token", accessToken)
			query.Add("refresh_token", refreshToken)
			query.Add("expires_in", strconv.Itoa(expiresIn))
			uri.RawQuery = query.Encode()
			stdResp.AddHeader("Location", uri.String())
			stdResp.WriteHeader(302)
		}
		return &resp, nil
	}
}
