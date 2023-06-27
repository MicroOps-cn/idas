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
	gohttp "net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/MicroOps-cn/fuck/http"
	logs "github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/common"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
)

//nolint:revive
const OAuthGrantType_proxy OAuthGrantType = 9

func (r *OAuthTokenRequest) GetRefreshToken() string {
	if r != nil && r.RefreshToken != nil {
		return string(*r.RefreshToken)
	}
	return ""
}

type UserToken struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}

func MakeOAuthTokensEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*OAuthTokenRequest)
		resp := OAuthTokenResponse{TokenType: req.TokenType}
		if restfulReq := request.(RestfulRequester).GetRestfulRequest(); restfulReq == nil {
			err = fmt.Errorf("invalid_grant")
		} else {
			app, ok := ctx.Value(global.MetaApp).(*models.App)
			if !ok {
				return "", errors.UnauthorizedError()
			}
			if req.GrantType == OAuthGrantType_refresh_token {
				if username, password, ok := restfulReq.Request.BasicAuth(); ok {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByPassword(ctx, req.GetRefreshToken(), username, password)
				} else if len(req.Username) != 0 && len(req.Password) != 0 {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByPassword(ctx, req.GetRefreshToken(), req.Username, req.Password)
				} else {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByAuthorizationCode(ctx, req.GetRefreshToken(), req.ClientId, req.ClientSecret)
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
					if app.GrantType&models.AppMeta_proxy == 0 {
						return nil, errors.NewServerError(500, "Unsupported authorization type.")
					}
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
					if app.GrantType&models.AppMeta_authorization_code == 0 {
						return nil, errors.NewServerError(500, "Unsupported authorization type.")
					}
					err = s.GetSessionByToken(ctx, req.Code, models.TokenTypeCode, user)
					if err == nil {
						_ = s.DeleteToken(ctx, models.TokenTypeCode, req.Code)
					}
				case OAuthGrantType_password:
					if app.GrantType&models.AppMeta_password == 0 {
						return nil, errors.NewServerError(500, "Unsupported authorization type.")
					}
					user, err = s.VerifyPassword(ctx, req.Username, req.Password, false)
				case OAuthGrantType_client_credentials:
					if app.GrantType&models.AppMeta_client_credentials == 0 {
						return nil, errors.NewServerError(500, "Unsupported authorization type.")
					}
					if username, password, ok := restfulReq.Request.BasicAuth(); ok {
						user, err = s.VerifyPassword(ctx, username, password, false)
					}
				}
				if err == nil && user != nil && len(user.Id) > 0 {
					if role, err := s.GetAppRoleByUserId(ctx, app.Id, user.Id); err != nil {
						user.Role = ""
					} else {
						user.Role = role.Name
					}
					jwtSecret := config.Get().Global.GetJwtSecret()

					at, err := s.CreateToken(ctx, tokenType, user)
					if err != nil {
						return "", errors.NewServerError(500, err.Error())
					}
					resp.ExpiresIn = int(time.Until(at.Expiry) / time.Minute)
					resp.AccessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
						Id:        at.Id,
						ExpiresAt: at.Expiry.Unix(),
						IssuedAt:  time.Now().UTC().Unix(),
						NotBefore: time.Now().UTC().Unix(),
						Subject:   user.Username,
					}).SignedString([]byte(jwtSecret))
					if err != nil {
						return "", errors.NewServerError(500, err.Error())
					}

					if !req.DisableRefreshToken {
						rt, err := s.CreateToken(ctx, models.TokenTypeRefreshToken, user)
						if err != nil {
							return "", errors.NewServerError(500, err.Error())
						}
						resp.RefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
							Id:        rt.Id,
							ExpiresAt: rt.Expiry.Unix(),
							IssuedAt:  time.Now().UTC().Unix(),
							NotBefore: time.Now().UTC().Unix(),
						}).SignedString([]byte(jwtSecret))
						if err != nil {
							return "", errors.NewServerError(500, err.Error())
						}
					}

					if app.GrantType&models.AppMeta_oidc == models.AppMeta_oidc {
						resp.IdToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, &UserToken{
							StandardClaims: jwt.StandardClaims{
								Id:        uuid.NewV4().String(),
								Audience:  app.Name,
								ExpiresAt: time.Now().UTC().Add(time.Minute * 10).Unix(),
								Issuer:    config.Get().GetGlobal().GetAppName(),
								IssuedAt:  time.Now().UTC().Unix(),
								NotBefore: time.Now().UTC().Unix(),
								Subject:   user.Username,
							},
							Username: user.Username,
							Email:    user.Email,
							FullName: user.FullName,
						}).SignedString([]byte(jwtSecret))
						if err != nil {
							return "", errors.NewServerError(500, err.Error())
						}
					}

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

		stdReq := request.(RestfulRequester).GetRestfulRequest()
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
			gohttp.Redirect(stdResp.ResponseWriter, stdReq.Request, w.M(common.GetWebURL(ctx, common.WithSubPages("404"))), gohttp.StatusFound)
			return nil, nil
		}

		query := uri.Query()
		app, err := s.GetAppInfo(ctx, opts.WithBasic, opts.WithAppId(appKey.AppId), opts.WithUsers(user.Id))
		if err != nil {
			return err, nil
		} else if app == nil {
			gohttp.Redirect(stdResp.ResponseWriter, stdReq.Request, w.M(common.GetWebURL(ctx, common.WithSubPages("404"))), gohttp.StatusFound)
			return nil, nil
		}

		//if code, err = s.GetAuthCodeByAppId(ctx, appKey.AppId, user, sessionId); err != nil && !errors.IsNotFount(err) {
		//	return nil, err
		//}

		if app.GrantType&models.AppMeta_authorization_code == 0 {
			return nil, errors.NewServerError(500, "Unsupported authorization type.")
		}
		if len(app.Users) == 0 {
			httpExternalURL, ok := ctx.Value(global.HTTPExternalURLKey).(string)
			if !ok {
				return nil, errors.StatusNotFound("Authorize")
			}
			webPrefix, ok := ctx.Value(global.HTTPWebPrefixKey).(string)
			if !ok {
				return nil, errors.StatusNotFound("Authorize")
			}
			u, err := url.Parse(httpExternalURL)
			if err != nil {
				return nil, errors.StatusNotFound("httpExternalURL")
			}
			u.Path = http.JoinPath(u.Path, webPrefix, "403")
			stdResp.AddHeader("Location", u.String())
			stdResp.WriteHeader(302)
		} else {
			user.Role = app.Users[0].Role
			token, err := s.CreateToken(ctx, models.TokenTypeCode, user)
			if err != nil {
				return "", err
			}
			appURL, err := url.Parse(app.Url)
			if err == nil && appURL.User != nil {
				uri.User = appURL.User
			}
			switch req.ResponseType {
			case OAuthAuthorizeRequest_code, OAuthAuthorizeRequest_default:
				query.Add("code", token.Id)
				query.Add("state", req.State)
				uri.RawQuery = query.Encode()
				stdResp.AddHeader("Location", uri.String())
				stdResp.WriteHeader(302)
			case OAuthAuthorizeRequest_token:
				accessToken, refreshToken, expiresIn, err := s.GetOAuthTokenByAuthorizationCode(ctx, token.Id, req.ClientId)
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
