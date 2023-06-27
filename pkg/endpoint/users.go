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
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MicroOps-cn/fuck/crypto"
	logs "github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"
	"github.com/gogo/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/xlzd/gotp"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
)

func MakeCurrentUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := SimpleResponseWrapper[interface{}]{}
		if user, ok := ctx.Value(global.MetaUser).(*models.User); ok && user != nil {
			return user, nil
		}
		if restfulRequester, ok := request.(RestfulRequester); ok {
			restfulRequest := restfulRequester.GetRestfulRequest()
			if auth := restfulRequest.Request.Header.Get("Authorization"); len(auth) > 0 {
				if strings.HasPrefix(auth, "Bearer ") {
					var users models.Users
					err = s.GetSessionByToken(ctx, strings.TrimPrefix(auth, "Bearer "), models.TokenTypeToken, &users)
					if err != nil {
						resp.Error = errors.NotLoginError()
						return resp, nil
					} else if len(users) > 0 {
						return users[0], nil
					}
				}
			}
		}
		return resp, nil
	}
}

func MakeUpdateCurrentUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdateUserRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		user, ok := ctx.Value(global.MetaUser).(*models.User)
		if !ok || user == nil {
			resp.Error = errors.NotLoginError()
			return resp, nil
		}
		if resp.Error = s.UpdateUser(ctx, &models.User{
			Model: models.Model{
				Id:       req.Id,
				IsDelete: req.IsDelete,
			},
			Username:    req.Username,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			FullName:    req.FullName,
			Avatar:      req.Avatar,
			Status:      req.Status,
		}, "email", "phone_number", "full_name", "avatar"); resp.Error != nil {
			return resp, nil
		}
		resp.Error = s.UpdateUserSession(ctx, user.Id)
		return resp, nil
	}
}

func (r PatchCurrentUserRequest) Map() map[string]interface{} {
	m := map[string]interface{}{}
	if r.EmailAsMfa != nil {
		m["email_as_mfa"] = r.EmailAsMfa
	}
	if r.SmsAsMfa != nil {
		m["sms_as_mfa"] = r.SmsAsMfa
	}
	if r.TotpAsMfa != nil {
		m["totp_as_mfa"] = r.TotpAsMfa
	}
	return m
}

func MakePatchCurrentUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchCurrentUserRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		user, ok := ctx.Value(global.MetaUser).(*models.User)
		if !ok || user == nil {
			resp.Error = errors.NotLoginError()
			return resp, nil
		}
		if patch := req.Map(); len(patch) > 0 {
			resp.Error = s.PatchUserExtData(ctx, user.Id, patch)
			if resp.Error == nil {
				resp.Error = s.UpdateUserSession(ctx, user.Id)
			}
		}
		return resp, nil
	}
}

func PasswordComplexityVerification(username, password string) (bool, error) {
	var (
		uppercase int
		lowercase int
		number    int
		special   int
	)
	allow := config.GetRuntimeConfig().Security.PasswordComplexity
	if minLength := config.GetRuntimeConfig().Security.PasswordMinLength; minLength > 0 {
		if len(password) < int(minLength) {
			return false, errors.NewServerError(400, "passwords too short. ", errors.CodePasswordTooShort)
		}
	}
	if allow == config.PasswordComplexity_unsafe {
		return true, nil
	}
	if len(username) != 0 && strings.Contains(password, username) {
		return false, errors.NewServerError(400, "passwords should not include username. ", errors.CodePasswordCannotContainUsername)
	}
	for _, chr := range password {
		if chr >= 'A' && chr <= 'Z' {
			uppercase++
		} else if chr >= 'a' && chr <= 'z' {
			lowercase++
		} else if chr >= '0' && chr <= '9' {
			number++
		} else {
			special++
		}
	}
	typeNum := 0
	if uppercase > 0 {
		typeNum++
	}
	if number > 0 {
		typeNum++
	}
	if lowercase > 0 {
		typeNum++
	}
	if special > 0 {
		typeNum++
	}
	switch allow {
	case config.PasswordComplexity_general:
		if typeNum < 2 {
			return false, errors.NewServerError(400, "it must be composed of at least two combinations of uppercase letters, lowercase letters, numbers, and special characters. ", errors.CodePasswordBaseGeneralTooSimple)
		}
	case config.PasswordComplexity_safe:
		if typeNum < 3 {
			return false, errors.NewServerError(400, "it must be composed of at least any three combinations of uppercase letters, lowercase letters, numbers, and special characters. ", errors.CodePasswordBaseSafeTooSimple)
		}
	case config.PasswordComplexity_very_safe:
		if typeNum < 4 {
			return false, errors.NewServerError(400, "must contain uppercase and lowercase letters, numbers, and special characters. ", errors.CodePasswordBaseVerySafeTooSimple)
		}
	}
	return true, nil
}

func (r *ResetUserPasswordRequest) GetToken() string {
	if r != nil && r.Token != nil {
		return string(*r.Token)
	}
	return ""
}

func MakeResetUserPasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		begin := time.Now()
		req := request.(Requester).GetRequestData().(*ResetUserPasswordRequest)
		resp := TotalResponseWrapper[interface{}]{}
		var user *models.User
		defer func() {
			var userId, username string
			if user != nil {
				userId = user.Id
				username = user.Username
			}
			eventId, message, status, took := GetEventMeta(ctx, "ResetUserPassword", begin, err, resp)
			if e := s.PostEventLog(ctx, eventId, userId, username, "", "ResetUserPassword", message, status, took); e != nil {
				level.Error(logs.GetContextLogger(ctx)).Log("failed to post event log", "err", e)
			}
		}()
		if _, err = PasswordComplexityVerification(req.Username, string(*req.NewPassword)); err != nil {
			resp.Error = err
			return resp, nil
		}
		if len(req.GetToken()) > 0 && len(req.UserId) > 0 {
			user = new(models.User)
			if s.VerifyToken(ctx, req.GetToken(), models.TokenTypeResetPassword, user, req.UserId) {
				if err = s.ResetPassword(ctx, user.Id, string(*req.NewPassword)); err != nil {
					return nil, err
				}
				if err = s.DeleteToken(ctx, models.TokenTypeResetPassword, req.GetToken()); err != nil {
					level.Error(logs.GetContextLogger(ctx)).Log("msg", "failed to delete token", "err", err)
				}
			} else {
				return nil, errors.ParameterError("invalid token")
			}
		} else if req.OldPassword != nil && len(*req.OldPassword) != 0 {
			if len(req.Username) > 0 {
				user, err = s.VerifyPassword(ctx, req.Username, string(*req.OldPassword), true)
				if err != nil {
					return nil, err
				}
				if err = s.VerifyUserStatus(ctx, user, true); err != nil {
					return nil, err
				}
				if user == nil {
					return nil, errors.NewServerError(http.StatusBadRequest, "Invalid old password")
				}
				if resp.Error = s.ResetPassword(ctx, user.Id, string(*req.NewPassword)); resp.Error != nil {
					return nil, err
				}
			} else if len(req.UserId) > 0 {
				if user = s.VerifyPasswordById(ctx, req.UserId, string(*req.OldPassword), true); user != nil {
					if resp.Error = s.VerifyUserStatus(ctx, user, true); resp.Error != nil {
						return resp, err
					}
					resp.Error = s.ResetPassword(ctx, req.UserId, string(*req.NewPassword))
				} else {
					return nil, errors.UnauthorizedError()
				}
			}
		}
		return resp, nil
	}
}

func MakeForgotPasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		begin := time.Now()
		req := request.(Requester).GetRequestData().(*ForgotUserPasswordRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		var user *models.User
		defer func() {
			var userId, username string
			if user != nil {
				userId = user.Id
				username = user.Username
			}
			eventId, message, status, took := GetEventMeta(ctx, "ForgotPassword", begin, err, resp)
			if e := s.PostEventLog(ctx, eventId, userId, username, "", "ForgotPassword", message, status, took); e != nil {
				level.Error(logs.GetContextLogger(ctx)).Log("failed to post event log", "err", e)
			}
		}()

		user, err = s.GetUserInfoByUsernameAndEmail(ctx, req.Username, req.Email)
		if err != nil {
			level.Warn(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to reset password", "username", req.Username, "email", req.Email)
			return resp, nil
		}

		token, err := s.CreateToken(ctx, models.TokenTypeResetPassword, user)
		if err != nil {
			return resp, err
		}

		to := fmt.Sprintf("%s<%s>", user.FullName, user.Email)

		httpExternalURL, _ := ctx.Value(global.HTTPExternalURLKey).(string)
		err = s.SendEmail(ctx, map[string]interface{}{
			"user":            user,
			"userId":          token.GetRelationId(),
			"token":           token,
			"httpExternalURL": httpExternalURL,
			"siteTitle":       config.Get().GetGlobal().GetTitle(),
			"adminEmail":      config.Get().GetGlobal().GetAdminEmail(),
		}, "User:ResetPassword", to)
		if err != nil {
			level.Error(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to send email")
			return nil, errors.NewServerError(500, "failed to send email")
		}
		return resp, nil
	}
}

func MakeGetUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetUsersRequest)
		resp := NewBaseListResponse[[]*models.User](&req.BaseListRequest)
		if req.Status == nil {
			req.Status = w.P[models.UserMeta_UserStatus](models.UserMetaStatusAll)
		}
		resp.Total, resp.Data, resp.Error = s.GetUsers(
			ctx, req.Keywords,
			*req.Status,
			req.App, req.Current, req.PageSize,
		)
		return &resp, nil
	}
}

type PatchUsersRequest []PatchUserRequest

func (m *PatchUsersRequest) Reset()         { *m = PatchUsersRequest{} }
func (m *PatchUsersRequest) String() string { return proto.CompactTextString(m) }

func (m PatchUsersRequest) ProtoMessage() {}

func MakePatchUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchUsersRequest)
		resp := TotalResponseWrapper[interface{}]{}

		var patchUsers []map[string]interface{}
		for _, u := range *req {
			if len(u.Id) == 0 {
				return nil, errors.ParameterError("There is an empty id in the patch.")
			}
			patch := map[string]interface{}{"id": u.Id}
			if u.Status != nil {
				patch["status"] = int32(*u.Status)
			}
			if u.IsDelete != nil {
				patch["isDelete"] = u.IsDelete
			}
			patchUsers = append(patchUsers, patch)
		}
		resp.Total, resp.Error = s.PatchUsers(ctx, patchUsers)
		return &resp, nil
	}
}

type DeleteUsersRequest []DeleteUserRequest

func MakeDeleteUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteUsersRequest)
		resp := TotalResponseWrapper[interface{}]{}
		var delUsers []string
		for _, u := range *req {
			if len(u.Id) == 0 {
				return nil, errors.ParameterError("There is an empty id in the request.")
			}
			delUsers = append(delUsers, u.Id)
		}
		resp.Total, resp.Error = s.DeleteUsers(ctx, delUsers)
		return &resp, nil
	}
}

func MakeUpdateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdateUserRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		if resp.Error = s.UpdateUser(ctx, &models.User{
			Model: models.Model{
				Id:       req.Id,
				IsDelete: req.IsDelete,
			},
			Username:    req.Username,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			FullName:    req.FullName,
			Avatar:      req.Avatar,
			Status:      req.Status,
			Apps: w.Map(req.Apps, func(app *UserApp) *models.App {
				return &models.App{
					Model:  models.Model{Id: app.Id},
					RoleId: app.RoleId,
				}
			}),
		}); resp.Error != nil {
			resp.Error = errors.WithServerError(500, resp.Error, "failed to update user")
		}
		return &resp, nil
	}
}

func MakeGetUserInfoEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetUserRequest)
		resp := SimpleResponseWrapper[*models.User]{}
		resp.Data, resp.Error = s.GetUser(ctx, opts.WithUserId(req.Id), opts.WithApps)
		return &resp, nil
	}
}

func MakeCreateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateUserRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		resp.Error = s.CreateUser(ctx, &models.User{
			Username:    req.Username,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			FullName:    req.FullName,
			Avatar:      req.Avatar,
			Status:      models.UserMeta_user_inactive,
			Apps: w.Map(req.Apps, func(app *UserApp) *models.App {
				return &models.App{
					Model:  models.Model{Id: app.Id},
					RoleId: app.RoleId,
				}
			}),
		})
		return &resp, nil
	}
}

type PatchUserResponse struct {
	User *models.User `json:",inline"`
}

func MakePatchUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchUserRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		if len(req.Id) == 0 {
			return nil, errors.ParameterError("There is an empty id in the patch.")
		}
		patch := map[string]interface{}{"id": req.Id}
		if req.Status != nil {
			patch["status"] = req.Status
		}
		if req.IsDelete != nil {
			patch["isDelete"] = req.IsDelete
		}
		resp.Error = s.PatchUser(ctx, patch)
		return &resp, nil
	}
}

func MakeDeleteUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteUserRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		resp.Error = s.DeleteUser(ctx, req.Id)
		return &resp, nil
	}
}

func MakeCreateUserKeyEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateUserKeyRequest)
		resp := SimpleResponseWrapper[*models.UserKey]{}
		resp.Data, resp.Error = s.CreateUserKey(ctx, req.UserId, req.Name)
		return &resp, nil
	}
}

func MakeDeleteUserKeyEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteUserKeyRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		resp.Error = s.DeleteUserKey(ctx, req.UserId, req.Id)
		return &resp, nil
	}
}

func MakeGetUserKeysEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetUserKeysRequest)
		resp := NewBaseListResponse[[]*models.UserKey](&req.BaseListRequest)
		resp.Total, resp.Data, resp.Error = s.GetUserKeys(ctx, req.UserId, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeCreateKeyEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateKeyRequest)
		resp := SimpleResponseWrapper[*models.UserKey]{}
		user, ok := ctx.Value(global.MetaUser).(*models.User)
		if !ok || user == nil {
			return nil, errors.NotLoginError()
		}
		if user.Id == req.UserId && user.Status == models.UserMeta_normal {
			resp.Data, resp.Error = s.CreateUserKey(ctx, req.UserId, req.Name)
			return &resp, nil
		}

		return nil, errors.StatusNotFound("user")
	}
}

func MakeSendActivationMailEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*SendActivationMailRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		user, err := s.GetUserInfo(ctx, req.UserId, "")
		if err != nil {
			return nil, err
		}

		if !user.Status.IsAnyOne(models.UserMeta_user_inactive, models.UserMeta_password_expired) {
			return nil, errors.NewServerError(http.StatusInternalServerError, "Unknown user's status")
		}
		if err = s.PatchUserExtData(ctx, req.UserId, map[string]interface{}{"activation_time": time.Now().UTC()}); err != nil {
			return nil, errors.NewServerError(http.StatusInternalServerError, "failed to active user.")
		}
		token, err := s.CreateToken(ctx, models.TokenTypeActive, user)
		if err != nil {
			return nil, errors.NewServerError(http.StatusInternalServerError, "Failed to create token")
		}
		to := fmt.Sprintf("%s<%s>", user.FullName, user.Email)
		httpExternalURL, _ := ctx.Value(global.HTTPExternalURLKey).(string)
		err = s.SendEmail(ctx, map[string]interface{}{
			"user":            user,
			"token":           token,
			"userId":          token.ParentId,
			"httpExternalURL": httpExternalURL,
			"siteTitle":       config.Get().GetGlobal().GetTitle(),
			"adminEmail":      config.Get().GetGlobal().GetAdminEmail(),
		}, "User:ActivateAccount", to)
		if err != nil {
			level.Error(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to send email")
			return nil, errors.NewServerError(500, "failed to send email")
		}
		return &resp, nil
	}
}

func (r *ActivateAccountRequest) GetToken() string {
	if r != nil && r.Token != nil {
		return string(*r.Token)
	}
	return ""
}

func MakeActivateAccountEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		begin := time.Now()
		req := request.(Requester).GetRequestData().(*ActivateAccountRequest)
		resp := TotalResponseWrapper[interface{}]{}
		defer func() {
			eventId, message, status, took := GetEventMeta(ctx, "ActivateAccount", begin, err, response)
			if e := s.PostEventLog(ctx, eventId, req.UserId, "", "", "ActivateAccount", message, status, took); e != nil {
				level.Error(logs.GetContextLogger(ctx)).Log("failed to post event log", "err", e)
			}
		}()
		if _, err = PasswordComplexityVerification("", string(*req.NewPassword)); err != nil {
			resp.Error = err
			return resp, nil
		}
		if len(req.GetToken()) > 0 {
			if s.VerifyToken(ctx, req.GetToken(), models.TokenTypeActive, nil, req.UserId) {
				resp.Error = s.ResetPassword(ctx, req.UserId, string(*req.NewPassword))
				if resp.Error == nil {
					if err = s.DeleteToken(ctx, models.TokenTypeActive, req.GetToken()); err != nil {
						level.Error(logs.GetContextLogger(ctx)).Log("msg", "failed to delete token", "err", err)
					}
				}
			} else {
				return nil, errors.ParameterError("invalid token")
			}
		} else {
			return nil, errors.ParameterError("invalid token")
		}

		return resp, nil
	}
}

type TOTPSecret struct {
	User   *models.User `json:"user"`
	Secret string       `json:"secret"`
	Salt   string       `json:"salt"`
}

func (s *TOTPSecret) SetSecret(secret string) (err error) {
	globalSecret := config.Get().GetGlobal().GetSecret()
	if globalSecret == "" {
		return errors.NewServerError(500, "global secret is not set")
	}
	s.Salt = uuid.NewV4().String()
	key := sha256.Sum256([]byte(s.Salt + (globalSecret)))
	sec, err := crypto.NewAESCipher(key[:]).CBCEncrypt([]byte(secret))
	s.Secret = base64.StdEncoding.EncodeToString(sec)
	return err
}

func (s TOTPSecret) GetSecret() (secret string, err error) {
	if len(s.Secret) == 0 || len(s.Salt) == 0 {
		return "", nil
	}
	globalSecret := config.Get().GetGlobal().GetSecret()
	if globalSecret == "" {
		return "", errors.NewServerError(500, "global secret is not set")
	}
	key := sha256.Sum256([]byte(s.Salt + (globalSecret)))
	decoded, err := base64.StdEncoding.DecodeString(s.Secret)
	if err != nil {
		return "", err
	}
	sec, err := crypto.NewAESCipher(key[:]).CBCDecrypt(decoded)
	return string(sec), err
}

func (r *CreateTOTPSecretRequest) GetToken() string {
	if r != nil && r.Token != nil {
		return string(*r.Token)
	}
	return ""
}

func MakeCreateTOTPSecretEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateTOTPSecretRequest)
		resp := &SimpleResponseWrapper[CreateTOTPSecretResponseData]{}
		var user *models.User
		if len(req.GetToken()) != 0 {
			user = new(models.User)
			if !s.VerifyToken(ctx, req.GetToken(), models.TokenTypeEnableMFA, user) {
				return nil, errors.NewServerError(400, "Invalid token.")
			}
		} else if user, _ = ctx.Value(global.MetaUser).(*models.User); user == nil {
			return nil, errors.NotLoginError()
		}
		randomSecret := gotp.RandomSecret(128)
		secret := TOTPSecret{
			User: user,
		}

		if err = secret.SetSecret(randomSecret); err != nil {
			resp.Error = errors.WithServerError(http.StatusInternalServerError, err, "Failed to general secret")
			return resp, nil
		}
		token, err := s.CreateToken(ctx, models.TokenTypeTotpSecret, &secret)
		if err != nil {
			resp.Error = errors.WithServerError(http.StatusInternalServerError, err, "Failed to general secret")
			return resp, nil
		}
		resp.Data.Secret = gotp.NewDefaultTOTP(randomSecret).ProvisioningUri(user.Username, config.Get().GetGlobal().GetAppName())
		resp.Data.Token = token.Id
		return resp, nil
	}
}

func (r *CreateTOTPRequest) GetToken() string {
	if r != nil && r.Token != nil {
		return string(*r.Token)
	}
	return ""
}

func MakeCreateTOTPEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateTOTPRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		user := ctx.Value(global.MetaUser).(*models.User)
		if user == nil {
			return nil, errors.NotLoginError()
		}
		var secret TOTPSecret
		if !s.VerifyToken(ctx, req.GetToken(), models.TokenTypeTotpSecret, &secret) {
			resp.Error = errors.ParameterError("token")
			return resp, nil
		}
		if user.Id != secret.User.Id {
			resp.Error = errors.ParameterError("token")
			return resp, nil
		}
		sec, err := secret.GetSecret()
		if err != nil {
			resp.Error = errors.WithServerError(http.StatusInternalServerError, err, "Failed to get secret")
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
			resp.Error = s.CreateTOTP(ctx, user.Id, sec)
			if resp.Error == nil {
				_ = s.DeleteToken(ctx, models.TokenTypeTotpSecret, req.GetToken())
			}
		}

		return resp, nil
	}
}
