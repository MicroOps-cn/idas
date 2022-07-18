package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"

	"idas/pkg/errors"
	"idas/pkg/service"
	"idas/pkg/service/models"
)

func MakeGetAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetAppsRequest)
		resp := NewBaseListResponse[[]*models.App](&req.BaseListRequest)
		resp.Total, resp.Data, resp.BaseResponse.Error = s.GetApps(ctx, req.Storage, req.Keywords, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeGetAppSourceRequestEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := TotalResponseWrapper[map[string]string]{}
		resp.Total, resp.Data, resp.Error = s.GetAppSource(ctx)
		return &resp, nil
	}
}

type PatchAppsRequest []PatchAppRequest

func MakePatchAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchAppsRequest)
		resp := TotalResponseWrapper[interface{}]{}
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
				patch["status"] = a.Status
			}
			if a.IsDelete != nil {
				patch["isDelete"] = a.IsDelete
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

func MakeDeleteAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteAppsRequest)
		resp := TotalResponseWrapper[interface{}]{}
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

func MakeUpdateAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdateAppRequest)
		resp := SimpleResponseWrapper[*models.App]{}
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

func MakeGetAppInfoEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetAppRequest)
		resp := SimpleResponseWrapper[*models.App]{}
		resp.Data, resp.Error = s.GetAppInfo(ctx, req.Storage, req.Id)
		return &resp, nil
	}
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

func MakeCreateAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateAppRequest)
		resp := SimpleResponseWrapper[interface{}]{}
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

func MakePatchAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchAppRequest)
		resp := SimpleResponseWrapper[interface{}]{}

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

func MakeDeleteAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteAppRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		resp.Error = s.DeleteApp(ctx, req.Storage, req.Id)
		return &resp, nil
	}
}
