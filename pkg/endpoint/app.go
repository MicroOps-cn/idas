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

type GetAppsResponse struct {
	BaseListResponse `json:"-"`
}

func MakeGetAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*GetAppsRequest)
		resp := GetAppsResponse{BaseListResponse: NewBaseListResponse(req.BaseListRequest)}
		resp.Data, resp.Total, resp.Error = s.GetApps(ctx, req.Storage, req.Keywords, req.Current, req.PageSize)
		return &resp, nil
	}
}

type GetAppSourceRequest struct {
	BaseRequest
}

type GetAppSourceResponse struct {
	BaseListResponse `json:"-"`
}

func MakeGetAppSourceRequestEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := GetAppSourceResponse{}
		resp.Data, resp.Total, resp.Error = s.GetAppSource(ctx)
		return &resp, nil
	}
}

type PatchAppsRequest struct {
	BaseRequest
	appPatch []map[string]interface{}
}

func (p *PatchAppsRequest) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &p.appPatch)
}

type PatchAppsResponse struct {
	BaseTotalResponse `json:"-"`
}

func MakePatchAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*PatchAppsRequest)
		resp := PatchAppsResponse{}
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
	BaseRequest
	Id      []string `valid:"required,notnull"`
	Storage string   `json:"storage" valid:"required"`
}

type DeleteAppsResponse struct {
	BaseTotalResponse `json:"-"`
}

func MakeDeleteAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*DeleteAppsRequest)
		resp := DeleteAppsResponse{}
		resp.Total, resp.Error = s.DeleteApps(ctx, req.Storage, req.Id)
		return &resp, nil
	}
}

type UpdateAppRequest struct {
	BaseRequest
	*models.App `json:",inline"`
}

type UpdateAppResponse struct {
	BaseResponse `json:"-"`
	App          *models.App `json:",inline"`
}

func MakeUpdateAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*UpdateAppRequest)
		resp := UpdateAppResponse{}
		if resp.App, resp.Error = s.UpdateApp(ctx, req.Storage, req.App); resp.Error != nil {
			resp.Error = errors.NewServerError(200, resp.Error.Error())
		}
		return &resp, nil
	}
}

type GetAppRequest struct {
	BaseRequest
	Id      string
	Storage string `json:"storage" valid:"required"`
}

type GetAppResponse struct {
	BaseResponse `json:"-"`
	App          *models.App `json:",inline"`
}

func MakeGetAppInfoEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetAppRequest)
		resp := GetAppResponse{}
		resp.App, resp.Error = s.GetAppInfo(ctx, req.Storage, req.Id)
		return &resp, nil
	}
}

type CreateAppRequest struct {
	BaseRequest
	Name        string           `json:"name" valid:"required"`
	Description string           `json:"description"`
	Avatar      string           `json:"avatar"`
	Storage     string           `json:"storage" valid:"required"`
	GrantType   models.GrantType `json:"grantType" valid:"required"`
	GrantMode   models.GrantMode `json:"grantMode"`
}

type CreateAppResponse struct {
	BaseResponse `json:"-"`
	App          *models.App `json:",inline"`
}

func MakeCreateAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CreateAppRequest)
		resp := CreateAppResponse{}
		resp.App, resp.Error = s.CreateApp(ctx, req.Storage, &models.App{
			Name:        req.Name,
			Description: req.Description,
			Avatar:      req.Avatar,
			GrantType:   req.GrantType,
			GrantMode:   req.GrantMode,
			Storage:     req.Storage,
		})
		return &resp, nil
	}
}

type PatchAppRequest struct {
	BaseRequest
	fields  map[string]interface{}
	Storage string `json:"storage" valid:"required"`
}

type PatchAppResponse struct {
	BaseResponse `json:"-"`
	App          *models.App `json:",inline"`
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
		req := request.(*PatchAppRequest)
		resp := PatchAppResponse{}
		resp.App, resp.Error = s.PatchApp(ctx, req.Storage, req.fields)
		return &resp, nil
	}
}

type DeleteAppRequest struct {
	BaseRequest
	Id      string `valid:"required"`
	Storage string `json:"storage" valid:"required"`
}

type DeleteAppResponse struct {
	BaseResponse `json:"-"`
}

func MakeDeleteAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*DeleteAppRequest)
		resp := DeleteAppResponse{}
		resp.Error = s.DeleteApp(ctx, req.Storage, req.Id)
		return &resp, nil
	}
}
