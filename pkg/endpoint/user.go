package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
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

type GetUsersRequest struct {
	BaseListRequest
	Status  models.UserStatus `json:"status"`
	Storage string            `json:"storage"`
	App     string            `json:"app"`
}

func MakeGetUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetUsersRequest)
		resp := NewBaseListResponse[interface{}](&req.BaseListRequest)
		resp.BaseResponse.Data, resp.Total, resp.BaseResponse.Error = s.GetUsers(ctx, req.Storage, req.Keywords, req.Status, req.App, req.Current, req.PageSize)
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

type PatchUsersRequest struct {
	userPatch []map[string]interface{}
}

func (p *PatchUsersRequest) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &p.userPatch)
}

type PatchUsersResponse struct {
}

func MakePatchUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchUsersRequest)
		resp := BaseTotalResponse[interface{}]{}
		var storage string
		for _, patch := range req.userPatch {
			if ss, ok := patch["storage"].(string); !ok || len(ss) == 0 {
				return nil, errors.ParameterError("storage is null")
			} else if patch["storage"] != storage && storage != "" {
				return nil, errors.ParameterError("storage is inconsistent")
			} else {
				storage = ss
			}
		}
		resp.Total, resp.Error = s.PatchUsers(ctx, storage, req.userPatch)
		return &resp, nil
	}
}

type DeleteUsersRequest struct {
	Id      []string `json:"id" valid:"required"`
	Storage string   `json:"storage" valid:"required"`
}

type DeleteUsersResponse struct {
}

func MakeDeleteUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteUsersRequest)
		resp := BaseTotalResponse[interface{}]{}
		resp.Total, resp.Error = s.DeleteUsers(ctx, req.Storage, req.Id)
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
	Username    string `gorm:"type:varchar(20);" json:"username"`
	Email       string `gorm:"type:varchar(50);" json:"email" `
	PhoneNumber string `json:"phoneNumber"`
	FullName    string `gorm:"type:varchar(20);" json:"fullName"`
	Avatar      string `json:"avatar"`
	Storage     string `gorm:"-" json:"storage"`
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
	fields  map[string]interface{}
	Storage string `json:"storage" valid:"required"`
}

type PatchUserResponse struct {
	User *models.User `json:",inline"`
}

func (p *PatchUserRequest) UnmarshalJSON(data []byte) error {
	fields := map[string]interface{}{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	p.fields = fields
	return nil
}

func MakePatchUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchUserRequest)
		resp := BaseResponse[interface{}]{}
		resp.Data, resp.Error = s.PatchUser(ctx, req.Storage, req.fields)
		return &resp, nil
	}
}

type DeleteUserRequest struct {
	Id      string `valid:"required"`
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
