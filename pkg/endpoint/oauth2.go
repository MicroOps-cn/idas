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
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	gohttp "net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MicroOps-cn/fuck/http"
	logs "github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"
	uuid "github.com/satori/go.uuid"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/common"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
	"github.com/MicroOps-cn/idas/pkg/utils/httputil"
	jwtutils "github.com/MicroOps-cn/idas/pkg/utils/jwt"
)

//nolint:revive
const OAuthGrantType_proxy OAuthGrantType = 9

func (r *OAuthTokenRequest) GetRefreshToken() string {
	if r != nil && r.RefreshToken != nil {
		return string(*r.RefreshToken)
	}
	return ""
}

func (m OAuthTokenType) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

type UserToken struct {
	*jwtutils.StandardClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}

func MakeOAuthTokensEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		begin := time.Now()
		req := request.(Requester).GetRequestData().(*OAuthTokenRequest)
		resp := OAuthTokenResponse{TokenType: req.TokenType}
		if restfulReq := request.(RestfulRequester).GetRestfulRequest(); restfulReq == nil {
			err = fmt.Errorf("invalid_grant")
		} else {
			var proxyConfig *models.AppProxyConfig
			app, ok := ctx.Value(global.MetaApp).(*models.App)
			if !ok {
				if req.GrantType == OAuthGrantType_proxy {
					proxyConfig = new(models.AppProxyConfig)
					if !s.VerifyToken(ctx, req.State, models.TokenTypeAppProxyLogin, proxyConfig) {
						return nil, errors.UnauthorizedError()
					}
					if err = s.DeleteToken(ctx, models.TokenTypeAppProxyLogin, req.State); err != nil {
						level.Warn(logs.GetContextLogger(ctx)).Log("msg", "failed to delete token", "err", err)
					}
					app, err = s.GetAppInfo(ctx, opts.WithBasic, opts.WithOAuth2, opts.WithAppId(proxyConfig.AppId))
					if err != nil {
						return err, nil
					} else if app == nil {
						return "", errors.UnauthorizedError()
					}
				} else {
					return "", errors.UnauthorizedError()
				}
			}
			if req.GrantType == OAuthGrantType_refresh_token {
				if username, password, ok := restfulReq.Request.BasicAuth(); ok {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByPassword(ctx, req.GetRefreshToken(), username, password)
				} else if len(req.Username) != 0 && len(req.Password) != 0 {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByPassword(ctx, req.GetRefreshToken(), req.Username, string(req.Password))
				} else {
					resp.AccessToken, resp.RefreshToken, resp.ExpiresIn, err = s.RefreshOAuthTokenByAuthorizationCode(ctx, req.GetRefreshToken(), req.ClientId, string(req.ClientSecret))
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
					session.User = user
					err = s.GetSessionByToken(ctx, req.Code, models.TokenTypeCode, session.User)
					if err != nil {
						return nil, errors.WithServerError(gohttp.StatusUnauthorized, err, "Invalid identity information")
					}
					if err = s.DeleteToken(ctx, models.TokenTypeCode, req.Code); err != nil {
						level.Warn(logs.GetContextLogger(ctx)).Log("msg", "failed to delete token", "err", err)
					}
					if err = s.VerifyUserStatus(ctx, user, false); err != nil {
						return nil, err
					}
					if role, err := s.GetAppRoleByUserId(ctx, app.Id, user.Id); err != nil {
						user.Role = ""
						user.RoleId = ""
					} else {
						user.Role = role.Name
						user.RoleId = role.Id
					}
					session.Proxy = proxyConfig
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
					if proxyConfig.JwtProvider && len(proxyConfig.JwtCookieName) > 0 {
						jwtIssuer := config.Get().GetJwtIssuer()
						if len(proxyConfig.JwtSecret) > 0 {
							secret, err := proxyConfig.GetJwtSecret()
							if err != nil {
								return nil, errors.WithServerError(500, err, "failed to get jwt token")
							}
							jwtSecret, err := base64.StdEncoding.DecodeString(secret)
							if err != nil {
								return nil, errors.WithServerError(500, err, "failed to get jwt token")
							}
							jwtIssuer, err = jwtutils.NewJWTConfigBySecret(string(jwtSecret))
							if err != nil {
								return nil, errors.WithServerError(500, err, "failed to get jwt token")
							}
						}

						token, err := jwtIssuer.SignedString(&jwtutils.StandardClaims{
							Id:        at.Id,
							ExpiresAt: at.Expiry.Unix(),
							IssuedAt:  time.Now().UTC().Unix(),
							NotBefore: time.Now().UTC().Unix(),
							Subject:   user.Username,
						})
						if err != nil {
							return "", errors.NewServerError(500, err.Error())
						}
						resp.Cookies = append(resp.Cookies, (&gohttp.Cookie{
							Name:  proxyConfig.JwtCookieName,
							Path:  "/",
							Value: token,
						}).String())
					}
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
					defer func() {
						if app == nil {
							app = new(models.App)
						}
						if user == nil {
							user = new(models.User)
						}
						status := true
						logContent := fmt.Sprintf("[OAuth2] Authorize login password mode to %s application", app.Name)
						if err != nil {
							status = false
							logContent += ":" + err.Error()
						}
						eventId := logs.GetTraceId(ctx)
						if u, e := uuid.FromString(eventId); e == nil {
							eventId = u.String()
						}
						took := time.Since(begin)
						if e := s.PostEventLog(ctx, eventId, user.Id, user.Username, "", "Authorize", logContent, status, took, logContent); e != nil {
							level.Error(logs.GetContextLogger(ctx)).Log("failed to post event log", "err", e)
						}
					}()
					if app.GrantType&models.AppMeta_password == 0 {
						return nil, errors.NewServerError(500, "Unsupported authorization type.")
					}
					user, err = s.VerifyPassword(ctx, req.Username, string(req.Password), false)
					if user.IsForceMfa() {
						if len(req.Code) > 0 {
							if err = user.VerifyTOTP(req.Code); err != nil {
								return nil, err
							}
						} else {
							resp.NextMethod = []LoginType{LoginType_mfa_totp}
							resp.Error = "TOTP code needs to be provided"
							return resp, nil
						}
					}
				case OAuthGrantType_client_credentials:
					if app.GrantType&models.AppMeta_client_credentials == 0 {
						return nil, errors.NewServerError(500, "Unsupported authorization type.")
					}
					if username, password, ok := restfulReq.Request.BasicAuth(); ok {
						user, err = s.VerifyPassword(ctx, username, password, false)
					}
				}
				if err == nil && user != nil && len(user.Id) > 0 {
					if err = s.VerifyUserStatus(ctx, user, false); err != nil {
						return nil, err
					}
					if role, err := s.GetAppRoleByUserId(ctx, app.Id, user.Id); err != nil {
						user.Role = ""
						user.RoleId = ""
					} else {
						user.RoleId = role.Id
						user.Role = role.Name
					}
					jwtIssuer := config.Get().GetJwtIssuer()

					if app.OAuth2.JwtSignatureKey != nil {
						sigKey, err := app.OAuth2.JwtSignatureKey.UnsafeString()
						if err != nil {
							return nil, errors.NewServerError(500, err.Error())
						}
						jwtIssuer, err = jwtutils.NewJWTIssuer(app.Id, app.OAuth2.JwtSignatureMethod.String(), sigKey)
						if err != nil {
							return nil, errors.NewServerError(500, err.Error())
						}
					}

					user.ExtendedData = nil
					at, err := s.CreateToken(ctx, tokenType, user)
					if err != nil {
						return "", errors.NewServerError(500, err.Error())
					}
					resp.ExpiresIn = int(time.Until(at.Expiry) / time.Minute)
					resp.AccessToken, err = jwtIssuer.SignedString(&jwtutils.StandardClaims{
						Id:        at.Id,
						ExpiresAt: at.Expiry.Unix(),
						IssuedAt:  time.Now().UTC().Unix(),
						NotBefore: time.Now().UTC().Unix(),
						Subject:   user.Username,
					})
					if err != nil {
						return "", errors.NewServerError(500, err.Error())
					}

					if !req.DisableRefreshToken {
						rt, err := s.CreateToken(ctx, models.TokenTypeRefreshToken, user)
						if err != nil {
							return "", errors.NewServerError(500, err.Error())
						}
						resp.RefreshToken, err = jwtIssuer.SignedString(&jwtutils.StandardClaims{
							Id:        rt.Id,
							ExpiresAt: rt.Expiry.Unix(),
							IssuedAt:  time.Now().UTC().Unix(),
							NotBefore: time.Now().UTC().Unix(),
						})
						if err != nil {
							return "", errors.NewServerError(500, err.Error())
						}
					}

					if app.GrantType&models.AppMeta_oidc == models.AppMeta_oidc {
						resp.IdToken, err = jwtIssuer.SignedString(&UserToken{
							StandardClaims: &jwtutils.StandardClaims{
								Id:        w.M(uuid.NewV4()).String(),
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
						})

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

func init() {
	httputil.RegisterTypes(&OAuthAuthorizeRequest_ResponseTypes{}, func(v string) (interface{}, error) {
		rt := OAuthAuthorizeRequest_ResponseTypes{}
		for _, name := range strings.Split(string(v), " ") {
			if len(name) > 0 {
				val, ok := OAuthAuthorizeRequest_ResponseType_value[name]
				if ok {
					rt.Types = append(rt.Types, OAuthAuthorizeRequest_ResponseType(val))
				}
			}
		}
		if len(rt.Types) == 0 {
			rt.Types = append(rt.Types, OAuthAuthorizeRequest_none)
		}
		return rt, nil
	})
}

func MakeOAuthAuthorizeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		begin := time.Now()
		logger := logs.GetContextLogger(ctx)
		req := request.(Requester).GetRequestData().(*OAuthAuthorizeRequest)
		resp := SimpleResponseWrapper[interface{}]{}

		stdReq := request.(RestfulRequester).GetRestfulRequest()
		stdResp := request.(RestfulRequester).GetRestfulResponse()
		var app *models.App
		var user *models.User
		defer func() {
			if app == nil {
				app = new(models.App)
			}
			if user == nil {
				user = new(models.User)
			}
			status := true
			if resp.Error != nil || err != nil {
				status = false
			}
			eventId := logs.GetTraceId(ctx)
			if u, e := uuid.FromString(eventId); e == nil {
				eventId = u.String()
			}
			took := time.Since(begin)
			logContent := fmt.Sprintf("[OAuth] Authorize login to %s application", app.Name)
			if e := s.PostEventLog(ctx, eventId, user.Id, user.Username, "", "", "", status, took, logContent); e != nil {
				level.Error(logs.GetContextLogger(ctx)).Log("failed to post event log", "err", e)
			}
		}()

		if len(req.ClientId) == 0 {
			return nil, errors.ParameterError("client_id")
		}
		var ok bool
		if user, ok = ctx.Value(global.MetaUser).(*models.User); !ok || user == nil {
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
		query := uri.Query()
		if req.AccessType == "proxy" {
			var proxyConfig models.AppProxyConfig
			if !s.VerifyToken(ctx, req.State, models.TokenTypeAppProxyLogin, &proxyConfig, req.ClientId) {
				return nil, errors.UnauthorizedError()
			}
			app, err = s.GetAppInfo(ctx, opts.WithBasic, opts.WithAppId(proxyConfig.AppId), opts.WithUsers(user.Id))
			if err != nil {
				return err, nil
			} else if app == nil {
				gohttp.Redirect(stdResp.ResponseWriter, stdReq.Request, w.M(common.GetWebURL(ctx, common.WithSubPages("404"))), gohttp.StatusFound)
				return nil, nil
			}
		} else {
			appKey, err := s.GetAppKeyFromKey(ctx, req.ClientId)
			if err != nil {
				level.Error(logger).Log("msg", "failed to get appId from client_id")
				gohttp.Redirect(stdResp.ResponseWriter, stdReq.Request, w.M(common.GetWebURL(ctx, common.WithSubPages("404"))), gohttp.StatusFound)
				return nil, nil
			}
			app, err = s.GetAppInfo(ctx, opts.WithBasic, opts.WithOAuth2, opts.WithAppId(appKey.AppId), opts.WithUsers(user.Id))
			if err != nil {
				return err, nil
			} else if app == nil {
				gohttp.Redirect(stdResp.ResponseWriter, stdReq.Request, w.M(common.GetWebURL(ctx, common.WithSubPages("404"))), gohttp.StatusFound)
				return nil, nil
			}
		}

		//if code, err = s.GetAuthCodeByAppId(ctx, appKey.AppId, user, sessionId); err != nil && !errors.IsNotFount(err) {
		//	return nil, err
		//}
		allowAccess := true
		if len(app.OAuth2.AuthorizedRedirectUrl) > 0 {
			if !w.Has(app.OAuth2.AuthorizedRedirectUrl, req.RedirectUri, func(a, b string) bool {
				if len(a) > 0 && len(b) > 0 {
					return strings.HasPrefix(b, a)
				}
				return false
			}) {
				allowAccess = false
			}
		} else {
			if !strings.HasPrefix(req.RedirectUri, app.Url) {
				allowAccess = false
			}
		}

		if app.GrantType&models.AppMeta_authorization_code == 0 {
			return nil, errors.NewServerError(500, "Unsupported authorization type.")
		}
		if !allowAccess || len(app.Users) == 0 {
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

			for _, rType := range req.ResponseType.Types {
				switch rType {
				case OAuthAuthorizeRequest_code, OAuthAuthorizeRequest_none:
					query.Add("code", token.Id)
					query.Add("state", req.State)
					if len(req.CodeChallenge) > 0 {
						codeVerifier := req.CodeChallenge
						if req.CodeChallengeMethod == OAuthAuthorizeRequest_S256 {
							hash := sha256.Sum256([]byte(req.CodeChallenge))
							codeVerifier = base64.URLEncoding.EncodeToString(hash[:])
						}
						query.Add("code_verifier", codeVerifier)
					}
				case OAuthAuthorizeRequest_token:
					accessToken, refreshToken, expiresIn, err := s.GetOAuthTokenByAuthorizationCode(ctx, token.Id, req.ClientId)
					if err != nil {
						return nil, err
					}
					query.Add("access_token", accessToken)
					query.Add("refresh_token", refreshToken)
					query.Add("expires_in", strconv.Itoa(expiresIn))
				case OAuthAuthorizeRequest_id_token:
					jwtIssuer := config.Get().GetJwtIssuer()
					idToken, err := jwtIssuer.SignedString(&UserToken{
						StandardClaims: &jwtutils.StandardClaims{
							Id:        w.M(uuid.NewV4()).String(),
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
					})
					if err != nil {
						return nil, err
					}
					query.Add("id_token", idToken)
				}
			}
			uri.RawQuery = query.Encode()
			stdResp.AddHeader("Location", uri.String())
			stdResp.WriteHeader(302)
		}
		return &resp, nil
	}
}

func MakeWellknownOpenidConfigurationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		logger := logs.GetContextLogger(ctx)
		req := request.(Requester).GetRequestData().(*OIDCWellKnownRequest)
		app, err := s.GetAppInfo(ctx, opts.WithBasic, opts.WithOAuth2, opts.WithAppId(req.ClientId))
		if err != nil {
			level.Error(logger).Log("msg", "failed to get app info", "err", err)
			return nil, err
		}
		return &OIDCWellKnownResponse{
			Issuer:                w.M(common.GetURL(ctx, common.WithRoot)),
			JwksUri:               w.M(common.GetURL(ctx, common.WithAPI("v1", "oauth", "jwks"), common.WithParam("client_id", req.ClientId))),
			TokenEndpoint:         w.M(common.GetURL(ctx, common.WithAPI("v1", "oauth", "token"))),
			AuthorizationEndpoint: w.M(common.GetURL(ctx, common.WithAPI("v1", "oauth", "authorize"))),
			UserinfoEndpoint:      w.M(common.GetURL(ctx, common.WithAPI("v1", "oauth", "userinfo"))),
			RevocationEndpoint:    w.M(common.GetURL(ctx, common.WithAPI("v1", "oauth", "revoke"))),
			CodeChallengeMethodsSupported: []string{
				"plain",
				"S256",
			},
			GrantTypesSupported: append([]string{"refresh_token", "urn:ietf:params:oauth:grant-type:jwt-bearer"}, app.GrantType.Name()...),
			ResponseTypesSupported: []string{
				"code",
				"token",
				"id_token",
				"code token",
				"code id_token",
				"token id_token",
				"code token id_token",
				"none",
			},
			SubjectTypesSupported: []string{
				"public",
			},
			ScopesSupported: []string{
				"openid",
				"profile",
				"email",
			},
			IdTokenSigningAlgValuesSupported: []string{
				"RS256",
				"HS256",
			},
			TokenEndpointAuthMethodsSupported: []string{
				"client_secret_basic",
				"client_secret_post",
			},
			ClaimsSupported: []string{
				"aud",
				"exp",
				"jti",
				"iat",
				"iss",
				"nbf",
				"sub",
				"username",
				"email",
				"fullName",
			},
		}, nil
	}
}

func MakeOAuthJWKSEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		jwtIssuer := config.Get().GetJwtIssuer()
		pk := jwtIssuer.GetPublicKey()
		var n, e []byte
		var pkSize int
		switch pubKey := pk.(type) {
		case *rsa.PublicKey:
			e = make([]byte, 4)
			binary.BigEndian.PutUint32(e, uint32(pubKey.E))
			pkSize = pubKey.Size()
			n = pubKey.N.Bytes()
			return &OAuthJWKSResponse{
				Keys: []*OAuthJWKSResponse_Key{
					{
						Kty: "RSA",
						Alg: "RS" + strconv.Itoa(pkSize),
						Use: "sig",
						N:   base64.URLEncoding.EncodeToString(n),
						E:   base64.URLEncoding.EncodeToString(e[1:]),
					},
				},
			}, nil

		case *ecdsa.PublicKey:
			switch pubKey.Curve.Params().Name {
			case "P-256":
				pkSize = 256
			case "P-384":
				pkSize = 384
			case "P-521":
				pkSize = 521
			}
			switch pubKey.Curve.Params().Name {
			case "P-256":
				return &OAuthJWKSResponse{
					Keys: []*OAuthJWKSResponse_Key{
						{
							Kty: "EC",
							Use: "sig",
							Alg: "ES" + strconv.Itoa(pkSize),
							Crv: pubKey.Curve.Params().Name,
							X:   base64.URLEncoding.EncodeToString(pubKey.X.Bytes()),
							Y:   base64.URLEncoding.EncodeToString(pubKey.Y.Bytes()),
						},
					},
				}, nil
			}
		}
		return nil, nil
	}
}
