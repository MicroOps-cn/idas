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
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/fuck/sets"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/xlzd/gotp"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/client/oauth2"
	"github.com/MicroOps-cn/idas/pkg/common"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
)

type LoginCode struct {
	UserId string    `json:"userId"`
	Code   string    `json:"code"`
	Type   LoginType `json:"type"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func MakeSendLoginCaptchaEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		begin := time.Now()
		req := request.(Requester).GetRequestData().(*SendLoginCaptchaRequest)
		resp := SimpleResponseWrapper[*SendLoginCaptchaResponseData]{}
		var user *models.User
		defer func() {
			var userId, username string
			if user != nil {
				userId = user.Id
				username = user.Username
			}
			eventId, message, status, took := GetEventMeta(ctx, "SendLoginCaptcha", begin, err, resp)
			if e := s.PostEventLog(ctx, eventId, userId, username, "", "SendLoginCaptcha", message, status, took); e != nil {
				level.Error(logs.GetContextLogger(ctx)).Log("failed to post event log", "err", e)
			}
		}()
		switch req.Type {
		case LoginType_mfa_email, LoginType_email:
			switch req.Type {
			case LoginType_mfa_email:
				email, username := req.GetEmail(), req.GetUsername()
				if user, err = s.GetUserInfoByUsernameAndEmail(ctx, username, email); err != nil {
					level.Warn(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to get user", "username", username, "email", email)
					resp.Error = errors.ParameterError("email")
					return resp, nil
				}
			case LoginType_email:
				if user, err = s.GetUser(ctx, opts.WithEmail(req.GetEmail()), opts.WithUserExt); err != nil {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
					level.Warn(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to get user", "username", req.GetEmail())
					return resp, nil
				}
			}
			if user.Status.Is(models.UserMeta_normal) {
				loginCode := LoginCode{UserId: user.Id, Type: req.Type, Code: strings.ToUpper(uuid.NewV4().String()[:6])}
				token, err := s.CreateToken(ctx, models.TokenTypeLoginCode, &loginCode)
				if err != nil {
					return nil, errors.NewServerError(http.StatusInternalServerError, "Failed to create token")
				}
				to := fmt.Sprintf("%s<%s>", user.FullName, user.Email)
				data := map[string]interface{}{
					"user":       user,
					"token":      token,
					"code":       loginCode.Code,
					"userId":     user.Id,
					"siteTitle":  config.Get().GetGlobal().GetTitle(),
					"adminEmail": config.Get().GetGlobal().GetAdminEmail(),
				}
				if err = s.SendEmail(ctx, data, "User:SendLoginCaptcha", to); err != nil {
					return nil, err
				}
				resp.Data = &SendLoginCaptchaResponseData{Token: token.Id}
				return &resp, nil
			}
		default:
			return nil, errors.NewServerError(http.StatusBadRequest, "Unsupported authentication method: "+req.Type.String())
		}

		return nil, errors.StatusNotFound("user")
	}
}

func getMFAMethod(user *models.User) sets.Set[LoginType] {
	method := sets.New[LoginType]()
	userExt := user.ExtendedData
	if userExt != nil {
		if userExt.EmailAsMFA {
			method.Insert(LoginType_mfa_email)
		}
		if userExt.TOTPAsMFA {
			method.Insert(LoginType_mfa_totp)
		}
		if userExt.SmsAsMFA {
			method.Insert(LoginType_mfa_sms)
		}
	}
	return method
}

func MakeUserOAuthLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		stdResp := request.(RestfulRequester).GetRestfulResponse()
		stdReq := request.(RestfulRequester).GetRestfulRequest()
		req := request.(Requester).GetRequestData().(*OAuthLoginRequest)
		oAuthOptions := config.Get().GetOAuthOptions(req.Id)
		if oAuthOptions != nil {
			if req.Code == "" && req.State == "" {
				redirectURL, err := oAuthOptions.GetRedirectURL(ctx, stdReq.QueryParameter("redirect_uri"))
				if err != nil {
					return nil, err
				}
				authURL, err := url.Parse(oAuthOptions.AuthUrl)
				if err != nil {
					return nil, err
				}
				q := authURL.Query()
				q.Set("response_type", "code")
				if scope := oAuthOptions.GetScope(); len(scope) != 0 {
					q.Set("scope", scope)
				}

				q.Set("redirect_uri", redirectURL.String())
				q.Set("client_id", oAuthOptions.ClientId)
				token, err := s.CreateToken(ctx, models.TokenTypeOAuthState, "")
				if err != nil {
					return nil, err
				}
				q.Set("state", token.Id)
				authURL.RawQuery = q.Encode()
				http.Redirect(stdResp.ResponseWriter, stdReq.Request, authURL.String(), 302)
			} else if req.Code != "" && req.State != "" {
				if s.VerifyToken(ctx, req.State, models.TokenTypeOAuthState, nil) {
					user, err := oauth2.NewClient(oAuthOptions).GetUserInfo(ctx, req.Code, stdReq.QueryParameter("redirect_uri"))
					if err != nil {
						return nil, err
					}
					var userInfo *models.User
				loop:
					for _, loginIdName := range strings.Split(oAuthOptions.LoginId, ",") {
						switch loginIdName {
						case "username":
							if len(user.Username) > 0 {
								userInfo, err = s.GetUser(ctx, opts.WithUsername(user.Username))
								if err != nil {
									level.Info(logs.GetContextLogger(ctx)).Log("msg", "failed to get user from username", "username", user.Username, "err", err)
								} else if len(userInfo.Id) > 0 {
									break loop
								}
							}
						case "email":
							if len(user.Email) > 0 {
								userInfo, err = s.GetUser(ctx, opts.WithEmail(user.Email))
								if err != nil {
									level.Info(logs.GetContextLogger(ctx)).Log("msg", "failed to get user from email", "email", user.Email, "err", err)
								} else if len(userInfo.Id) > 0 {
									break loop
								}
							}
						case "phoneNumber":
							if len(user.PhoneNumber) > 0 {
								userInfo, err = s.GetUser(ctx, opts.WithPhoneNumber(user.PhoneNumber))
								if err != nil {
									level.Info(logs.GetContextLogger(ctx)).Log("msg", "failed to get user from phoneNumber", "phoneNumber", user.PhoneNumber, "err", err)
								} else if len(userInfo.Id) > 0 {
									break loop
								}
							}
						}
					}
					if userInfo == nil || userInfo.Id == "" {
						level.Warn(logs.GetContextLogger(ctx)).Log("msg", "user is not exists", "user", user.String())
						http.Redirect(stdResp.ResponseWriter, stdReq.Request, w.M(common.GetWebURL(ctx, common.WithSubPages("403"))), http.StatusFound)
						return nil, nil
					}
					app, err := s.GetAppInfo(ctx, opts.WithBasic, opts.WithUsers(userInfo.Id), opts.WithAppName(config.Get().GetGlobal().GetAppName()))
					if err != nil && !errors.IsNotFount(err) {
						level.Error(logs.GetContextLogger(ctx)).Log("msg", "failed to get app info", "err", err)
					} else if app != nil {
						role, err := s.GetAppRoleByUserId(ctx, app.Id, userInfo.Id)
						if err == nil {
							userInfo.RoleId = role.Id
							userInfo.Role = role.Name
						} else if !errors.IsNotFount(err) {
							level.Error(logs.GetContextLogger(ctx)).Log("msg", "failed to get app role", "err", err)
						}
					}
					if len(userInfo.Role) == 0 {
						userInfo.Role = user.Role
					}
					if len(userInfo.Avatar) == 0 {
						userInfo.Avatar = user.Avatar
						_ = s.PatchUser(ctx, map[string]interface{}{"id": userInfo.Id, "avatar": user.Avatar})
					}
					sess, err := s.CreateToken(ctx, models.TokenTypeLoginSession, userInfo)
					if err != nil {
						return "", err
					}

					jwtSecret := config.Get().Global.GetJwtSecret()
					token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
						Id:        sess.Id,
						ExpiresAt: sess.Expiry.Unix(),
						IssuedAt:  time.Now().Unix(),
						NotBefore: time.Now().Unix(),
					}).SignedString([]byte(jwtSecret))
					if err != nil {
						return "", errors.NewServerError(500, err.Error())
					}

					cookie := http.Cookie{
						Name:  global.LoginSession,
						Value: token,
						Path:  "/",
					}
					if httpExternalURL, ok := ctx.Value(global.HTTPExternalURLKey).(string); ok && len(httpExternalURL) > 0 {
						if extURL, err := url.Parse(httpExternalURL); err == nil {
							cookie.Path = extURL.Path
						}
					}

					if oAuthOptions.AutoLogin {
						cookie.Expires = sess.Expiry
						http.SetCookie(stdResp, &http.Cookie{
							Name:    global.CookieAutoLogin,
							Value:   "true",
							Path:    cookie.Path,
							Expires: sess.Expiry,
						})
					} else {
						http.SetCookie(stdResp, &http.Cookie{
							Name:  global.CookieAutoLogin,
							Value: "false",
							Path:  cookie.Path,
						})
					}
					http.SetCookie(stdResp, &cookie)
					if oriRedirectURI := stdReq.QueryParameter("redirect_uri"); len(oriRedirectURI) > 0 {
						http.Redirect(stdResp.ResponseWriter, stdReq.Request, oriRedirectURI, http.StatusFound)
					} else {
						http.Redirect(stdResp.ResponseWriter, stdReq.Request, w.M(common.GetWebURL(ctx)), http.StatusFound)
					}
				} else {
					http.Redirect(stdResp.ResponseWriter, stdReq.Request, w.M(common.GetWebURL(ctx, common.WithSubPages("warning"), common.WithQuery(url.Values{"message": []string{"The session has expired."}}))), http.StatusFound)
				}
			}
		}
		return nil, nil
	}
}

func MakeUserLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		begin := time.Now()
		req := request.(Requester).GetRequestData().(*UserLoginRequest)
		resp := SimpleResponseWrapper[*UserLoginResponseData]{}
		logger := logs.GetContextLogger(ctx)
		var user *models.User

		defer func() {
			var userId, username string
			if user != nil {
				userId = user.Id
				username = user.Username
			}
			eventId, message, status, took := GetEventMeta(ctx, "UserLogin", begin, err, response)
			if e := s.PostEventLog(ctx, eventId, userId, username, "", "UserLogin", message, status, took); e != nil {
				level.Error(logs.GetContextLogger(ctx)).Log("failed to post event log", "err", e)
			}
		}()

		if config.Get().GetGlobal().DisableLoginForm {
			return nil, errors.ParameterError("unsupported login type")
		}
		switch req.Type {
		case LoginType_mfa_sms, LoginType_mfa_email, LoginType_normal, LoginType_mfa_totp:
			if len(req.Username) == 0 || len(req.Password) == 0 {
				resp.Error = errors.ParameterError("username or password")
				return resp, err
			}
			user, resp.Error = s.VerifyPassword(ctx, req.Username, string(req.Password), false)
			if user == nil || resp.Error != nil {
				return resp, nil
			}

			forceMFA := user.IsForceMfa()
			if forceMFA {
				method := getMFAMethod(user)
				switch req.Type {
				case LoginType_mfa_totp:
					if !method.Has(req.Type) {
						resp.Error = errors.NewServerError(500, "The authentication method is not supported.")
						return resp, nil
					}
					secret, err := user.ExtendedData.GetSecret()
					if err != nil {
						resp.Error = errors.NewServerError(500, "failed to get totp settings")
						return resp, nil
					} else if len(secret) == 0 {
						resp.Error = errors.NewServerError(500, "can't get totp settings")
						return resp, nil
					}
					nowTime := time.Now()
					ts := nowTime.Add(time.Second * time.Duration(-(nowTime.Second() % 30))).Unix()

					totp := gotp.NewDefaultTOTP(secret)
					if !totp.Verify(req.Code, ts) {
						if !totp.Verify(req.Code, ts-30) {
							resp.Error = errors.NewServerError(http.StatusBadRequest, "The verification code is invalid or expired")
							return resp, nil
						}
					}
				case LoginType_mfa_sms, LoginType_mfa_email:
					if !method.Has(req.Type) {
						resp.Error = errors.NewServerError(500, "The authentication method is not supported.")
						return resp, nil
					}
					if len(req.Token) == 0 {
						resp.Error = errors.ParameterError("token")
						return resp, nil
					}
					var code LoginCode
					if !s.VerifyToken(ctx, req.Token, models.TokenTypeLoginCode, &code) {
						resp.Error = errors.ParameterError("token")
						return resp, nil
					}
					if code.Code != req.Code || user.Id != code.UserId || req.Type != code.Type {
						resp.Error = errors.ParameterError("code")
						return resp, nil
					}
				default:
					if method.Len() != 0 {
						resp.Data = &UserLoginResponseData{NextMethod: method.List()}
					} else {
						token, err := s.CreateToken(ctx, models.TokenTypeEnableMFA, user)
						if err != nil {
							return nil, errors.NewServerError(http.StatusInternalServerError, "Failed to create token")
						}
						resp.Data = &UserLoginResponseData{
							NextMethod:  []LoginType{LoginType_enable_mfa_totp},
							Token:       token.Id,
							PhoneNumber: user.PhoneNumber,
							Email:       user.Email,
						}
						if len(user.Email) > 0 {
							resp.Data.NextMethod = append(resp.Data.NextMethod, LoginType_enable_mfa_email)
						}
						if len(user.PhoneNumber) > 0 {
							resp.Data.NextMethod = append(resp.Data.NextMethod, LoginType_enable_mfa_sms)
						}
					}
					resp.Success = false
					return resp, nil
				}
			}
		case LoginType_sms, LoginType_email:
			var code LoginCode
			if !s.VerifyToken(ctx, req.Token, models.TokenTypeLoginCode, &code) {
				resp.Error = errors.ParameterError("token")
				return resp, nil
			}

			if code.Code != req.Code || req.Type != code.Type {
				resp.Error = errors.ParameterError("code")
				return resp, nil
			}

			switch req.Type {
			case LoginType_email:
				user, err = s.GetUser(ctx, opts.WithEmail(req.Email), opts.WithUserExt)
				if err != nil {
					level.Warn(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to get user", "email", req.Email)
					return resp, nil
				}
				if user.ExtendedData != nil && user.ExtendedData.ForceMFA {
					return nil, errors.NewServerError(500, "The user has opened MFA and does not support logging in using this method.")
				}
				if user.Id != code.UserId {
					resp.Error = errors.ParameterError("token or email")
					return resp, nil
				}
			case LoginType_sms:
				user, err = s.GetUser(ctx, opts.WithEmail(req.Phone), opts.WithUserExt)
				if err != nil {
					level.Warn(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to get user", "username", req.Username)
					return resp, nil
				}
				if user.ExtendedData != nil && user.ExtendedData.ForceMFA {
					return nil, errors.NewServerError(500, "The user has opened MFA and does not support logging in using this method.")
				}
				if user.Id != code.UserId {
					resp.Error = errors.ParameterError("token or phone")
					return resp, nil
				}
			default:
				return nil, errors.ParameterError("type")
			}
		case LoginType_enable_mfa_email:
			user = new(models.User)
			if !s.VerifyToken(ctx, req.Token, models.TokenTypeEnableMFA, user) {
				resp.Error = errors.ParameterError("token")
				return resp, nil
			}
			if resp.Error = s.PatchUserExtData(ctx, user.Id, map[string]interface{}{
				"email_as_mfa": true,
			}); resp.Error != nil {
				return resp, nil
			}
		case LoginType_enable_mfa_totp:
			var secret TOTPSecret
			if !s.VerifyToken(ctx, req.Token, models.TokenTypeTotpSecret, &secret) {
				resp.Error = errors.ParameterError("token")
				return resp, nil
			}

			sec, err := secret.GetSecret()
			if err != nil {
				resp.Error = errors.WithServerError(http.StatusInternalServerError, err, "Failed to general secret")
				return resp, nil
			}
			nowTime := time.Now()
			ts := nowTime.Add(time.Second * time.Duration(-(nowTime.Second() % 30))).Unix()
			totp := gotp.NewDefaultTOTP(sec)
			if !totp.Verify(req.FirstCode, ts-30) {
				resp.Error = errors.NewServerError(http.StatusBadRequest, "The first code is invalid or expired")
			} else if !totp.Verify(req.SecondCode, ts) {
				resp.Error = errors.NewServerError(http.StatusBadRequest, "The second code is invalid or expired")
			} else {
				resp.Error = s.CreateTOTP(ctx, secret.User.Id, sec)
				if resp.Error == nil {
					_ = s.DeleteToken(ctx, models.TokenTypeTotpSecret, req.Token)
				}
			}

		default:
			resp.Error = errors.ParameterError("type")
			return resp, nil
		}
		if err = s.VerifyUserStatus(ctx, user, false); err != nil {
			return nil, err
		}
		app, err := s.GetAppInfo(ctx, opts.WithBasic, opts.WithUsers(user.Id), opts.WithAppName(config.Get().GetGlobal().GetAppName()))
		if err != nil && !errors.IsNotFount(err) {
			level.Error(logger).Log("msg", "failed to get app info", "err", err)
		} else if app != nil {
			role, err := s.GetAppRoleByUserId(ctx, app.Id, user.Id)
			if err == nil {
				user.RoleId = role.Id
				user.Role = role.Name
			} else if !errors.IsNotFount(err) {
				level.Error(logger).Log("msg", "failed to get app role", "err", err)
			}
		}
		if user.ExtendedData == nil {
			user.ExtendedData = &models.UserExt{}
		}
		user.ExtendedData.LoginTime = time.Now()

		if err = s.PatchUserExtData(ctx, user.Id, map[string]interface{}{"login_time": time.Now()}); err != nil {
			return nil, errors.NewServerError(http.StatusInternalServerError, "failed to active user.")
		}
		stdResp := request.(RestfulRequester).GetRestfulResponse().ResponseWriter
		sess, err := s.CreateToken(ctx, models.TokenTypeLoginSession, user)
		if err != nil {
			return "", err
		}

		jwtSecret := config.Get().Global.GetJwtSecret()
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Id:        sess.Id,
			ExpiresAt: sess.Expiry.Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		}).SignedString([]byte(jwtSecret))
		if err != nil {
			return "", errors.NewServerError(500, err.Error())
		}

		cookie := http.Cookie{
			Name:  global.LoginSession,
			Value: token,
			Path:  "/",
		}
		if httpExternalURL, ok := ctx.Value(global.HTTPExternalURLKey).(string); ok && len(httpExternalURL) > 0 {
			if extURL, err := url.Parse(httpExternalURL); err == nil {
				cookie.Path = extURL.Path
			}
		}

		if req.AutoLogin {
			cookie.Expires = sess.Expiry
			http.SetCookie(stdResp, &http.Cookie{
				Name:    global.CookieAutoLogin,
				Value:   "true",
				Path:    cookie.Path,
				Expires: sess.Expiry,
			})
		} else {
			http.SetCookie(stdResp, &http.Cookie{
				Name:  global.CookieAutoLogin,
				Value: "false",
				Path:  cookie.Path,
			})
		}
		http.SetCookie(stdResp, &cookie)

		return &resp, nil
	}
}

func MakeUserLogoutEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := SimpleResponseWrapper[interface{}]{}
		loginCookie, err := request.(RestfulRequester).GetRestfulRequest().Request.Cookie(global.LoginSession)
		if err != nil {
			resp.Error = errors.BadRequestError()
		} else if len(loginCookie.Value) > 0 {
			var claims jwt.StandardClaims
			token, err := jwt.ParseWithClaims(loginCookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected login method %v", token.Header["alg"])
				}
				return []byte(config.Get().Global.GetJwtSecret()), nil
			})
			if err == nil && token.Valid {
				if err = s.DeleteLoginSession(ctx, claims.Id); err != nil {
					resp.Error = errors.InternalServerError()
					return resp, nil
				}
			}

			if httpExternalURL, ok := ctx.Value(global.HTTPExternalURLKey).(string); ok && len(httpExternalURL) > 0 {
				if extURL, err := url.Parse(httpExternalURL); err == nil {
					loginCookie.Path = extURL.Path
				}
			}
			respWriter := request.(RestfulRequester).GetRestfulResponse().ResponseWriter
			http.SetCookie(respWriter, &http.Cookie{
				Name:    global.LoginSession,
				Value:   loginCookie.Value,
				Path:    loginCookie.Path,
				Expires: time.Now(),
			})
			if loginCookie.Path != "/" {
				http.SetCookie(respWriter, &http.Cookie{
					Name:    global.LoginSession,
					Value:   loginCookie.Value,
					Path:    "/",
					Expires: time.Now(),
				})
			}
		} else {
			resp.Error = errors.NewServerError(http.StatusUnauthorized, "Invalid identity information")
		}
		return &resp, nil
	}
}

func MakeAuthenticationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*AuthenticationRequest)
		return s.Authentication(ctx, req.AuthMethod, req.AuthAlgorithm, req.AuthKey, req.AuthSecret, req.Payload, req.AuthSign)
	}
}

type GetSessionParams struct {
	Token     string
	TokenType models.TokenType
}

func UpdateSessionExpires(ctx context.Context, s service.Service, id string, stdResp http.ResponseWriter, loginTime time.Time) {
	expiry := models.TokenTypeLoginSession.GetExpiry()
	maxExpire := loginTime.Add(time.Hour * time.Duration(config.GetRuntimeConfig().GetLoginSessionMaxTime()))
	if expiry.After(maxExpire) {
		expiry = maxExpire
	}
	if err := s.UpdateTokenExpires(ctx, id, expiry); err != nil {
		logger := logs.WithPrint(fmt.Sprintf("%+v", err))(logs.GetContextLogger(ctx))
		level.Error(logger).Log("msg", "failed to update session expires")
	} else {
		jwtSecret := config.Get().Global.GetJwtSecret()
		newToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Id:        id,
			ExpiresAt: expiry.Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		}).SignedString([]byte(jwtSecret))
		if err != nil {
			logger := logs.WithPrint(fmt.Sprintf("%+v", err))(logs.GetContextLogger(ctx))
			level.Error(logger).Log("msg", "failed to create jwt token.")
		}

		cookie := http.Cookie{
			Name:  global.LoginSession,
			Value: newToken,
			Path:  "/",
		}
		if httpExternalURL, ok := ctx.Value(global.HTTPExternalURLKey).(string); ok && len(httpExternalURL) > 0 {
			if extURL, err := url.Parse(httpExternalURL); err == nil {
				cookie.Path = extURL.Path
			}
		}
		cookie.Expires = expiry

		http.SetCookie(stdResp, &http.Cookie{
			Name:    global.CookieAutoLogin,
			Value:   "true",
			Path:    cookie.Path,
			Expires: expiry,
		})

		http.SetCookie(stdResp, &cookie)
	}
}

func MakeGetSessionByTokenEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		stdReq := request.(RestfulRequester).GetRestfulRequest().Request
		params := request.(Requester).GetRequestData().(*GetSessionParams)
		var resp *models.User
		if len(params.Token) > 0 {
			var claims jwt.StandardClaims
			token, err := jwt.ParseWithClaims(params.Token, &claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected login method %v", token.Header["alg"])
				}
				return []byte(config.Get().Global.GetJwtSecret()), nil
			})
			if err != nil {
				logger := logs.WithPrint(fmt.Sprintf("%+v", err))(logs.GetContextLogger(ctx))
				level.Error(logger).Log("err", err, "msg", "Invalid token")
				err = errors.NotLoginError()
			} else if !token.Valid {
				logger := logs.WithPrint(fmt.Sprintf("%+v", err))(logs.GetContextLogger(ctx))
				level.Error(logger).Log("msg", "Invalid token")
				err = errors.NotLoginError()
			} else if err = s.GetSessionByToken(ctx, claims.Id, params.TokenType, &resp); err != nil {
				if err != errors.NotLoginError() {
					logger := logs.WithPrint(fmt.Sprintf("%+v", err))(logs.GetContextLogger(ctx))
					level.Error(logger).Log("err", err, "msg", "failed to get session")
					err = errors.NotLoginError()
				}
			}
			if err == nil && params.TokenType == models.TokenTypeLoginSession {
				if resp.ExtendedData == nil {
					resp.ExtendedData = new(models.UserExt)
				}
				if !resp.ExtendedData.LoginTime.IsZero() && time.Since(resp.ExtendedData.LoginTime) > time.Hour {
					maxTime := config.GetRuntimeConfig().GetLoginSessionMaxTime()
					if time.Since(resp.ExtendedData.LoginTime) < (time.Hour * time.Duration(maxTime)) {
						if autoLoginCookie, _ := stdReq.Cookie(global.CookieAutoLogin); autoLoginCookie != nil && autoLoginCookie.Value == "true" {
							stdResp := request.(RestfulRequester).GetRestfulResponse().ResponseWriter
							UpdateSessionExpires(ctx, s, claims.Id, stdResp, resp.ExtendedData.LoginTime)
						}
					}
				}
			}
		} else {
			err = errors.NotLoginError()
		}

		return resp, err
	}
}

func MakeGetProxySessionByTokenEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		params := request.(*GetSessionParams)
		var session []*models.ProxySession
		if len(params.Token) > 0 {
			if err = s.GetSessionByToken(ctx, params.Token, params.TokenType, &session); err != nil {
				if err != errors.NotLoginError() {
					level.Error(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to get session")
				}
			} else if len(session) == 1 {
				return session[0], nil
			} else if len(session) > 1 {
				return nil, errors.NewServerError(500, "proxy configuration exception")
			}
		} else {
			err = errors.NotLoginError()
		}
		return nil, err
	}
}

func MakeGetSessionsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetSessionsRequest)
		resp := NewBaseListResponse[[]*models.Token](&req.BaseListRequest)
		resp.Total, resp.Data, resp.BaseResponse.Error = s.GetSessions(ctx, req.UserId, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeGetCurrentUserSessionsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetSessionsRequest)
		resp := NewBaseListResponse[[]*models.Token](&req.BaseListRequest)
		user, ok := ctx.Value(global.MetaUser).(*models.User)
		if !ok || user == nil {
			return nil, errors.NotLoginError()
		}

		resp.Total, resp.Data, resp.BaseResponse.Error = s.GetSessions(ctx, user.Id, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeDeleteSessionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteSessionRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		resp.Error = s.DeleteToken(ctx, models.TokenTypeLoginSession, req.Id)
		return &resp, nil
	}
}

func MakeDeleteCurrentUserSessionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteSessionRequest)
		user := ctx.Value(global.MetaUser).(*models.User)
		if user == nil {
			return nil, errors.NotLoginError()
		}
		resp := SimpleResponseWrapper[interface{}]{}
		if s.VerifyToken(ctx, req.Id, models.TokenTypeLoginSession, nil, user.Id) {
			resp.Error = s.DeleteToken(ctx, models.TokenTypeLoginSession, req.Id)
		}
		return &resp, nil
	}
}
