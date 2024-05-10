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

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/fuck/sets"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/go-kit/log/level"
	"github.com/golang/groupcache/lru"
	uuid "github.com/satori/go.uuid"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
	"github.com/MicroOps-cn/idas/pkg/utils/image"
	jwtutils "github.com/MicroOps-cn/idas/pkg/utils/jwt"
)

func (s Set) GetApps(ctx context.Context, keywords string, filter map[string]interface{}, current, pageSize int64) (total int64, apps []*models.App, err error) {
	count, apps, err := s.GetUserAndAppService().GetApps(ctx, keywords, filter, current, pageSize)
	for _, app := range apps {
		app.Roles, err = s.commonService.GetAppRoles(ctx, app.Id)
		if err != nil {
			return 0, nil, errors.WithServerError(500, err, "failed to get app roles")
		}
		if app.Proxy != nil {
			app.Proxy.JwtSecretSalt = nil
			app.Proxy.JwtSecret = nil
		}
		app.I18N = new(models.AppI18NOptions)
		var i18n map[string]string
		if i18n, err = s.commonService.GetI18n(ctx, "app", app.Id, "description"); err != nil {
			return 0, nil, err
		}
		app.I18N.Description = i18n
		if i18n, err = s.commonService.GetI18n(ctx, "app", app.Id, "displayName"); err != nil {
			return 0, nil, err
		}
		app.I18N.DisplayName = i18n
	}
	if len(filter) == 0 {
		appIds, count2 := s.commonService.FindAppByKeywords(ctx, keywords, (current-1)*pageSize-count, pageSize-int64(len(apps)))
		if count2 > 0 {
			count += count2
		loop:
			for _, id := range appIds {
				app, err := s.GetAppInfo(ctx, opts.WithAppId(id), opts.WithBasic)
				if err != nil {
					continue
				}
				for _, a := range apps {
					if a.Id == app.Id {
						continue loop
					}
				}
				apps = append(apps, app)
			}

		}
	}
	return count, apps, err
}

func (s Set) PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error) {
	return s.GetUserAndAppService().PatchApps(ctx, patch)
}

func (s Set) DeleteApps(ctx context.Context, id ...string) (total int64, err error) {
	if total, err = s.GetUserAndAppService().DeleteApps(ctx, id...); err != nil {
		return total, err
	}
	if err = s.commonService.DeleteAppProxy(ctx, id...); err != nil {
		return total, err
	}
	if err = s.commonService.DeleteAppAccessControl(ctx, id...); err != nil {
		return total, err
	}
	if err = s.commonService.DeleteI18nBySourceId(ctx, id...); err != nil {
		return total, err
	}

	return total, nil
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
			continue
		}
		if role := app.Roles.GetRoleById(rel.RoleId); role != nil {
			user.Role = role.Name
		}
	}
	if userIds.Len() > 0 && !opt.DisableGetUsers {
		var tmpUsers models.Users
		tmpUsers, err = s.userAndAppService.GetUsersById(ctx, userIds.List())
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
	return nil
}

func (s Set) PatchAppI18n(ctx context.Context, appId string, options *models.AppI18NOptions) (err error) {
	if options == nil {
		return nil
	}
	var i18ns []models.I18nTranslate

	for lang, val := range options.DisplayName {
		i18ns = append(i18ns, models.I18nTranslate{
			Source:   "app",
			Field:    "displayName",
			SourceId: appId,
			Lang:     lang,
			Value:    val,
		})
	}
	for lang, val := range options.Description {
		i18ns = append(i18ns, models.I18nTranslate{
			Source:   "app",
			Field:    "description",
			SourceId: appId,
			Lang:     lang,
			Value:    val,
		})
	}

	return s.BatchPatchI18n(ctx, i18ns)
}

func (s Set) UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (err error) {
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

	if err = s.GetUserAndAppService().UpdateApp(ctx, app, updateColumns...); err != nil {
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
	if len(app.Avatar) > 0 {
		if keyUUID, err := uuid.FromString(app.Avatar); err == nil {
			if err := s.commonService.UpdateFileOwner(ctx, keyUUID.String(), "app:icon:public"); err != nil {
				logger := logs.GetContextLogger(ctx)
				level.Error(logger).Log("err", err, "msg", "failed to update file owner")
			}
		}
	}
	if app.GrantType&models.AppMeta_authorization_code > 0 {
		app.OAuth2.AppId = app.Id
		if err = s.commonService.PatchAppOAuthConfig(ctx, app.OAuth2); err != nil {
			logger := logs.GetContextLogger(ctx)
			level.Error(logger).Log("err", err, "msg", "failed to update oauth config")
		}
		issuerCache.Remove(app.Id)
	}
	return s.PatchAppI18n(ctx, app.Id, app.I18N)
}

var issuerCache = lru.New(16)

func (s Set) GetIssuerByAppId(ctx context.Context, appId string) (jwtutils.JWTIssuer, error) {
	if value, ok := issuerCache.Get(appId); ok {
		if issuer, ok := value.(jwtutils.JWTIssuer); ok {
			return issuer, nil
		}
	}
	oAuthConfig, err := s.commonService.GetAppOAuthConfig(ctx, appId)
	if err != nil {
		return nil, err
	}
	return oAuthConfig.GetJWTIssuer(ctx), nil
}

func (s Set) GetAppInfo(ctx context.Context, o ...opts.WithGetAppOptions) (app *models.App, err error) {
	service := s.GetUserAndAppService()
	if app, err = service.GetAppInfo(ctx, o...); err != nil {
		return nil, err
	} else if app == nil {
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
	if !opts.NewAppOptions(o...).DisableGetI18n {
		app.I18N = new(models.AppI18NOptions)
		var i18n map[string]string
		if i18n, err = s.commonService.GetI18n(ctx, "app", app.Id, "description"); err != nil {
			return nil, err
		}
		app.I18N.Description = i18n
		if i18n, err = s.commonService.GetI18n(ctx, "app", app.Id, "displayName"); err != nil {
			return nil, err
		}
		app.I18N.DisplayName = i18n
	}
	if !opts.NewAppOptions(o...).DisableGetOAuth2 {
		if app.OAuth2, err = s.commonService.GetAppOAuthConfig(ctx, app.Id); err != nil {
			return nil, err
		}
	}
	return app, nil
}

func (s Set) CreateApp(ctx context.Context, app *models.App) (err error) {
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
	service := s.GetUserAndAppService()
	if len(app.Avatar) == 0 {
		if avatar, err := image.GenerateAvatar(ctx, app.Name); err != nil {
			level.Error(logger).Log("err", err, "msg", "failed to generate avatar")
		} else if fileKey, err := s.UploadFile(ctx, app.Name+".png", "image/png", avatar); err != nil {
			level.Error(logger).Log("err", err, "msg", "failed to save avatar")
		} else {
			app.Avatar = fileKey
			if err := s.commonService.UpdateFileOwner(ctx, fileKey, "app:icon:public"); err != nil {
				level.Error(logger).Log("err", err, "msg", "failed to update file owner")
			}
		}
	} else {
		if keyUUID, err := uuid.FromString(app.Avatar); err == nil {
			if err := s.commonService.UpdateFileOwner(ctx, keyUUID.String(), "app:icon:public"); err != nil {
				level.Error(logger).Log("err", err, "msg", "failed to update file owner")
			}
		}
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
	if err = s.commonService.UpdateAppAccessControl(ctx, app); err != nil {
		return err
	}

	if app.GrantType&models.AppMeta_authorization_code > 0 {
		app.OAuth2.AppId = app.Id
		if err = s.commonService.PatchAppOAuthConfig(ctx, app.OAuth2); err != nil {
			logger := logs.GetContextLogger(ctx)
			level.Error(logger).Log("err", err, "msg", "failed to update oauth config")
		}
		issuerCache.Remove(app.Id)
	}
	return s.PatchAppI18n(ctx, app.Id, app.I18N)
}

func (s Set) PatchApp(ctx context.Context, fields map[string]interface{}) (err error) {
	return s.GetUserAndAppService().PatchApp(ctx, fields)
}

func (s Set) DeleteApp(ctx context.Context, id string) (err error) {
	return w.E(s.DeleteApps(ctx, id))
}

func (s Set) AppAuthentication(ctx context.Context, key string, secret string) (*models.App, error) {
	logger := logs.GetContextLogger(ctx)
	if appId, err := s.commonService.AppAuthorization(ctx, key, secret); err != nil {
		level.Error(logger).Log("msg", "failed to authorization app", "err", err)
	} else if len(appId) > 0 {
		if app, err := s.GetAppInfo(ctx, opts.WithAppId(appId), opts.WithBasic, opts.WithOAuth2); err != nil {
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

func (s Set) GetAppIcons(ctx context.Context, current int64, pageSize int64) (count int64, keys []*models.Model, err error) {
	return s.commonService.GetFilesByOwner(ctx, "app:icon:public", current, pageSize)
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
