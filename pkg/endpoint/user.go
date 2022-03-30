package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/service"
	"idas/pkg/service/models"
)

type UserLoginRequest struct {
	BaseRequest
	Username   string `json:"username,omitempty"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

type UserLoginResponse struct {
	BaseResponse `json:"-"`
}

func MakeUserLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*UserLoginRequest)
		resp := UserLoginResponse{}
		if loginCookie, err := s.CreateLoginSession(ctx, req.Username, req.Password); err == nil {
			req.restfulResponse.AddHeader("Set-Cookie", loginCookie)
		} else {
			resp.Error = errors.NewServerError(http.StatusUnauthorized, "Wrong user name or password")
		}
		return &resp, nil
	}
}

type UserLogoutRequest struct {
	BaseRequest
}

type UserLogoutResponse struct {
	BaseResponse `json:",inline"`
}

func MakeUserLogoutEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*UserLogoutRequest)
		resp := UserLogoutResponse{}
		cookie, err := req.GetRestfulRequest().Request.Cookie(global.LoginSession)
		if err != nil {
			resp.Error = errors.BadRequestError
		} else if len(cookie.Value) > 0 {
			if loginCookie, err := s.DeleteLoginSession(ctx, cookie.Value); err == nil {
				req.restfulResponse.AddHeader("Set-Cookie", fmt.Sprintf(loginCookie))
			} else {
				resp.Error = errors.InternalServerError
			}
		} else {
			resp.Error = errors.NewServerError(http.StatusUnauthorized, "Invalid identity information")
		}
		return &resp, nil
	}
}

type CurrentUserRequest struct {
	BaseRequest
}

type CurrentUserResponse struct {
	BaseResponse `json:"-"`
	*models.User
}

func MakeCurrentUserEndpoint(_ service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CurrentUserRequest)
		resp := CurrentUserResponse{}
		if user, ok := req.GetRestfulRequest().Attribute(global.AttrUser).(*models.User); ok {
			resp.User = user
		} else {
			resp.Error = errors.NotLoginError
		}
		return &resp, nil
	}
}

type GetUsersRequest struct {
	BaseListRequest
	Status  models.UserStatus `json:"status"`
	Storage string            `json:"storage"`
}

type GetUsersResponse struct {
	BaseListResponse `json:"-"`
}

func MakeGetUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*GetUsersRequest)
		resp := GetUsersResponse{BaseListResponse: NewBaseListResponse(req.BaseListRequest)}
		resp.Data, resp.Total, resp.Error = s.GetUsers(ctx, req.Storage, req.Keywords, req.Status, req.Current, req.PageSize)
		return &resp, nil
	}
}

type GetUserSourceRequest struct {
	BaseRequest
}

type GetUserSourceResponse struct {
	BaseListResponse `json:"-"`
}

func MakeGetUserSourceRequestEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := GetUserSourceResponse{}
		resp.Data, resp.Total, resp.Error = s.GetUserSource(ctx)
		return &resp, nil
	}
}

type PatchUsersRequest struct {
	BaseRequest
	userPatch []map[string]interface{}
	Storage   string `json:"storage" valid:"required"`
}

func (p *PatchUsersRequest) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &p.userPatch)
}

type PatchUsersResponse struct {
	BaseTotalResponse `json:"-"`
}

func MakePatchUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*PatchUsersRequest)
		resp := PatchUsersResponse{}
		resp.Total, resp.ErrorMessage, resp.Error = s.PatchUsers(ctx, req.Storage, req.userPatch)
		return &resp, nil
	}
}

type DeleteUsersRequest struct {
	BaseRequest
	Id      []string `valid:"required,notnull"`
	Storage string   `json:"storage" valid:"required"`
}

type DeleteUsersResponse struct {
	BaseTotalResponse `json:"-"`
}

func MakeDeleteUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*DeleteUsersRequest)
		resp := DeleteUsersResponse{}
		resp.Total, resp.ErrorMessage, resp.Error = s.DeleteUsers(ctx, req.Storage, req.Id)
		return &resp, nil
	}
}

type UpdateUserRequest struct {
	BaseRequest
	User    *models.User `json:",inline"`
	Storage string       `json:"storage" valid:"required"`
}

type UpdateUserResponse struct {
	BaseResponse `json:"-"`
	User         *models.User `json:",inline"`
}

func MakeUpdateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*UpdateUserRequest)
		resp := UpdateUserResponse{}
		if resp.User, resp.ErrorMessage, resp.Error = s.UpdateUser(ctx, req.Storage, req.User); resp.Error != nil {
			resp.Error = errors.NewServerError(200, resp.Error.Error())
		}
		return &resp, nil
	}
}

type GetUserRequest struct {
	BaseRequest
	Id       string
	Username string
	Storage  string `json:"storage" valid:"required"`
}

type GetUserResponse struct {
	BaseResponse `json:"-"`
	User         *models.User `json:",inline"`
}

func MakeGetUserInfoEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequest)
		resp := GetUserResponse{}
		resp.User, resp.ErrorMessage, resp.Error = s.GetUserInfo(ctx, req.Storage, req.Id, req.Username)
		return &resp, nil
	}
}

type CreateUserRequest struct {
	BaseRequest
	User    *models.User `json:",inline"`
	Storage string       `json:"storage" valid:"required"`
}

type CreateUserResponse struct {
	BaseResponse `json:"-"`
	User         *models.User `json:",inline"`
}

func MakeCreateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CreateUserRequest)
		resp := CreateUserResponse{}
		resp.User, resp.ErrorMessage, resp.Error = s.CreateUser(ctx, req.Storage, req.User)
		return &resp, nil
	}
}

type PatchUserRequest struct {
	BaseRequest
	fields  map[string]interface{}
	Storage string `json:"storage" valid:"required"`
}

type PatchUserResponse struct {
	BaseResponse `json:"-"`
	User         *models.User `json:",inline"`
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
		req := request.(*PatchUserRequest)
		resp := PatchUserResponse{}
		resp.User, resp.ErrorMessage, resp.Error = s.PatchUser(ctx, req.Storage, req.fields)
		return &resp, nil
	}
}

type DeleteUserRequest struct {
	BaseRequest
	Id      string `valid:"required"`
	Storage string `json:"storage" valid:"required"`
}

type DeleteUserResponse struct {
	BaseResponse `json:"-"`
}

func MakeDeleteUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*DeleteUserRequest)
		resp := DeleteUserResponse{}
		resp.ErrorMessage, resp.Error = s.DeleteUser(ctx, req.Storage, req.Id)
		return &resp, nil
	}
}

type GetLoginSession struct {
	BaseResponse `json:"-"`
	User         *models.User `json:",inline"`
}

func MakeGetLoginSessionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		sessionId := request.(string)
		var resp *models.User
		if sessionId != "" {
			if resp, _, err = s.GetLoginSession(ctx, sessionId); err != nil {
				err = errors.NotLoginError
			}
		} else {
			err = errors.NotLoginError
		}
		return resp, err
	}
}
