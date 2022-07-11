package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"

	"idas/pkg/errors"
	"idas/pkg/service"
	"idas/pkg/service/models"
)

type GetAppsRequest struct {
	BaseListRequest
	Storage string `json:"storage"`
}

func MakeGetAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetAppsRequest)
		resp := NewBaseListResponse[[]*models.App](&req.BaseListRequest)
		resp.BaseResponse.Data, resp.Total, resp.BaseResponse.Error = s.GetApps(ctx, req.Storage, req.Keywords, req.Current, req.PageSize)
		return &resp, nil
	}
}

type GetAppSourceRequest struct {
}

type GetAppSourceResponse map[string]string

func MakeGetAppSourceRequestEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := BaseTotalResponse[GetAppSourceResponse]{}
		resp.Data, resp.Total, resp.Error = s.GetAppSource(ctx)
		return &resp, nil
	}
}

type PatchAppsRequest []PatchAppRequest

type PatchAppsResponse struct {
}

func MakePatchAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchAppsRequest)
		resp := BaseTotalResponse[PatchAppsResponse]{}
		var patchApps = map[string][]map[string]interface{}{}
		for _, a := range *req {
			if len(a.Storage) == 0 {
				return nil, errors.ParameterError("There is an empty storage in the patch.")
			}
			if len(a.Id) == 0 {
				return nil, errors.ParameterError("There is an empty id in the patch.")
			}
			var patch = map[string]interface{}{"id": a.Id}
			if a.Status != nil {
				patch["status"] = *a.Status
			}
			if a.IsDelete != nil {
				patch["isDelete"] = *a.IsDelete
			}
			patchApps[a.Storage] = append(patchApps[a.Storage], patch)
		}

		errs := errors.NewMultipleServerError(500, "Multiple errors have occurred: ")
		for storage, patch := range patchApps {
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

type DeleteAppsRequest []DeleteAppRequest

type DeleteAppsResponse struct {
}

func MakeDeleteAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteAppsRequest)
		resp := BaseTotalResponse[DeleteAppsResponse]{}
		var delApps = map[string][]string{}
		for _, app := range *req {
			if len(app.Storage) == 0 {
				return nil, errors.ParameterError("There is an empty storage in the request.")
			}
			if len(app.Id) == 0 {
				return nil, errors.ParameterError("There is an empty id in the request.")
			}
			delApps[app.Storage] = append(delApps[app.Storage], app.Id)
		}
		errs := errors.NewMultipleServerError(500, "Multiple errors have occurred: ")
		for storage, ids := range delApps {
			total, err := s.DeleteApps(ctx, storage, ids)
			resp.Total += total
			if err != nil {
				errs.Append(err)
				resp.Error = err
			}
		}
		return &resp, nil
	}
}

type UpdateAppRequest struct {
	Id          string             `json:"id" valid:"required"`
	Name        string             `json:"name" valid:"required"`
	Description string             `json:"description,omitempty"`
	Avatar      string             `json:"avatar,omitempty"`
	Storage     string             `json:"storage" valid:"required"`
	GrantType   models.GrantType   `json:"grantType" valid:"required"`
	GrantMode   models.GrantMode   `json:"grantMode" valid:"required"`
	Status      models.GroupStatus `json:"status,omitempty"`
	User        []AppUser          `json:"user,omitempty"`
	Role        []AppRole          `json:"role,omitempty"`
}

func (r UpdateAppRequest) GetUsers() (users []*models.User) {
	for _, u := range r.User {
		users = append(users, &models.User{Model: models.Model{Id: u.Id}, RoleId: u.RoleId})
	}
	return users
}

func (r UpdateAppRequest) GetRoles() (roles []*models.AppRole) {
	for _, role := range r.Role {
		roles = append(roles, &models.AppRole{Model: models.Model{Id: role.Id}, Name: role.Name, Config: role.Config, IsDefault: role.IsDefault})
	}
	return roles
}

type UpdateAppResponse struct {
}

func MakeUpdateAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdateAppRequest)
		resp := BaseResponse[interface{}]{}
		if resp.Data, resp.Error = s.UpdateApp(ctx, req.Storage, &models.App{
			Model:       models.Model{Id: req.Id},
			Name:        req.Name,
			Description: req.Description,
			Avatar:      req.Avatar,
			GrantType:   req.GrantType,
			GrantMode:   req.GrantMode,
			Storage:     req.Storage,
			User:        req.GetUsers(),
			Role:        req.GetRoles(),
		}); resp.Error != nil {
			resp.Error = errors.NewServerError(200, resp.Error.Error())
		}
		return &resp, nil
	}
}

type GetAppRequest struct {
	Id      string `json:"id" valid:"required"`
	Storage string `json:"storage" valid:"required"`
}

type GetAppResponse struct {
}

func MakeGetAppInfoEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetAppRequest)
		resp := BaseResponse[*models.App]{}
		resp.Data, resp.Error = s.GetAppInfo(ctx, req.Storage, req.Id)
		return &resp, nil
	}
}

type AppUser struct {
	Id     string `json:"id" valid:"required"`
	RoleId string `json:"roleId,omitempty"`
}

type AppRole struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name" valid:"required"`
	Config    string `json:"config,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty"`
}
type CreateAppRequest struct {
	Name        string           `json:"name" valid:"required"`
	Description string           `json:"description"`
	Avatar      string           `json:"avatar"`
	Storage     string           `json:"storage" valid:"required"`
	GrantType   models.GrantType `json:"grantType" valid:"required"`
	GrantMode   models.GrantMode `json:"grantMode" valid:"required"`
	User        []AppUser        `json:"user,omitempty"`
	Role        []AppRole        `json:"role,omitempty"`
}

func (r CreateAppRequest) GetUsers() (users []*models.User) {
	for _, u := range r.User {
		users = append(users, &models.User{Model: models.Model{Id: u.Id}, RoleId: u.RoleId})
	}
	return users
}

func (r CreateAppRequest) GetRoles() (roles []*models.AppRole) {
	for _, role := range r.Role {
		roles = append(roles, &models.AppRole{Model: models.Model{Id: role.Id}, Name: role.Name, Config: role.Config, IsDefault: role.IsDefault})
	}
	return roles
}

type CreateAppResponse struct {
	App *models.App `json:",inline"`
}

func MakeCreateAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateAppRequest)
		resp := BaseResponse[*models.App]{}
		resp.Data, resp.Error = s.CreateApp(ctx, req.Storage, &models.App{
			Name:        req.Name,
			Description: req.Description,
			Avatar:      req.Avatar,
			GrantType:   req.GrantType,
			GrantMode:   req.GrantMode,
			Storage:     req.Storage,
			User:        req.GetUsers(),
			Role:        req.GetRoles(),
		})
		return &resp, nil
	}
}

type PatchAppRequest struct {
	Id          string              `json:"id" valid:"required"`
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Avatar      *string             `json:"avatar,omitempty"`
	GrantType   *models.GrantType   `json:"grantType,omitempty"`
	GrantMode   *models.GrantMode   `json:"grantMode,omitempty"`
	Status      *models.GroupStatus `json:"status,omitempty"`
	Storage     string              `json:"storage" valid:"required"`
	IsDelete    *bool               `json:"isDelete,omitempty"`
}

type PatchAppResponse struct {
}

func MakePatchAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchAppRequest)
		resp := BaseResponse[*models.App]{}

		if len(req.Storage) == 0 {
			return nil, errors.ParameterError("There is an empty storage in the patch.")
		}
		if len(req.Id) == 0 {
			return nil, errors.ParameterError("There is an empty id in the patch.")
		}
		var tmpPatch = map[string]interface{}{
			"id":          req.Id,
			"storage":     req.Storage,
			"name":        req.Name,
			"description": req.Description,
			"avatar":      req.Avatar,
			"grantType":   req.GrantType,
			"grantMode":   req.GrantMode,
			"status":      req.Status,
			"isDelete":    req.IsDelete,
		}
		var patch = map[string]interface{}{}
		for name, val := range tmpPatch {
			if val != nil {
				patch[name] = val
			}
		}
		resp.Data, resp.Error = s.PatchApp(ctx, req.Storage, patch)
		return &resp, nil
	}
}

type DeleteAppRequest struct {
	Id      string `valid:"required"`
	Storage string `json:"storage" valid:"required"`
}

type DeleteAppResponse struct {
}

func MakeDeleteAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteAppRequest)
		resp := BaseResponse[*DeleteAppResponse]{}
		resp.Error = s.DeleteApp(ctx, req.Storage, req.Id)
		return &resp, nil
	}
}
