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
	"time"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
)

//nolint:revive
const OAuthGrantType_proxy OAuthGrantType = 9

func MakeOAuthTokensEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*OAuthTokenRequest)
		resp := OAuthTokenResponse{TokenType: req.TokenType}
		if restfulReq := request.(RestfulRequester).GetRestfulRequest(); restfulReq == nil {
			err = fmt.Errorf("invalid_grant")
		} else {
			if req.GrantType == OAuthGrantType_refresh_token {
				if username, password, ok := restfulReq.Request.BasicAuth(); ok {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByPassword(ctx, req.RefreshToken, username, password)
				} else if len(req.Username) != 0 && len(req.Password) != 0 {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByPassword(ctx, req.RefreshToken, req.Username, req.Password)
				} else {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByAuthorizationCode(ctx, req.RefreshToken, req.ClientId, req.ClientSecret)
				}
			} else {
				tokenType := models.TokenTypeToken
				if req.TokenType == OAuthTokenType_Cookie {
					req.DisableRefreshToken = true
					tokenType = models.TokenTypeLoginSession
				}
				user := new(models.User)
				switch req.GrantType {
				case OAuthGrantType_proxy:
					var session models.ProxySession
					var proxyConfig []*models.AppProxyConfig
					err = s.GetSessionByToken(ctx, req.Code, models.TokenTypeCode, &session.User)
					if err != nil {
						return nil, errors.UnauthorizedError()
					}
					if err = s.DeleteToken(ctx, models.TokenTypeCode, req.Code); err != nil {
						level.Warn(logs.GetContextLogger(ctx)).Log("msg", "failed to delete token", "err", err)
					}
					if !s.VerifyToken(ctx, req.State, models.TokenTypeAppProxyLogin, &proxyConfig, req.ClientId) {
						return nil, errors.UnauthorizedError()
					}
					if len(proxyConfig) > 1 {
						return nil, errors.NewServerError(500, "proxy configuration exception")
					} else if len(proxyConfig) == 0 {
						return nil, errors.UnauthorizedError()
					}
					session.Proxy = proxyConfig[0]
					at, err := s.CreateToken(ctx, tokenType, &session)
					if err != nil {
						return "", errors.NewServerError(500, err.Error())
					}
					resp.AccessToken = at.Id
					if !req.DisableRefreshToken {
						rt, err := s.CreateToken(ctx, models.TokenTypeRefreshToken, session)
						if err != nil {
							return "", errors.NewServerError(500, err.Error())
						}
						resp.RefreshToken = rt.Id
					}
					resp.ExpiresIn = int(global.TokenExpiration / time.Minute)
					return &resp, nil
				case OAuthGrantType_authorization_code:
					err = s.GetSessionByToken(ctx, req.Code, models.TokenTypeCode, user)
				case OAuthGrantType_password:
					user, err = s.VerifyPassword(ctx, req.Username, req.Password)
				case OAuthGrantType_client_credentials:
					if username, password, ok := restfulReq.Request.BasicAuth(); ok {
						user, err = s.VerifyPassword(ctx, username, password)
					}
				}
				if err == nil && user != nil && len(user.Id) > 0 {
					at, err := s.CreateToken(ctx, tokenType, user)
					if err != nil {
						return "", errors.NewServerError(500, err.Error())
					}
					resp.AccessToken = at.Id
					if !req.DisableRefreshToken {
						rt, err := s.CreateToken(ctx, models.TokenTypeRefreshToken, user)
						if err != nil {
							return "", errors.NewServerError(500, err.Error())
						}
						resp.RefreshToken = rt.Id
					}
					resp.ExpiresIn = int(global.TokenExpiration / time.Minute)
					return &resp, nil
				}
				return nil, errors.UnauthorizedError()
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
		user, ok := ctx.Value(global.MetaUser).(*models.User)
		if !ok || user == nil {
			level.Warn(logger).Log("msg", "failed to get user from context")
			resp.Error = errors.NotLoginError()
			return resp, nil
		}
		sessionId, ok := ctx.Value(global.LoginSession).(string)
		if !ok || len(sessionId) == 0 {
			level.Warn(logger).Log("msg", "failed to get session from context")
			resp.Error = errors.NotLoginError()
			return resp, nil
		}
		uri, err := url.Parse(req.RedirectUri)
		if err != nil {
			return nil, errors.ParameterError("redirect_uri")
		}
		appKey, err := s.GetAppKeyFromKey(ctx, req.ClientId)
		if err != nil {
			level.Error(logger).Log("msg", "failed to get appId from client_id")
			return nil, errors.ParameterError("client_id")
		}

		query := uri.Query()
		app, err := s.GetAppInfo(ctx, opts.WithAppId(appKey.AppId), opts.WithBasic)
		if err != nil {
			return err, nil
		} else if app == nil {
			return nil, errors.StatusNotFound("Authorize")
		}
		if code, err = s.GetAuthCodeByAppId(ctx, appKey.AppId, user, sessionId); err != nil {
			return nil, err
		}
		if code == "" {
			httpExternalURL, ok := ctx.Value(global.HTTPExternalURLKey).(string)
			if !ok {
				return nil, errors.StatusNotFound("Authorize")
			}
			stdResp.AddHeader("Location", httpExternalURL+"403")
			stdResp.WriteHeader(302)
		} else {
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
		}
		return &resp, nil
	}
}
