package endpoint

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"
	"idas/pkg/logs"
	"net/http"
	"strings"
	"time"

	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/service"
	"idas/pkg/service/models"
)

type UserLoginRequest struct {
	Username   string `json:"username,omitempty"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

type UserLoginResponse struct {
}

func MakeUserLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UserLoginRequest)
		resp := BaseResponse[*UserLoginResponse]{}
		if loginCookie, err := s.CreateLoginSession(ctx, req.Username, req.Password); err == nil {
			request.(RestfulRequester).GetRestfulResponse().AddHeader("Set-Cookie", strings.Join(loginCookie, ","))
		} else {
			resp.Error = errors.NewServerError(http.StatusUnauthorized, "Wrong user name or password")
		}
		return &resp, nil
	}
}

type UserLogoutRequest struct {
}

type UserLogoutResponse struct {
}

func MakeUserLogoutEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//req := request.(Requester).GetRequestData().(*UserLogoutRequest)
		resp := BaseResponse[*UserLogoutResponse]{}
		cookie, err := request.(RestfulRequester).GetRestfulRequest().Request.Cookie(global.LoginSession)
		if err != nil {
			resp.Error = errors.BadRequestError
		} else if len(cookie.Value) > 0 {
			for _, id := range strings.Split(cookie.Value, ",") {
				if err = s.DeleteLoginSession(ctx, id); err != nil {
					resp.Error = errors.InternalServerError
					return resp, nil
				}
			}
			loginCookie := fmt.Sprintf("%s=%s; Path=/;Expires=%s", global.LoginSession, cookie.Value, time.Now().UTC().Format(global.LoginSessionExpiresFormat))
			request.(RestfulRequester).GetRestfulResponse().AddHeader("Set-Cookie", fmt.Sprintf(loginCookie))
		} else {
			resp.Error = errors.NewServerError(http.StatusUnauthorized, "Invalid identity information")
		}
		return &resp, nil
	}
}

type CurrentUserRequest struct {
}

type CurrentUserResponse struct {
}

func MakeCurrentUserEndpoint(_ service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := BaseResponse[interface{}]{}
		if users, ok := request.(RestfulRequester).GetRestfulRequest().Attribute(global.AttrUser).([]*models.User); ok && len(users) > 0 {
			resp.Data = users
			for _, user := range users {
				return user, nil
			}
		} else {
			resp.Error = errors.NotLoginError
		}
		return resp, nil
	}
}

type ResetUserPasswordRequest struct {
	UserId      string `json:"userId" valid:"required"`
	Storage     string `json:"storage" valid:"required"`
	Token       string `json:"token,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword" valid:"required"`
}

type ResetUserPasswordResponse struct {
}

func MakeResetUserPasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*ResetUserPasswordRequest)
		resp := BaseTotalResponse[interface{}]{}
		if len(req.Token) > 0 {
			if s.VerifyToken(ctx, req.Token, req.UserId, models.TokenTypeResetPassword) {
				resp.Error = s.ResetPassword(ctx, req.UserId, req.Storage, req.NewPassword)
			}
		} else {

		}

		return resp, nil
	}
}

type ForgotUserPasswordRequest struct {
	Username string `json:"username" valid:"required"`
	Email    string `json:"email" valid:"required"`
}

func MakeForgotPasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*ForgotUserPasswordRequest)
		resp := BaseResponse[interface{}]{}
		users := s.GetUserInfoByUsernameAndEmail(ctx, req.Username, req.Email)
		if len(users) == 0 {
			return resp, nil
		}
		token, err := s.CreateToken(ctx, users, models.TokenTypeResetPassword)
		if err != nil {
			return resp, err
		}
		err = s.SendResetPasswordLink(ctx, users, token)
		if err != nil {
			level.Error(logs.GetContextLogger(ctx)).Log("err", err, "msg", "failed to send email")
		}
		return resp, nil
	}
}

func MakeGetUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetUsersRequest)
		resp := NewBaseListResponse[interface{}](&req.BaseListRequest)
		resp.BaseResponse.Data, resp.Total, resp.BaseResponse.Error = s.GetUsers(
			ctx, req.Storage, req.Keywords,
			models.UserStatus(req.Status),
			req.App, req.Current, req.PageSize,
		)
		return &resp, nil
	}
}

type GetUserSourceRequest struct {
}

type GetUserSourceResponse map[string]string

func MakeGetUserSourceRequestEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := BaseListResponse[GetUserSourceResponse]{}
		resp.BaseResponse.Data, resp.Total, resp.BaseResponse.Error = s.GetUserSource(ctx)
		return &resp, nil
	}
}

type PatchUsersRequest []PatchUserRequest

type PatchUsersResponse struct {
}

func MakePatchUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchUsersRequest)
		resp := BaseTotalResponse[interface{}]{}

		var patchUsers = map[string][]map[string]interface{}{}
		for _, u := range *req {
			if len(u.Storage) == 0 {
				return nil, errors.ParameterError("There is an empty storage in the patch.")
			}
			if len(u.Id) == 0 {
				return nil, errors.ParameterError("There is an empty id in the patch.")
			}
			var patch = map[string]interface{}{"id": u.Id}
			if u.Status != nil {
				patch["status"] = *u.Status
			}
			if u.IsDelete != nil {
				patch["isDelete"] = *u.IsDelete
			}
			patchUsers[u.Storage] = append(patchUsers[u.Storage], patch)
		}
		errs := errors.NewMultipleServerError(500, "Multiple errors have occurred: ")
		for storage, patch := range patchUsers {
			total, err := s.PatchUsers(ctx, storage, patch)
			resp.Total += total
			if err != nil {
				errs.Append(err)
				resp.Error = err
			}
		}

		return &resp, nil
	}
}

type DeleteUsersRequest []DeleteUserRequest

type DeleteUsersResponse struct {
}

func MakeDeleteUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteUsersRequest)
		resp := BaseTotalResponse[interface{}]{}
		var delUsers = map[string][]string{}
		for _, u := range *req {
			if len(u.Storage) == 0 {
				return nil, errors.ParameterError("There is an empty storage in the request.")
			}
			if len(u.Id) == 0 {
				return nil, errors.ParameterError("There is an empty id in the request.")
			}
			delUsers[u.Storage] = append(delUsers[u.Storage], u.Id)
		}
		errs := errors.NewMultipleServerError(500, "Multiple errors have occurred: ")
		for storage, ids := range delUsers {
			total, err := s.DeleteUsers(ctx, storage, ids)
			resp.Total += total
			if err != nil {
				errs.Append(err)
				resp.Error = err
			}
		}
		return &resp, nil
	}
}

type UpdateUserRequest struct {
	Id          string            `json:"id" valid:"required"`
	Username    string            `json:"username"`
	Email       string            `json:"email" valid:"email,optional"`
	PhoneNumber string            `json:"phoneNumber" valid:"numeric,optional"`
	FullName    string            `json:"fullName"`
	Avatar      string            `json:"avatar"`
	Status      models.UserStatus `json:"status"`
	Storage     string            `json:"storage"`
}

func MakeUpdateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdateUserRequest)
		resp := BaseResponse[interface{}]{}
		if resp.Data, resp.Error = s.UpdateUser(ctx, req.Storage, &models.User{
			Model: models.Model{
				Id: req.Id,
			},
			Username:    req.Username,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			FullName:    req.FullName,
			Avatar:      req.Avatar,
			Status:      req.Status,
			Storage:     req.Storage,
		}); resp.Error != nil {
			resp.Error = errors.NewServerError(200, resp.Error.Error())
		}
		return &resp, nil
	}
}

type GetUserRequest struct {
	Id       string
	Username string
	Storage  string `json:"storage" valid:"required"`
}

type GetUserResponse struct {
	User *models.User `json:",inline"`
}

func MakeGetUserInfoEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequest)
		resp := BaseResponse[interface{}]{}
		resp.Data, resp.Error = s.GetUserInfo(ctx, req.Storage, req.Id, req.Username)
		return &resp, nil
	}
}

type CreateUserRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email" `
	PhoneNumber string `json:"phoneNumber,omitempty"`
	FullName    string `json:"fullName,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Storage     string `json:"storage"`
}

func MakeCreateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateUserRequest)
		resp := BaseResponse[interface{}]{}
		resp.Data, resp.Error = s.CreateUser(ctx, req.Storage, &models.User{
			Username:    req.Username,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			FullName:    req.FullName,
			Avatar:      req.Avatar,
			Storage:     req.Storage,
			Status:      models.UserStatusNormal,
		})
		return &resp, nil
	}
}

type PatchUserRequest struct {
	Id       string             `json:"id" valid:"required"`
	Storage  string             `json:"storage" valid:"required"`
	IsDelete *bool              `json:"isDelete,omitempty"`
	Status   *models.UserStatus `json:"status,omitempty"`
}

type PatchUserResponse struct {
	User *models.User `json:",inline"`
}

func MakePatchUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchUserRequest)
		resp := BaseResponse[interface{}]{}
		if len(req.Storage) == 0 {
			return nil, errors.ParameterError("There is an empty storage in the patch.")
		}
		if len(req.Id) == 0 {
			return nil, errors.ParameterError("There is an empty id in the patch.")
		}
		var patch = map[string]interface{}{"id": req.Id}
		if req.Status != nil {
			patch["status"] = *req.Status
		}
		if req.IsDelete != nil {
			patch["isDelete"] = *req.IsDelete
		}
		resp.Data, resp.Error = s.PatchUser(ctx, req.Storage, patch)
		return &resp, nil
	}
}

type DeleteUserRequest struct {
	Id      string `json:"id" valid:"required"`
	Storage string `json:"storage" valid:"required"`
}

type DeleteUserResponse struct {
}

func MakeDeleteUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteUserRequest)
		resp := BaseResponse[interface{}]{}
		resp.Error = s.DeleteUser(ctx, req.Storage, req.Id)
		return &resp, nil
	}
}

type GetLoginSession struct {
	User *models.User `json:",inline"`
}

func MakeGetLoginSessionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		sessionId := request.([]string)
		var resp []*models.User
		if len(sessionId) > 0 {
			if resp, err = s.GetLoginSession(ctx, sessionId); err != nil {
				err = errors.NotLoginError
			}
		} else {
			err = errors.NotLoginError
		}
		return resp, err
	}
}

type GetSessionsRequest struct {
	BaseListRequest
	UserId string `json:"userId" valid:"required"`
}

type GetSessionsResponse struct {
}

func MakeGetSessionsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetSessionsRequest)
		resp := NewBaseListResponse[interface{}](&req.BaseListRequest)
		resp.BaseResponse.Data, resp.Total, resp.BaseResponse.Error = s.GetSessions(ctx, req.UserId, req.Current, req.PageSize)
		return &resp, nil
	}
}

type DeleteSessionRequest struct {
	Id string `valid:"required"`
}

type DeleteSessionResponse struct {
}

func MakeDeleteSessionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteSessionRequest)
		resp := BaseResponse[interface{}]{}
		resp.Error = s.DeleteSession(ctx, req.Id)
		return &resp, nil
	}
}
