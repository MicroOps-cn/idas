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

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
)

func MakeGetAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetAppsRequest)
		resp := NewBaseListResponse[[]*models.App](&req.BaseListRequest)
		resp.Total, resp.Data, resp.BaseResponse.Error = s.GetApps(ctx, req.Keywords, nil, req.Current, req.PageSize)
		return &resp, nil
	}
}

type PatchAppsRequest []PatchAppRequest

func MakePatchAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchAppsRequest)
		resp := TotalResponseWrapper[interface{}]{}
		var patchApps []map[string]interface{}
		for _, a := range *req {
			if len(a.Id) == 0 {
				return nil, errors.ParameterError("There is an empty id in the patch.")
			}
			patch := map[string]interface{}{"id": a.Id}
			if a.Status != nil {
				patch["status"] = a.Status
			}
			if a.IsDelete != nil {
				patch["isDelete"] = a.IsDelete
			}
			patchApps = append(patchApps, patch)
		}

		resp.Total, resp.Error = s.PatchUsers(ctx, patchApps)

		return &resp, nil
	}
}

type DeleteAppsRequest []DeleteAppRequest

func MakeDeleteAppsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteAppsRequest)
		resp := TotalResponseWrapper[interface{}]{}
		delApps := []string{}
		for _, app := range *req {
			if len(app.Id) == 0 {
				return nil, errors.ParameterError("There is an empty id in the request.")
			}
			delApps = append(delApps, app.Id)
		}
		resp.Total, resp.Error = s.DeleteApps(ctx, delApps)
		return &resp, nil
	}
}

func (m UpdateAppRequest) GetUsers() (users []*models.User) {
	for _, u := range m.Users {
		users = append(users, &models.User{Model: models.Model{Id: u.Id}, RoleId: u.RoleId})
	}
	return users
}

func (m UpdateAppRequest) GetRoles() (roles []*models.AppRole) {
	for _, role := range m.Roles {
		appRole := &models.AppRole{
			Model:     models.Model{Id: role.Id},
			Name:      role.Name,
			IsDefault: role.IsDefault,
		}
		for _, urlId := range role.Urls {
			appRole.Urls = append(appRole.Urls, &models.AppProxyUrl{
				Model: models.Model{Id: urlId},
			})
		}
		roles = append(roles, appRole)
	}
	return roles
}

func MakeUpdateAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdateAppRequest)
		resp := BaseResponse{}
		if resp.Error = s.UpdateApp(ctx, &models.App{
			Model:       models.Model{Id: req.Id},
			Name:        req.Name,
			Description: req.Description,
			Avatar:      req.Avatar,
			GrantType:   models.NewGrantType(req.GrantType...),
			GrantMode:   req.GrantMode,
			Users:       req.GetUsers(),
			Roles:       req.GetRoles(),
			Proxy:       req.GetProxyConfig(),
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
		resp.Data, resp.Error = s.GetAppInfo(ctx, opts.WithAppId(req.Id))
		return &resp, nil
	}
}

func (r CreateAppRequest) GetUsers() (users []*models.User) {
	for _, u := range r.Users {
		users = append(users, &models.User{Model: models.Model{Id: u.Id}, RoleId: u.RoleId})
	}
	return users
}

func (r CreateAppRequest) GetRoles() (roles []*models.AppRole) {
	for _, role := range r.Roles {
		appRole := &models.AppRole{
			Model:     models.Model{Id: role.Id},
			Name:      role.Name,
			IsDefault: role.IsDefault,
		}
		for _, urlId := range role.Urls {
			appRole.Urls = append(appRole.Urls, &models.AppProxyUrl{
				Model: models.Model{Id: urlId},
			})
		}
		roles = append(roles, appRole)
	}
	return roles
}

func (r CreateAppRequest) GetProxyConfig() *models.AppProxy {
	if r.Proxy == nil {
		return nil
	}
	proxy := &models.AppProxy{
		Domain:                r.Proxy.Domain,
		Upstream:              r.Proxy.Upstream,
		InsecureSkipVerify:    r.Proxy.InsecureSkipVerify,
		TransparentServerName: r.Proxy.TransparentServerName,
	}
	for _, url := range r.Proxy.Urls {
		proxy.Urls = append(
			proxy.Urls,
			&models.AppProxyUrl{
				Model:  models.Model{Id: url.Id},
				Name:   url.Name,
				Method: url.Method,
				Url:    url.Url,
			})
	}
	return proxy
}

func MakeCreateAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateAppRequest)
		resp := BaseResponse{}
		resp.Error = s.CreateApp(ctx, &models.App{
			Name:        req.Name,
			Description: req.Description,
			Avatar:      req.Avatar,
			GrantType:   models.NewGrantType(req.GrantType...),
			GrantMode:   req.GrantMode,
			Users:       req.GetUsers(),
			Roles:       req.GetRoles(),
			Proxy:       req.GetProxyConfig(),
		})
		return &resp, nil
	}
}

func MakePatchAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchAppRequest)
		resp := BaseResponse{}

		if len(req.Id) == 0 {
			return nil, errors.ParameterError("There is an empty id in the patch.")
		}
		tmpPatch := map[string]interface{}{
			"id":          req.Id,
			"name":        req.Name,
			"description": req.Description,
			"avatar":      req.Avatar,
			"grantType":   req.GrantType,
			"grantMode":   req.GrantMode,
			"status":      req.Status,
			"isDelete":    req.IsDelete,
		}
		patch := map[string]interface{}{}
		for name, val := range tmpPatch {
			if val != nil {
				patch[name] = val
			}
		}
		resp.Error = s.PatchApp(ctx, patch)
		return &resp, nil
	}
}

func MakeDeleteAppEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteAppRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		resp.Error = s.DeleteApp(ctx, req.Id)
		return &resp, nil
	}
}

type AppProxyUrls []*AppProxyUrl

func (m *UpdateAppRequest) GetProxyConfig() *models.AppProxy {
	if m.Proxy == nil {
		return nil
	}
	proxy := &models.AppProxy{
		AppId:                 m.Id,
		Domain:                m.Proxy.Domain,
		Upstream:              m.Proxy.Upstream,
		InsecureSkipVerify:    m.Proxy.InsecureSkipVerify,
		TransparentServerName: m.Proxy.TransparentServerName,
	}
	for _, url := range m.Proxy.Urls {
		proxy.Urls = append(
			proxy.Urls,
			&models.AppProxyUrl{
				Model:  models.Model{Id: url.Id},
				Name:   url.Name,
				Method: url.Method,
				Url:    url.Url,
			})
	}
	return proxy
}

func MakeAppAuthenticationEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*AuthenticationRequest)
		return s.AppAuthentication(ctx, req.AuthKey, req.AuthSecret)
	}
}

func MakeCreateAppKeyEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreateAppKeyRequest)
		resp := SimpleResponseWrapper[*models.AppKey]{}
		resp.Data, resp.Error = s.CreateAppKey(ctx, req.AppId, req.Name)
		return &resp, nil
	}
}

func MakeDeleteAppKeyEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeleteAppKeysRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		_, resp.Error = s.DeleteAppKey(ctx, req.AppId, req.Id)
		return &resp, nil
	}
}

func MakeGetAppKeysEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetAppKeysRequest)
		resp := NewBaseListResponse[[]*models.AppKey](&req.BaseListRequest)
		resp.Total, resp.Data, resp.Error = s.GetAppKeys(ctx, req.AppId, req.Current, req.PageSize)
		return &resp, nil
	}
}
