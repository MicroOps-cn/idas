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

package service

import (
	"context"
	"fmt"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/log/level"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/image"
	"github.com/MicroOps-cn/idas/pkg/utils/sets"
)

func (s Set) SafeGetUserAndAppService(name string) UserAndAppService {
	for _, svc := range s.userAndAppService {
		if svc.Name() == name {
			return svc
		}
	}
	for _, svc := range s.userAndAppService {
		return svc
	}
	return newNullService("", name)
}

func (s Set) GetApps(ctx context.Context, storage string, keywords string, current, pageSize int64) (total int64, apps []*models.App, err error) {
	return s.SafeGetUserAndAppService(storage).GetApps(ctx, keywords, current, pageSize)
}

func (s Set) GetAppSource(ctx context.Context) (total int64, data map[string]string, err error) {
	total = int64(len(s.userAndAppService))
	data = make(map[string]string, total)
	for _, appService := range s.userAndAppService {
		data[appService.Name()] = appService.Name()
	}

	return
}

func (s Set) PatchApps(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, err error) {
	return s.GetUserAndAppService(storage).PatchApps(ctx, patch)
}

func (s Set) DeleteApps(ctx context.Context, storage string, id []string) (total int64, err error) {
	return s.GetUserAndAppService(storage).DeleteApps(ctx, id)
}

func (s Set) UpdateApp(ctx context.Context, storage string, app *models.App, updateColumns ...string) (a *models.App, err error) {
	if len(app.Name) == 0 {
		return nil, errors.ParameterError("app name is null")
	}

	if len(app.Id) == 0 {
		return nil, errors.ParameterError("app id is null")
	}

	if app.GrantType&models.AppMeta_proxy > 0 {
		if app.Proxy == nil {
			return nil, errors.ParameterError("app proxy cannot be empty")
		}
		if len(app.Proxy.Urls) == 0 {
			return nil, errors.ParameterError("app proxy.urls cannot be empty")
		}
		if len(app.Proxy.Domain) == 0 {
			return nil, errors.ParameterError("app proxy.domain cannot be empty")
		}
		if len(app.Proxy.Upstream) == 0 {
			return nil, errors.ParameterError("app proxy.upstream cannot be empty")
		}
	}

	proxy := app.Proxy
	app.Proxy = nil

	roles := sets.New[string]()

	for _, role := range app.Roles {
		if len(role.Name) == 0 && len(role.Id) == 0 {
			return nil, errors.ParameterError("role name and id is nil")
		}
		if len(role.Name) != 0 {
			if roles.Has(role.Name) {
				return nil, errors.ParameterError("duplicate role: " + role.Name)
			}
			roles.Insert(role.Name)
		}
		if len(role.Id) != 0 {
			if roles.Has(role.Id) {
				return nil, errors.ParameterError("duplicate role: " + role.Id)
			}
			roles.Insert(role.Id)
		}
	}

	newApp, err := s.GetUserAndAppService(storage).UpdateApp(ctx, app, updateColumns...)
	if err != nil {
		return nil, err
	}

	if app.GrantType&models.AppMeta_proxy > 0 {
		proxy.Storage = app.Storage
		proxy.AppId = app.Id
		newApp.Proxy, err = s.commonService.UpdateProxyConfig(ctx, proxy)
		if err != nil {
			return nil, err
		}
	}

	return newApp, nil
}

func (s Set) GetAppInfo(ctx context.Context, storage string, id string) (app *models.App, err error) {
	if len(storage) == 0 {
		for _, service := range s.userAndAppService {
			info, err := service.GetAppInfo(ctx, id, "")
			if err != nil {
				continue
			}
			return info, nil
		}
	} else {
		service := s.GetUserAndAppService(storage)
		if service == nil {
			err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
			return
		}
		info, err := service.GetAppInfo(ctx, id, "")
		if err != nil {
			return nil, err
		}
		if info.GrantType&models.AppMeta_proxy == models.AppMeta_proxy {
			info.Proxy, err = s.commonService.GetAppProxyConfig(ctx, info.Id)
			if err != nil {
				return nil, err
			}
		}
		return info, nil
	}
	return nil, errors.StatusNotFound("app")
}

func (s Set) CreateApp(ctx context.Context, storage string, app *models.App) (a *models.App, err error) {
	if len(app.Name) == 0 {
		return nil, errors.ParameterError("app name cannot be empty")
	}

	if app.GrantType&models.AppMeta_proxy > 0 {
		if app.Proxy == nil {
			return nil, errors.ParameterError("app proxy cannot be empty")
		}
		if len(app.Proxy.Urls) == 0 {
			return nil, errors.ParameterError("app proxy.urls cannot be empty")
		}
		if len(app.Proxy.Domain) == 0 {
			return nil, errors.ParameterError("app proxy.domain cannot be empty")
		}
		if len(app.Proxy.Upstream) == 0 {
			return nil, errors.ParameterError("app proxy.upstream cannot be empty")
		}
	}

	proxy := app.Proxy
	app.Proxy = nil

	roles := sets.New[string]()

	for _, role := range app.Roles {
		if len(role.Name) == 0 && len(role.Id) == 0 {
			return nil, errors.ParameterError("role name and id is nil")
		}
		if len(role.Name) != 0 {
			if roles.Has(role.Name) {
				return nil, errors.ParameterError("duplicate role: " + role.Name)
			}
			roles.Insert(role.Name)
		}
		if len(role.Id) != 0 {
			if roles.Has(role.Id) {
				return nil, errors.ParameterError("duplicate role: " + role.Id)
			}
			roles.Insert(role.Id)
		}
	}
	logger := logs.GetContextLogger(ctx)
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	if avatar, err := image.GenerateAvatar(ctx, app.Name); err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to generate avatar")
	} else if fileKey, err := s.UploadFile(ctx, app.Name+".png", "image/png", avatar); err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to save avatar")
	} else {
		app.Avatar = fileKey
	}
	newApp, err := service.CreateApp(ctx, app)
	if err != nil {
		return nil, err
	}
	if app.GrantType&models.AppMeta_proxy > 0 {
		proxy.Storage = app.Storage
		proxy.AppId = app.Id
		newApp.Proxy, err = s.commonService.UpdateProxyConfig(ctx, proxy)
		if err != nil {
			return nil, err
		}
	}

	return newApp, nil
}

func (s Set) PatchApp(ctx context.Context, storage string, fields map[string]interface{}) (app *models.App, err error) {
	return s.GetUserAndAppService(storage).PatchApp(ctx, fields)
}

func (s Set) DeleteApp(ctx context.Context, storage string, id string) (err error) {
	return s.GetUserAndAppService(storage).DeleteApp(ctx, id)
}
