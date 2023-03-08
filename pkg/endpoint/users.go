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
	"net/http"
	"strings"
	"time"

	logs "github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"
	"github.com/gogo/protobuf/proto"
	"github.com/xlzd/gotp"
)

func MakeCurrentUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := SimpleResponseWrapper[interface{}]{}
		if user, ok := ctx.Value(global.MetaUser).(*models.User); ok && user != nil {
			return user, nil
		}
		if restfulRequester, ok := request.(RestfulRequester); ok {
			restfulRequester.GetRestfulRequest()
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

func MakeResetUserPasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*ResetUserPasswordRequest)
		resp := TotalResponseWrapper[interface{}]{}
		if len(req.Token) > 0 && len(req.UserId) > 0 {
			var users models.Users
			if s.VerifyToken(ctx, req.Token, models.TokenTypeResetPassword, &users, strings.Split(req.UserId, ",")...) {
				errs := errors.NewMultipleServerError(500, "")
				for _, user := range users {
					if e := s.ResetPassword(ctx, user.Id, req.NewPassword); e != nil {
						errs.Append(e)
					}
				}
				if errs.HasError() {
					return nil, errs
				}
			} else {
				return nil, errors.ParameterError("invalid token")
			}
		} else if len(req.OldPassword) != 0 {
			if len(req.Username) > 0 {
				user, err := s.VerifyPassword(ctx, req.Username, req.OldPassword)
				if err != nil {
					return nil, err
				}
				if user == nil {
					return nil, errors.NewServerError(http.StatusBadRequest, "Invalid old password")
				}
				if resp.Error = s.ResetPassword(ctx, user.Id, req.NewPassword); resp.Error != nil {
					return nil, err
				}
			} else if len(req.UserId) > 0 {
				if user := s.VerifyPasswordById(ctx, req.UserId, req.OldPassword); user != nil {
					resp.Error = s.ResetPassword(ctx, req.UserId, req.NewPassword)
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
		req := request.(Requester).GetRequestData().(*ForgotUserPasswordRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		user, err := s.GetUserInfoByUsernameAndEmail(ctx, req.Username, req.Email)
		if err != nil {
			level.Warn(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to reset password", "username", req.Username, "email", req.Email)
			return resp, nil
		}

		token, err := s.CreateToken(ctx, models.TokenTypeResetPassword, user)
		if err != nil {
			return resp, err
		}

		to := fmt.Sprintf("%s<%s>", user.FullName, user.Email)
		err = s.SendEmail(ctx, map[string]interface{}{
			"user":   user,
			"userId": token.GetRelationId(),
			"token":  token,
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
		resp := BaseResponse{}
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
		}); resp.Error != nil {
			resp.Error = errors.NewServerError(200, resp.Error.Error())
		}
		return &resp, nil
	}
}

func MakeGetUserInfoEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequest)
		resp := SimpleResponseWrapper[*models.User]{}
		resp.Data, resp.Error = s.GetUserInfo(ctx, req.Id, req.Username)
		return &resp, nil
	}
}

func MakeCreateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateUserRequest)
		resp := BaseResponse{}
		resp.Error = s.CreateUser(ctx, &models.User{
			Username:    req.Username,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			FullName:    req.FullName,
			Avatar:      req.Avatar,
			Status:      models.UserMeta_user_inactive,
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
		resp := BaseResponse{}
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
		token, err := s.CreateToken(ctx, models.TokenTypeActive, user)
		if err != nil {
			return nil, errors.NewServerError(http.StatusInternalServerError, "Failed to create token")
		}
		to := fmt.Sprintf("%s<%s>", user.FullName, user.Email)
		err = s.SendEmail(ctx, map[string]interface{}{
			"user":   user,
			"token":  token,
			"userId": token.ParentId,
		}, "User:ActivateAccount", to)
		if err != nil {
			level.Error(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to send email")
			return nil, errors.NewServerError(500, "failed to send email")
		}
		return &resp, nil
	}
}

func MakeActivateAccountEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*ActivateAccountRequest)
		resp := TotalResponseWrapper[interface{}]{}
		if len(req.Token) > 0 {
			if s.VerifyToken(ctx, req.Token, models.TokenTypeActive, nil, req.UserId) {
				resp.Error = s.ResetPassword(ctx, req.UserId, req.NewPassword)
			} else {
				return nil, errors.ParameterError("invalid token")
			}
		} else {
			return nil, errors.ParameterError("invalid token")
		}

		return resp, nil
	}
}

func MakeCreateTOTPSecretEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		secret := gotp.RandomSecret(128)
		resp := &SimpleResponseWrapper[CreateTOTPSecretResponseData]{}

		user := ctx.Value(global.MetaUser).(*models.User)
		if user == nil {
			return nil, errors.NotLoginError()
		}
		token, err := s.CreateToken(ctx, models.TokenTypeTotpSecret, secret)
		if err != nil {
			resp.Error = err
			return nil, err
		}
		resp.Data.Secret = gotp.NewDefaultTOTP(secret).ProvisioningUri(user.Username, "IDAS")
		resp.Data.Token = token.Id
		return resp, nil
	}
}

func MakeCreateTOTPEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateTOTPRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		user := ctx.Value(global.MetaUser).(*models.User)
		if user == nil {
			return nil, errors.NotLoginError()
		}
		var secret string
		if !s.VerifyToken(ctx, req.Token, models.TokenTypeTotpSecret, &secret) {
			resp.Error = errors.ParameterError("token")
			return resp, nil
		}

		nowTime := time.Now()
		ts := nowTime.Add(time.Second * time.Duration(-(nowTime.Second() % 30))).Unix()
		totp := gotp.NewDefaultTOTP(secret)
		if !totp.Verify(req.FirstCode, ts-30) {
			resp.Error = errors.NewServerError(http.StatusBadRequest, "The first code is invalid or expired")
		} else if !totp.Verify(req.SecondCode, ts) {
			resp.Error = errors.NewServerError(http.StatusBadRequest, "The second code is invalid or expired")
		} else {
			resp.Error = s.CreateTOTP(ctx, user.Id, secret)
			if resp.Error == nil {
				s.DeleteToken(ctx, models.TokenTypeTotpSecret, req.Token)
			}
		}

		return resp, nil
	}
}
