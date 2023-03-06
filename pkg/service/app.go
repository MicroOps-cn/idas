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
	"github.com/MicroOps-cn/fuck/sets"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
	"github.com/MicroOps-cn/idas/pkg/utils/image"
	"github.com/go-kit/log/level"
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

func (s Set) GetApps(ctx context.Context, storage string, keywords string, filter map[string]interface{}, current, pageSize int64) (total int64, apps []*models.App, err error) {
	return s.SafeGetUserAndAppService(storage).GetApps(ctx, keywords, filter, current, pageSize)
}

func (s Set) GetAppSource(_ context.Context) (total int64, data map[string]string, err error) {
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

func (s Set) GetAppAccessControl(ctx context.Context, app *models.App, o ...opts.WithGetAppOptions) (err error) {
	var appUsersRel models.AppUsers
	appUsersRel, app.Roles, err = s.commonService.GetAppAccessControl(ctx, app.Id, o...)
	opt := opts.NewAppOptions(o...)
	if err != nil {
		return err
	}
	userIds := sets.New[string](appUsersRel.UserId()...).Difference(sets.New[string](app.Users.Id()...))
	for _, rel := range appUsersRel {
		user := app.Users.GetById(rel.UserId)
		if user != nil {
			user.RoleId = rel.RoleId
		} else {
			user = &models.User{Model: models.Model{Id: rel.UserId}, RoleId: rel.RoleId}
			app.Users = append(app.Users, user)
		}
		if role := app.Roles.GetRoleById(rel.RoleId); role != nil {
			user.Role = role.Name
		}
	}
	if userIds.Len() > 0 && !opt.DisableGetUsers {
		var tmpUsers models.Users
		for _, svc := range s.userAndAppService {
			tmpUsers, err = svc.GetUsersById(ctx, userIds.List())
			if err != nil {
				return err
			}
			for _, user := range tmpUsers {
				userIds.Delete(user.Id)
				if appUserRel := appUsersRel.GetByUserId(user.Id); appUserRel != nil {
					user.RoleId = appUserRel.RoleId
					if role := app.Roles.GetRoleById(appUserRel.RoleId); role != nil {
						user.Role = role.Name
					}
				}
				if u := app.Users.GetById(user.Id); u != nil {
					*u = *user
				} else {
					app.Users = append(app.Users, user)
				}
			}
		}
	}
	return nil
}

func (s Set) UpdateApp(ctx context.Context, storage string, app *models.App, updateColumns ...string) (err error) {
	if len(app.Name) == 0 {
		return errors.ParameterError("app name is null")
	}

	if len(app.Id) == 0 {
		return errors.ParameterError("app id is null")
	}

	if app.GrantType&models.AppMeta_proxy > 0 {
		if app.Proxy == nil {
			return errors.ParameterError("app proxy cannot be empty")
		}
		if len(app.Proxy.Urls) == 0 {
			return errors.ParameterError("app proxy.urls cannot be empty")
		}
		if len(app.Proxy.Domain) == 0 {
			return errors.ParameterError("app proxy.domain cannot be empty")
		}
		if len(app.Proxy.Upstream) == 0 {
			return errors.ParameterError("app proxy.upstream cannot be empty")
		}
	}

	proxy := app.Proxy
	app.Proxy = nil

	roles := sets.New[string]()

	for _, role := range app.Roles {
		if len(role.Name) == 0 && len(role.Id) == 0 {
			return errors.ParameterError("role name and id is nil")
		}
		if len(role.Name) != 0 {
			if roles.Has(role.Name) {
				return errors.ParameterError("duplicate role: " + role.Name)
			}
			roles.Insert(role.Name)
		}
		if len(role.Id) != 0 {
			if roles.Has(role.Id) {
				return errors.ParameterError("duplicate role: " + role.Id)
			}
			roles.Insert(role.Id)
		}
	}

	if err = s.GetUserAndAppService(storage).UpdateApp(ctx, app, updateColumns...); err != nil {
		return err
	}

	if app.GrantType&models.AppMeta_proxy > 0 {
		proxy.AppId = app.Id
		err = s.commonService.UpdateAppProxyConfig(ctx, proxy)
		if err != nil {
			return err
		}
	}

	err = s.commonService.UpdateAppAccessControl(ctx, app)
	if err != nil {
		return err
	}

	return nil
}

func (s Set) GetAppInfo(ctx context.Context, storage string, o ...opts.WithGetAppOptions) (app *models.App, err error) {
	if len(storage) == 0 {
		for _, service := range s.userAndAppService {
			app, err = service.GetAppInfo(ctx, o...)
			if err != nil {
				continue
			}
			break
		}
	} else {
		service := s.GetUserAndAppService(storage)
		if service == nil {
			err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
			return
		}
		app, err = service.GetAppInfo(ctx, o...)
		if err != nil {
			return nil, err
		}
	}
	if app == nil {
		return nil, errors.StatusNotFound("app")
	}
	if !opts.NewAppOptions(o...).DisableGetProxy {
		if app.GrantType&models.AppMeta_proxy == models.AppMeta_proxy {
			if app.Proxy, err = s.commonService.GetAppProxyConfig(ctx, app.Id); err != nil {
				return nil, err
			}
		}
	}
	if !opts.NewAppOptions(o...).DisableGetAccessController {
		if err = s.GetAppAccessControl(ctx, app, o...); err != nil {
			return nil, err
		}
	}
	return app, nil
}

func (s Set) CreateApp(ctx context.Context, storage string, app *models.App) (err error) {
	if len(app.Name) == 0 {
		return errors.ParameterError("app name cannot be empty")
	}

	if app.GrantType&models.AppMeta_proxy > 0 {
		if app.Proxy == nil {
			return errors.ParameterError("app proxy cannot be empty")
		}
		if len(app.Proxy.Urls) == 0 {
			return errors.ParameterError("app proxy.urls cannot be empty")
		}
		if len(app.Proxy.Domain) == 0 {
			return errors.ParameterError("app proxy.domain cannot be empty")
		}
		if len(app.Proxy.Upstream) == 0 {
			return errors.ParameterError("app proxy.upstream cannot be empty")
		}
	}

	proxy := app.Proxy
	app.Proxy = nil

	roles := sets.New[string]()

	for _, role := range app.Roles {
		if len(role.Name) == 0 && len(role.Id) == 0 {
			return errors.ParameterError("role name and id is nil")
		}
		if len(role.Name) != 0 {
			if roles.Has(role.Name) {
				return errors.ParameterError("duplicate role: " + role.Name)
			}
			roles.Insert(role.Name)
		}
		if len(role.Id) != 0 {
			if roles.Has(role.Id) {
				return errors.ParameterError("duplicate role: " + role.Id)
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

	if err = service.CreateApp(ctx, app); err != nil {
		return err
	}
	if app.GrantType&models.AppMeta_proxy > 0 {
		proxy.AppId = app.Id
		err = s.commonService.UpdateAppProxyConfig(ctx, proxy)
		if err != nil {
			return err
		}
	}

	return s.commonService.UpdateAppAccessControl(ctx, app)
}

func (s Set) PatchApp(ctx context.Context, storage string, fields map[string]interface{}) (err error) {
	return s.GetUserAndAppService(storage).PatchApp(ctx, fields)
}

func (s Set) DeleteApp(ctx context.Context, storage string, id string) (err error) {
	return s.GetUserAndAppService(storage).DeleteApp(ctx, id)
}

func (s Set) AppAuthentication(ctx context.Context, key string, secret string) (*models.App, error) {
	logger := logs.GetContextLogger(ctx)
	if appId, err := s.commonService.AppAuthorization(ctx, key, secret); err != nil {
		level.Error(logger).Log("msg", "failed to authorization app", "err", err)
	} else if len(appId) > 0 {
		if app, err := s.GetAppInfo(ctx, "", opts.WithAppId(appId), opts.WithBasic); err != nil {
			level.Error(logger).Log("msg", "failed to get app info", "err", err)
		} else {
			return app, nil
		}
	}
	return nil, nil
}

func (s Set) CreateAppKey(ctx context.Context, appId string, name string) (appKey *models.AppKey, err error) {
	return s.commonService.CreateAppKey(ctx, appId, name)
}

func (s Set) GetAppKeys(ctx context.Context, appId string, current int64, pageSize int64) (count int64, keys []*models.AppKey, err error) {
	return s.commonService.GetAppKeys(ctx, appId, current, pageSize)
}

func (s Set) DeleteAppKey(ctx context.Context, appId string, id []string) (affected int64, err error) {
	return s.commonService.DeleteAppKeys(ctx, appId, id)
}

func (s Set) GetAppKeyFromKey(ctx context.Context, key string) (appKey *models.AppKey, err error) {
	return s.commonService.GetAppKeyFromKey(ctx, key)
}

func (s Set) GetAppRoleByUserId(ctx context.Context, appId string, userId string) (role *models.AppRole, err error) {
	return s.commonService.GetAppRoleByUserId(ctx, appId, userId)
}
