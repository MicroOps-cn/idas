package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"idas/pkg/service"
	"idas/pkg/service/models"
)

type DeleteRolesRequest []DeleteRoleRequest

func MakeGetPermissionsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetPermissionsRequest)
		resp := NewBaseListResponse[[]*models.Permission](&req.BaseListRequest)
		resp.Total, resp.Data, resp.BaseResponse.Error = s.GetPermissions(ctx, req.Keywords, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeGetRolesEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetRolesRequest)
		resp := NewBaseListResponse[[]*models.Role](&req.BaseListRequest)
		resp.Total, resp.Data, resp.BaseResponse.Error = s.GetRoles(ctx, req.Keywords, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeCreateRoleEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateRoleRequest)
		resp := SimpleResponseWrapper[*models.Role]{}
		role := &models.Role{
			Name:        req.Name,
			Description: req.Description,
		}
		for _, pid := range req.Permission {
			role.Permission = append(role.Permission, &models.Permission{Model: models.Model{Id: pid}})
		}
		resp.Data, resp.Error = s.CreateRole(ctx, role)
		return &resp, nil
	}
}

func MakeUpdateRoleEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdateRoleRequest)
		resp := SimpleResponseWrapper[*models.Role]{}
		role := &models.Role{
			Model:       models.Model{Id: req.Id},
			Name:        req.Name,
			Description: req.Description,
		}
		for _, pid := range req.Permission {
			role.Permission = append(role.Permission, &models.Permission{Model: models.Model{Id: pid}})
		}
		resp.Data, resp.Error = s.UpdateRole(ctx, role)
		return &resp, nil
	}
}
func MakeDeleteRoleEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteRoleRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		resp.Error = s.DeleteRoles(ctx, []string{req.Id})
		return &resp, nil
	}
}
func MakeDeleteRolesEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(DeleteRolesRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		var ids []string
		for _, role := range req {
			ids = append(ids, role.Id)
		}
		resp.Error = s.DeleteRoles(ctx, ids)
		return &resp, nil
	}
}
