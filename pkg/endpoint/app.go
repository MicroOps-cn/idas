package endpoint

import (
	"context"
	"encoding/json"
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

type PatchAppsRequest struct {
	appPatch []map[string]interface{}
}

func (p *PatchAppsRequest) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &p.appPatch)
}

type PatchAppsResponse struct {
}

func MakePatchAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchAppsRequest)
		resp := BaseTotalResponse[PatchAppsResponse]{}
		var storage string
		for _, patch := range req.appPatch {
			if ss, ok := patch["storage"].(string); !ok || len(ss) == 0 {
				return nil, errors.ParameterError("storage is null")
			} else if patch["storage"] != storage && storage != "" {
				return nil, errors.ParameterError("storage is inconsistent")
			} else {
				storage = ss
			}
		}
		resp.Total, resp.Error = s.PatchApps(ctx, storage, req.appPatch)
		return &resp, nil
	}
}

type DeleteAppsRequest struct {
	Id      []string `valid:"required,notnull"`
	Storage string   `json:"storage" valid:"required"`
}

type DeleteAppsResponse struct {
}

func MakeDeleteAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteAppsRequest)
		resp := BaseTotalResponse[DeleteAppsResponse]{}
		resp.Total, resp.Error = s.DeleteApps(ctx, req.Storage, req.Id)
		return &resp, nil
	}
}

type UpdateAppRequest struct {
	*models.App `json:",inline"`
}

type UpdateAppResponse struct {
}

func MakeUpdateAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdateAppRequest)
		resp := BaseResponse[interface{}]{}
		if resp.Data, resp.Error = s.UpdateApp(ctx, req.Storage, req.App); resp.Error != nil {
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

type CreateAppRequest struct {
	Name        string            `json:"name" valid:"required"`
	Description string            `json:"description"`
	Avatar      string            `json:"avatar"`
	Storage     string            `json:"storage" valid:"required"`
	GrantType   models.GrantType  `json:"grantType" valid:"required"`
	GrantMode   models.GrantMode  `json:"grantMode"`
	User        []*models.User    `gorm:"many2many:app_user" json:"user,omitempty"`
	Role        []*models.AppRole `gorm:"-" json:"role,omitempty"`
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
			User:        req.User,
			Role:        req.Role,
		})
		return &resp, nil
	}
}

type PatchAppRequest struct {
	fields  map[string]interface{}
	Storage string `json:"storage" valid:"required"`
}

type PatchAppResponse struct {
}

func (p *PatchAppRequest) UnmarshalJSON(data []byte) error {
	fields := map[string]interface{}{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	p.fields = fields
	return nil
}

func MakePatchAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchAppRequest)
		resp := BaseResponse[*models.App]{}
		resp.Data, resp.Error = s.PatchApp(ctx, req.Storage, req.fields)
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
