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
	"github.com/MicroOps-cn/fuck/sets"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
	uuid "github.com/satori/go.uuid"
	"github.com/xlzd/gotp"
	"net/http"
	"strings"
	"time"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

type LoginCode struct {
	UserId string `json:"userId"`
	Code   string `json:"code"`
}

func MakeSendLoginCaptchaEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*SendLoginCaptchaRequest)
		resp := SimpleResponseWrapper[*SendLoginCaptchaResponseData]{}
		switch req.Type {
		case LoginType_mfa_email:
			user, err := s.GetUserInfoByUsernameAndEmail(ctx, req.Username, req.Target)
			if err != nil {
				level.Warn(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed send login captcha", "username", req.Username)
				return resp, nil
			}
			if user.Status.Is(models.UserMeta_normal) {
				loginCode := LoginCode{UserId: user.Id, Code: strings.ToUpper(uuid.NewV4().String()[:6])}
				token, err := s.CreateToken(ctx, models.TokenTypeLoginCode, &loginCode)
				if err != nil {
					return nil, errors.NewServerError(http.StatusInternalServerError, "Failed to create token")
				}
				to := fmt.Sprintf("%s<%s>", user.FullName, user.Email)
				err = s.SendEmail(ctx, map[string]interface{}{
					"user":   user,
					"token":  token,
					"code":   loginCode.Code,
					"userId": user.Id,
				}, "User:SendLoginCaptcha", to)
				if err != nil {
					level.Error(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to send email")
					return nil, errors.NewServerError(500, "failed to send email")
				}
				resp.Data = &SendLoginCaptchaResponseData{Token: token.Id}
				return &resp, nil
			}
		default:
			return nil, errors.ParameterError("type")
		}

		return nil, errors.StatusNotFound("user")
	}
}

func getMFAMethod(user *models.User) sets.Set[LoginType] {
	method := sets.New[LoginType]()
	userExt := user.ExtendedData
	if userExt.EmailAsMFA {
		method.Insert(LoginType_mfa_email)
	}
	if userExt.TOTPAsMFA {
		method.Insert(LoginType_mfa_totp)
	}
	if userExt.SmsAsMFA {
		method.Insert(LoginType_mfa_sms)
	}
	return method
}

func MakeUserLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UserLoginRequest)
		resp := SimpleResponseWrapper[*UserLoginResponseData]{}
		user, err := s.VerifyPassword(ctx, req.Username, string(req.Password))
		if user == nil || err != nil {
			resp.ErrorMessage = "Wrong user name or password. "
			if err != nil {
				resp.Error = err
			}
			return resp, nil
		}
		logger := logs.GetContextLogger(ctx)
		var forceMFA = user.ExtendedData != nil && user.ExtendedData.ForceMFA

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
			case LoginType_mfa_sms:
				resp.Error = errors.NewServerError(500, "The authentication method is not supported.")
				return resp, nil

			case LoginType_mfa_email:
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
					resp.Error = errors.ParameterError("code")
					return resp, nil
				}
				if code.Code != req.Code || user.Id != code.UserId {
					resp.Error = errors.ParameterError("code")
					return resp, nil
				}
				_ = s.DeleteToken(ctx, models.TokenTypeLoginCode, req.Token)
			default:
				resp.Data = &UserLoginResponseData{NextMethod: method.List()}
				resp.Success = false
				return resp, nil
			}
		}
		app, err := s.GetAppInfo(ctx, opts.WithBasic, opts.WithUsers(user.Id), opts.WithAppName("IDAS"))
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

		token, err := s.CreateToken(ctx, models.TokenTypeLoginSession, user)
		if err != nil {
			return "", err
		}
		cookie := http.Cookie{
			Name:  global.LoginSession,
			Value: token.Id,
			Path:  "/",
		}
		if req.AutoLogin {
			cookie.Expires = token.Expiry
		}
		request.(RestfulRequester).GetRestfulResponse().AddHeader("Set-Cookie", cookie.String())

		return &resp, nil
	}
}

func MakeUserLogoutEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := SimpleResponseWrapper[interface{}]{}
		cookie, err := request.(RestfulRequester).GetRestfulRequest().Request.Cookie(global.LoginSession)
		if err != nil {
			resp.Error = errors.BadRequestError()
		} else if len(cookie.Value) > 0 {
			for _, id := range strings.Split(cookie.Value, ",") {
				if err = s.DeleteLoginSession(ctx, id); err != nil {
					resp.Error = errors.InternalServerError()
					return resp, nil
				}
			}
			loginCookie := fmt.Sprintf("%s=%s; Path=/;Expires=%s", global.LoginSession, cookie.Value, time.Now().UTC().Format(global.LoginSessionExpiresFormat))
			request.(RestfulRequester).GetRestfulResponse().AddHeader("Set-Cookie", loginCookie)
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

func MakeGetSessionByTokenEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		params := request.(*GetSessionParams)
		var resp *models.User
		if len(params.Token) > 0 {
			if err = s.GetSessionByToken(ctx, params.Token, params.TokenType, &resp); err != nil {
				if err != errors.NotLoginError() {
					logger := logs.WithPrint(fmt.Sprintf("%+v", err))(logs.GetContextLogger(ctx))
					level.Error(logger).Log("err", err, "msg", "failed to get session")
					err = errors.NotLoginError()
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

func MakeDeleteSessionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteSessionRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		resp.Error = s.DeleteToken(ctx, models.TokenTypeLoginSession, req.Id)
		return &resp, nil
	}
}
