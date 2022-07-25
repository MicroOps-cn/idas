package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"

	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/logs"
	"idas/pkg/service"
	"idas/pkg/service/models"
	w "idas/pkg/utils/wrapper"
)

func MakeCurrentUserEndpoint(_ service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := SimpleResponseWrapper[interface{}]{}
		if users, ok := ctx.Value(global.MetaUser).([]*models.User); ok && len(users) > 0 {
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

func MakeResetUserPasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*ResetUserPasswordRequest)
		resp := TotalResponseWrapper[interface{}]{}
		if len(req.Token) > 0 {
			if s.VerifyToken(ctx, req.Token, req.UserId, models.TokenTypeResetPassword) {
				resp.Error = s.ResetPassword(ctx, req.UserId, req.Storage, req.NewPassword)
			}
		} else {
			// implement me
		}
		// if auth, ok := req.Auth.(*ResetUserPasswordRequest_Token); ok && len(auth.Token) > 0 {
		// 	if s.VerifyToken(ctx, auth.Token, req.UserId, models.TokenTypeResetPassword) {
		// 		resp.Error = s.ResetPassword(ctx, req.UserId, req.Storage, req.NewPassword)
		// 	}
		// } else {
		//
		// }

		return resp, nil
	}
}

func MakeForgotPasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*ForgotUserPasswordRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		users := s.GetUserInfoByUsernameAndEmail(ctx, req.Username, req.Email)
		if len(users) == 0 {
			level.Warn(logs.GetContextLogger(ctx)).Log("err", "user not found", "msg", "failed to reset password", "username", req.Username, "email", req.Email)
			return resp, nil
		}

		token, err := s.CreateToken(ctx, models.TokenTypeResetPassword, w.ToInterfaces[*models.User](users)...)
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
		resp := NewBaseListResponse[[]*models.User](&req.BaseListRequest)
		resp.Total, resp.Data, resp.Error = s.GetUsers(
			ctx, req.Storage, req.Keywords,
			models.UserStatus(req.Status),
			req.App, req.Current, req.PageSize,
		)
		return &resp, nil
	}
}

func MakeGetUserSourceRequestEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := TotalResponseWrapper[map[string]string]{}
		resp.Total, resp.Data, resp.BaseResponse.Error = s.GetUserSource(ctx)
		return &resp, nil
	}
}

type PatchUsersRequest []PatchUserRequest

func MakePatchUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchUsersRequest)
		resp := TotalResponseWrapper[interface{}]{}

		patchUsers := map[string][]map[string]interface{}{}
		for _, u := range *req {
			if len(u.Storage) == 0 {
				return nil, errors.ParameterError("There is an empty storage in the patch.")
			}
			if len(u.Id) == 0 {
				return nil, errors.ParameterError("There is an empty id in the patch.")
			}
			patch := map[string]interface{}{"id": u.Id}
			if u.Status != nil {
				patch["status"] = u.Status
			}
			if u.IsDelete != nil {
				patch["isDelete"] = u.IsDelete
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

func MakeDeleteUsersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteUsersRequest)
		resp := TotalResponseWrapper[interface{}]{}
		delUsers := map[string][]string{}
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

func MakeUpdateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdateUserRequest)
		resp := SimpleResponseWrapper[*models.User]{}
		if resp.Data, resp.Error = s.UpdateUser(ctx, req.Storage, &models.User{
			Model: models.Model{
				Id:       req.Id,
				IsDelete: req.IsDelete,
			},
			Username:    req.Username,
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			FullName:    req.FullName,
			Avatar:      req.Avatar,
			Status:      models.UserStatus(req.Status),
			Storage:     req.Storage,
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
		resp.Data, resp.Error = s.GetUserInfo(ctx, req.Storage, req.Id, req.Username)
		return &resp, nil
	}
}

func MakeCreateUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateUserRequest)
		resp := SimpleResponseWrapper[*models.User]{}
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

type PatchUserResponse struct {
	User *models.User `json:",inline"`
}

func MakePatchUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchUserRequest)
		resp := SimpleResponseWrapper[*models.User]{}
		if len(req.Storage) == 0 {
			return nil, errors.ParameterError("There is an empty storage in the patch.")
		}
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
		resp.Data, resp.Error = s.PatchUser(ctx, req.Storage, patch)
		return &resp, nil
	}
}

func MakeDeleteUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteUserRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		resp.Error = s.DeleteUser(ctx, req.Storage, req.Id)
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

func MakeCreateKeyEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateKeyRequest)
		resp := SimpleResponseWrapper[*models.UserKey]{}
		users, ok := ctx.Value(global.MetaUser).([]*models.User)
		if !ok || len(users) == 0 {
			return nil, errors.NotLoginError
		}
		for _, user := range users {
			if user.Id == req.UserId {
				resp.Data, resp.Error = s.CreateUserKey(ctx, req.UserId, req.Name)
				return &resp, nil
			}
		}
		return nil, errors.StatusNotFound("user")
	}
}
