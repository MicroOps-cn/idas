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

	"github.com/go-kit/kit/endpoint"

	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
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
		req := request.(Requester).GetRequestData().(*DeleteRolesRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		var ids []string
		for _, role := range *req {
			ids = append(ids, role.Id)
		}
		resp.Error = s.DeleteRoles(ctx, ids)
		return &resp, nil
	}
}
