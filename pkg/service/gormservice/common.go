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

package gormservice

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/MicroOps-cn/fuck/clients/gorm"
	"github.com/MicroOps-cn/fuck/safe"
	gogorm "gorm.io/gorm"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
	"github.com/MicroOps-cn/idas/pkg/utils/jwt"
	"github.com/MicroOps-cn/idas/pkg/utils/sign"
)

func NewCommonService(name string, client *gorm.Client) *CommonService {
	return &CommonService{name: name, Client: client}
}

type CommonService struct {
	*gorm.Client
	name string
}

var findAppSql = `
SELECT DISTINCT
        (app_id)
    FROM
        t_app_key AS t1
    WHERE
        t1.key LIKE ? UNION SELECT DISTINCT
        (source_id) AS app_id
    FROM
        t_i18n_translate AS t2
    WHERE
        t2.source = 'app'
            AND t2.value LIKE ?
`

func (c CommonService) FindAppByKeywords(ctx context.Context, keywords string, skip int64, limit int64) (ids []string, count int64) {
	err := c.Session(ctx).Raw(fmt.Sprintf("SELECT count(1) from (%s) t", findAppSql), fmt.Sprintf("%%%s%%", keywords), fmt.Sprintf("%%%s%%", keywords)).Scan(&count).Error
	if err != nil {
		return nil, 0
	}
	if skip < 0 {
		skip = 0
	}
	if limit < 0 {
		return nil, count
	}
	err = c.Session(ctx).Raw(fmt.Sprintf("SELECT app_id from (%s) t limit %d,%d", findAppSql, skip, limit), fmt.Sprintf("%%%s%%", keywords), fmt.Sprintf("%%%s%%", keywords)).Scan(&ids).Error
	if err != nil {
		return nil, count
	}
	return ids, count
}

func (c CommonService) Name() string {
	return c.name
}

func (c CommonService) AutoMigrate(ctx context.Context) error {
	err := c.Session(ctx).AutoMigrate(
		&models.File{},
		&models.Permission{},
		&models.Role{},
		&models.AppKey{},
		&models.AppProxy{},
		&models.AppProxyUrl{},
		&models.AppUser{},
		&models.AppRole{},
		&models.PageConfig{},
		&models.PageData{},
		&models.UserExt{},
		&models.SystemConfig{},
		&models.UserPasswordHistory{},
		&models.WeakPassword{},
		&models.I18nTranslate{},
		&models.AppOAuth2{},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c CommonService) RecordUploadFile(ctx context.Context, name string, path string, contentType string, size int64) (id string, err error) {
	file := &models.File{MimiType: contentType, Name: name, Path: path, Size: size}
	if err = c.Session(ctx).Create(file).Error; err != nil {
		return
	}
	return file.Id, err
}

func (c CommonService) UpdateFileOwner(ctx context.Context, fileId string, owner string) (err error) {
	return c.Session(ctx).Model(&models.File{}).Where("id = ?", fileId).Update("owner", owner).Error
}

func (c CommonService) GetFilesByOwner(ctx context.Context, owner string, current int64, pageSize int64) (count int64, files []*models.Model, err error) {
	tb := c.Session(ctx).Model(&models.File{})
	if len(owner) > 0 {
		tb.Where("`owner` = ?", owner)
	}
	if err = tb.Count(&count).Error; err != nil {
		return 0, nil, err
	} else if count == 0 {
		return 0, nil, nil
	}

	if err = tb.Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&files).Error; err != nil {
		return 0, nil, err
	}
	return
}

func (c CommonService) GetFileInfoFromId(ctx context.Context, id string) (fileName, mimiType, filePath string, err error) {
	file := &models.File{Model: models.Model{Id: id}}
	if err = c.Session(ctx).First(file).Error; err != nil {
		return "", "", "", err
	}
	return file.Name, file.MimiType, file.Path, nil
}

func (c CommonService) GetProxyConfig(ctx context.Context, host string) (*models.AppProxyConfig, error) {
	conn := c.Session(ctx)
	var proxy models.AppProxyConfig

	if err := conn.Preload("Urls").
		Order("id desc").
		First(&proxy.AppProxy, "domain = ?", host).Error; err != nil {
		if err == gogorm.ErrRecordNotFound {
			return nil, errors.StatusNotFound("page")
		}
		return nil, errors.NewServerError(500, err.Error())
	}
	if err := conn.Table("t_app_role_url").
		Select("app_role_id", "`name` as `app_role_name`", "app_proxy_url_id").
		Joins("JOIN `t_app_role` ON `t_app_role`.id = `t_app_role_url`.app_role_id").
		Where("app_proxy_url_id in ?", proxy.Urls.Id()).
		Scan(&proxy.URLRoles).Error; err != nil {
		return nil, err
	}
	sort.Sort(proxy.Urls)
	return &proxy, nil
}

func (c CommonService) UpdateAppProxyConfig(ctx context.Context, proxy *models.AppProxy) (err error) {
	conn := c.Session(ctx)
	if len(proxy.JwtSecretSalt) == 0 && len(proxy.JwtSecret) > 0 {
		if err = proxy.SetJwtSecret(string(proxy.JwtSecret)); err != nil {
			return err
		}
	}
	var model models.Model
	for i, url := range proxy.Urls {
		url.Index = uint32(i)
	}
	if err = conn.Model(&models.AppProxy{}).Select("id").Where("app_id = ?", proxy.AppId).First(&model).Error; err != nil && err != gogorm.ErrRecordNotFound {
		return err
	}
	proxy.Id = model.Id
	if len(proxy.Id) > 0 {
		if err = conn.Select(
			"Urls", "update_time", "domain", "upstream",
			"insecure_skip_verify", "transparent_server_name",
			"jwt_provider", "jwt_cookie_name", "jwt_secret", "jwt_secret_salt", "hsts_offload",
		).Session(&gogorm.Session{FullSaveAssociations: true}).
			Updates(proxy).Error; err != nil {
			return err
		}
	} else {
		if err = conn.Create(proxy).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c CommonService) GetAppProxyConfig(ctx context.Context, appId string) (proxy *models.AppProxy, err error) {
	if err = c.Session(ctx).Where("app_id = ?", appId).Preload("Urls").First(&proxy).Error; err != nil && !errors.IsNotFount(err) {
		return nil, err
	}
	proxy.JwtSecret = nil
	proxy.JwtSecretSalt = nil
	sort.Sort(proxy.Urls)
	return proxy, nil
}

// PatchAppRole
//
//	@Description[en-US]: Update App Role.
//	@Description[zh-CN]: 更新应用角色。
//	@param ctx     context.Context
//	@param dn      string
//	@param patch   *models.AppRole
//	@return err    error
func (c CommonService) PatchAppRole(ctx context.Context, role *models.AppRole) error {
	conn := c.Session(ctx)
	var r models.AppRole
	if len(role.Id) != 0 {
		if err := conn.Where("id = ? and app_id = ?", role.Id, role.AppId).First(&r).Error; err == gogorm.ErrRecordNotFound {
			if err = conn.Create(&role).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			if role.Name != r.Name || role.IsDefault != r.IsDefault || role.IsDelete != r.IsDelete {
				if err = conn.Select("name", "config", "is_delete", "is_default").Updates(&role).Error; err != nil {
					return err
				}
			}
			if err = conn.Model(&role).Association("Urls").Replace(role.Urls); err != nil {
				return err
			}
		}
	} else if role.Name != "" {
		if err := conn.Create(&role).Error; err != nil {
			return err
		}
	}
	var userIds []string
	var oriUsers models.AppUsers
	if err := conn.Unscoped().Select("user_id", "role_id", "delete_time").Where("app_id = ?", role.AppId).Find(&oriUsers).Error; err != nil {
		return err
	}
	for _, user := range role.Users {
		if oriUser := oriUsers.GetByUserId(user.Id); oriUser == nil {
			appUser := models.AppUser{AppId: role.AppId, UserId: user.Id, RoleId: role.Id}
			if err := conn.Create(&appUser).Error; err != nil {
				return err
			}
		} else if oriUser.RoleId != role.Id || oriUser.DeleteTime.Valid {
			if err := conn.Unscoped().Model(&models.AppUser{}).Where("app_id = ? and user_id = ?", role.AppId, user.Id).Updates(map[string]interface{}{
				"role_id":     role.Id,
				"delete_time": gogorm.Expr("null"),
			}).Error; err != nil {
				return err
			}
		}
		user.RoleId = role.Id
		userIds = append(userIds, user.Id)
	}
	return conn.Delete(&models.AppUser{}, "app_id = ? and role_id = ? and user_id not in ? ", role.AppId, role.Id, userIds).Error
}

func (c CommonService) GetAppRoleByUserId(ctx context.Context, appId string, userId string) (role *models.AppRole, err error) {
	conn := c.Session(ctx)
	query := conn.Where("`t_app_role`.`app_id` = ?", appId)

	query = query.Select("`t_app_role`.*").
		Joins("JOIN `t_app_user` ON `t_app_user`.`role_id` = `t_app_role`.`id`").
		Where("`t_app_user`.`user_id` = ?", userId)

	if err = query.Find(&role).Error; err != nil {
		return nil, err
	}
	return
}

func (c CommonService) GetAppRoles(ctx context.Context, appId string) (roles models.AppRoles, err error) {
	conn := c.Session(ctx)
	err = conn.Where("app_id = ?", appId).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (c CommonService) GetAppAccessControl(ctx context.Context, appId string, o ...opts.WithGetAppOptions) (users models.AppUsers, roles models.AppRoles, err error) {
	conn := c.Session(ctx)
	opt := opts.NewAppOptions(o...)
	{
		query := conn.Where("app_id = ?", appId)
		if len(opt.UserId) != 0 {
			query = query.Where("user_id in ?", opt.UserId)
		}
		if err = query.Find(&users).Error; err != nil {
			return nil, nil, err
		}
	}
	{
		query := conn.Where("`t_app_role`.`app_id` = ?", appId)
		if len(opt.UserId) != 0 {
			query = query.Select("`t_app_role`.*").
				Joins("JOIN `t_app_user` ON `t_app_user`.`role_id` = `t_app_role`.`id`").
				Where("`t_app_user`.`user_id` in ?", opt.UserId)
		}
		if err = query.Find(&roles).Error; err != nil {
			return nil, nil, err
		}
	}
	if !opt.DisableGetProxy {
		var results []models.AppRoleURL
		if err = conn.Table("t_app_role_url").
			Select("app_role_id", "`t_app_role`.`name` as `app_role_name`", "app_proxy_url_id").
			Joins("JOIN `t_app_role` ON `t_app_role`.id = `t_app_role_url`.app_role_id").
			Joins("JOIN `t_app_proxy_url` ON `t_app_proxy_url`.id = `t_app_role_url`.app_proxy_url_id").
			Joins("JOIN `t_app_proxy` ON `t_app_proxy`.id = `t_app_proxy_url`.app_proxy_id").
			Where("app_role_id IN ?", roles.GetId()).
			Where("t_app_proxy.app_id = ?", appId).
			Where("t_app_proxy.`delete_time` IS NULL").
			Scan(&results).Error; err != nil {
			return nil, nil, err
		}
		for _, role := range roles {
			for _, result := range results {
				if result.AppRoleId == role.Id {
					role.UrlsId = append(role.UrlsId, result.AppProxyURLId)
				}
			}
		}
	}
	return
}

func (c CommonService) UpdateUserAccessControl(ctx context.Context, userId string, apps models.Apps) (err error) {
	tx := c.Session(ctx).Begin()
	defer tx.Rollback()
	var appUsers models.AppUsers
	if err = tx.Where("user_id = ?", userId).Find(&appUsers).Error; err != nil {
		return err
	}
	for _, app := range apps {
		if appUser := appUsers.GetByAppId(app.Id); appUser == nil {
			if err = tx.Create(&models.AppUser{AppId: app.Id, UserId: userId, RoleId: app.RoleId}).Error; err != nil {
				return err
			}
		} else {
			if appUser.RoleId != app.RoleId {
				if err = tx.Model(&models.AppUser{}).Where("id = ?", appUser.Id).Update("role_id", app.RoleId).Error; err != nil {
					return err
				}
			}
		}
	}
	if err = tx.Where("user_id = ? AND app_id NOT IN ?", userId, apps.Id()).Delete(&models.AppUser{}).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

func (c CommonService) UpdateAppAccessControl(ctx context.Context, app *models.App) (err error) {
	tx := c.Session(ctx).Begin()
	defer tx.Rollback()
	if len(app.Roles) > 0 {
		var roleIds []string
		for _, role := range app.Roles {
			for _, user := range app.Users {
				if len(user.RoleId) == 0 && len(user.Role) == 0 && role.IsDefault {
					role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
				} else if len(user.RoleId) != 0 && user.RoleId == role.Id {
					role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
				} else if len(user.Role) != 0 && user.Role == role.Name {
					role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
				}
			}
			role.AppId = app.Id
			if err = c.PatchAppRole(gorm.WithConnContext(ctx, tx), role); err != nil {
				return err
			}
			roleIds = append(roleIds, role.Id)
		}
		if err = tx.Delete(&models.AppRole{}, "app_id = ? and id not in ? ", app.Id, roleIds).Error; err != nil {
			return err
		}
	} else if len(app.Users) > 0 {
		role := models.AppRole{AppId: app.Id}
		for _, user := range app.Users {
			role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
		}
		if err = c.PatchAppRole(gorm.WithConnContext(ctx, tx), &role); err != nil {
			return err
		}
		if err = tx.Delete(&models.AppRole{}, "app_id = ?", app.Id).Error; err != nil {
			return err
		}
	} else {
		if err = tx.Delete(&models.AppUser{}, "app_id = ?", app.Id).Error; err != nil {
			return err
		}
		if err = tx.Delete(&models.AppRole{}, "app_id = ?", app.Id).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (c CommonService) CreateAppKey(ctx context.Context, appId, name string) (*models.AppKey, error) {
	conn := c.Session(ctx)
	pub1, _, privateKey, err := sign.GenerateECDSAKeyPair()
	if err != nil {
		return nil, err
	}
	appKey := &models.AppKey{
		Name:       name,
		AppId:      appId,
		Key:        pub1,
		Secret:     sign.SumSha256Hmac(pub1, privateKey),
		PrivateKey: privateKey,
	}
	if err = conn.Create(&appKey).Error; err != nil {
		return nil, err
	}
	return appKey, nil
}

func (c CommonService) DeleteAppKeys(ctx context.Context, appId string, id []string) (affected int64, err error) {
	deleted := c.Session(ctx).Model(&models.AppKey{}).Where("id in ? and app_id = ?", id, appId).Update("delete_time", time.Now().UTC())
	if err = deleted.Error; err != nil {
		return deleted.RowsAffected, err
	}
	return deleted.RowsAffected, nil
}

func (c CommonService) GetAppKeys(ctx context.Context, appId string, current, pageSize int64) (count int64, keyPairs []*models.AppKey, err error) {
	query := c.Session(ctx).Model(&models.AppKey{}).Where("app_id = ? and delete_time is NULL", appId)
	if err = query.Count(&count).Error; err != nil || count == 0 {
		return 0, nil, err
	} else if err = query.Select("id", "name", "create_time", "key").
		Order("id").Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).
		Find(&keyPairs).Error; err != nil {
		return 0, nil, err
	} else {
		for _, keyPair := range keyPairs {
			keyPair.AppId = appId
		}
		return count, keyPairs, nil
	}
}

func (c CommonService) GetAppKeyFromKey(ctx context.Context, key string) (*models.AppKey, error) {
	var appKey models.AppKey
	if err := c.Session(ctx).Model(&models.AppKey{}).Where("`key` = ? ", key).First(&appKey).Error; err != nil {
		return nil, err
	}
	return &appKey, nil
}

func (c CommonService) AppAuthorization(ctx context.Context, key string, secret string) (id string, err error) {
	conn := c.Session(ctx)
	var appKey models.AppKey
	if err = conn.Select("id", "app_id").
		Where("`key` = ? and secret = ?", key, sign.SumSha256Hmac(key, secret)).First(&appKey).Error; err != nil {
		return "", err
	}
	return appKey.AppId, nil
}

func (c *CommonService) CreateTOTP(ctx context.Context, id string, secret string) error {
	tx := c.Session(ctx).Begin()
	defer tx.Rollback()
	ext := new(models.UserExt)
	if err := tx.Where("user_id = ?", id).First(&ext).Error; err == gogorm.ErrRecordNotFound {

		totp := models.UserExt{UserId: id, TOTPAsMFA: true}
		err = totp.SetSecret(secret)
		if err != nil {
			return err
		}
		if err = tx.Create(&totp).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	sec, err := ext.GetSecret()
	if err != nil || sec != secret {
		if err = ext.SetSecret(secret); err != nil {
			return err
		}
	}
	ext.TOTPAsMFA = true
	if err = tx.Where("user_id = ?", ext.UserId).
		Select("totp_salt", "totp_secret", "totp_as_mfa").Updates(ext).Error; err != nil {
		return errors.NewServerError(500, "failed to update totp setting: "+err.Error())
	}
	return tx.Commit().Error
}

func (c *CommonService) GetTOTPSecrets(ctx context.Context, ids []string) (secrets []string, err error) {
	conn := c.Session(ctx)
	var totps []models.UserExt
	err = conn.Where("user_id in ?", ids).Find(&totps).Error
	if err != nil {
		return nil, err
	}
	for _, totp := range totps {
		secret, err := totp.GetSecret()
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secret)
	}
	return secrets, nil
}

func (c CommonService) PatchSystemConfig(ctx context.Context, prefix string, patch map[string]interface{}) error {
	tx := c.Session(ctx).Begin()
	defer tx.Rollback()
	for name, value := range patch {
		fullName := name
		if len(prefix) != 0 {
			fullName = fmt.Sprintf("%s.%s", prefix, name)
		}
		switch value.(type) {
		case string, uint, uint64, uint32, uint16, uint8, int, int64, int32, int16, int8, bool, float64, float32:
		default:
			return fmt.Errorf("unknown value type: %T", value)
		}
		valType := fmt.Sprintf("%T", value)
		var option models.SystemConfig
		if err := tx.Where("name = ?", fullName).First(&option).Error; err != nil {
			if err != gogorm.ErrRecordNotFound {
				return err
			} else if err = tx.Create(&models.SystemConfig{Name: fullName, Type: valType, Value: fmt.Sprintf("%v", value)}).Error; err != nil {
				return err
			}
			continue
		}
		if err := tx.Model(&models.SystemConfig{}).Where("name = ?", fullName).Updates(map[string]interface{}{
			"value": fmt.Sprintf("%v", value),
			"type":  valType,
		}).Error; err != nil {
			return err
		}
	}
	return tx.Commit().Error
}

func (c CommonService) GetSystemConfig(ctx context.Context, prefix string) (map[string]interface{}, error) {
	conn := c.Session(ctx)
	var options []models.SystemConfig
	var count int64
	query := conn.Model(&models.SystemConfig{})
	if len(prefix) != 0 {
		query = conn.Where("name like ?", prefix+".%")
	}
	if err := query.Model(&models.SystemConfig{}).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 2000 {
		return nil, fmt.Errorf("There are too many configurations, please check. ")
	}
	if err := query.Limit(2000).Find(&options).Error; err != nil {
		return nil, err
	}

	cfgMap := map[string]interface{}{}
	for _, option := range options {
		name := option.Name[len(prefix)+1:]
		switch option.Type {
		case "string":
			cfgMap[name] = option.Value
		case "float64":
			if val, err := strconv.ParseFloat(option.Value, 64); err == nil {
				cfgMap[name] = val
			}
		case "float32":
			if val, err := strconv.ParseFloat(option.Value, 32); err == nil {
				cfgMap[name] = val
			}
		case "uint":
			if val, err := strconv.ParseUint(option.Value, 10, 32); err == nil {
				cfgMap[name] = uint(val)
			}
		case "uint64":
			if val, err := strconv.ParseUint(option.Value, 10, 64); err == nil {
				cfgMap[name] = val
			}
		case "uint32":
			if val, err := strconv.ParseUint(option.Value, 10, 32); err == nil {
				cfgMap[name] = uint32(val)
			}
		case "uint16":
			if val, err := strconv.ParseUint(option.Value, 10, 16); err == nil {
				cfgMap[name] = uint16(val)
			}
		case "uint8":
			if val, err := strconv.ParseUint(option.Value, 10, 8); err == nil {
				cfgMap[name] = uint8(val)
			}
		case "int":
			if val, err := strconv.ParseInt(option.Value, 10, 32); err == nil {
				cfgMap[name] = int(val)
			}
		case "int64":
			if val, err := strconv.ParseInt(option.Value, 10, 64); err == nil {
				cfgMap[name] = int64(val)
			}
		case "int32":
			if val, err := strconv.ParseInt(option.Value, 10, 32); err == nil {
				cfgMap[name] = int32(val)
			}
		case "int16":
			if val, err := strconv.ParseInt(option.Value, 10, 16); err == nil {
				cfgMap[name] = int16(val)
			}
		case "int8":
			if val, err := strconv.ParseInt(option.Value, 10, 8); err == nil {
				cfgMap[name] = int8(val)
			}
		case "bool":
			if val, err := strconv.ParseBool(option.Value); err == nil {
				cfgMap[name] = val
			}
		}
	}
	return cfgMap, nil
}

func (c CommonService) BatchPatchI18n(ctx context.Context, i18ns []models.I18nTranslate) (err error) {
	conn := c.Session(ctx)
	for _, i18n := range i18ns {
		oriI18n := i18n
		if err = conn.Where(models.I18nTranslate{Source: i18n.Source, Field: i18n.Field, SourceId: i18n.SourceId, Lang: i18n.Lang}).Order("id desc").FirstOrCreate(&oriI18n).Error; err != nil {
			return err
		}
		if oriI18n.Value != i18n.Value {
			if err = conn.Model(&models.I18nTranslate{}).Where("id = ?", oriI18n.Id).Updates(&models.I18nTranslate{Value: i18n.Value}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CommonService) DeleteAppProxy(ctx context.Context, ids ...string) error {
	if deleted := c.Session(ctx).Model(&models.AppRoleURL{}).
		Where("app_id in ?", ids).
		Update("delete_time", time.Now().UTC()); deleted.Error != nil {
		return deleted.Error
	}
	if deleted := c.Session(ctx).Model(&models.AppProxyUrl{}).
		Where("app_id in ?", ids).
		Update("delete_time", time.Now().UTC()); deleted.Error != nil {
		return deleted.Error
	}
	deleted := c.Session(ctx).Model(&models.AppProxy{}).Where("app_id in ?", ids).Update("delete_time", time.Now().UTC())
	return deleted.Error
}

func (c *CommonService) DeleteAppAccessControl(ctx context.Context, ids ...string) error {
	if deleted := c.Session(ctx).Model(&models.AppUser{}).
		Where("app_id in ?", ids).
		Update("delete_time", time.Now().UTC()); deleted.Error != nil {
		return deleted.Error
	}
	if deleted := c.Session(ctx).Model(&models.AppRole{}).
		Where("app_id in ?", ids).
		Update("delete_time", time.Now().UTC()); deleted.Error != nil {
		return deleted.Error
	}
	return nil
}

func (c *CommonService) DeleteI18nBySourceId(ctx context.Context, ids ...string) error {
	if deleted := c.Session(ctx).Model(&models.I18nTranslate{}).
		Where("source_id in ?", ids).
		Update("delete_time", time.Now().UTC()); deleted.Error != nil {
		return deleted.Error
	}
	return nil
}

func (c *CommonService) GetI18n(ctx context.Context, source string, sourceId string, field string) (map[string]string, error) {
	var i18ns []models.I18nTranslate
	err := c.Session(ctx).Model(&models.I18nTranslate{}).
		Where("`source` = ? and `source_id` = ? and `field` = ?", source, sourceId, field).
		Find(&i18ns).Error
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string, len(i18ns))
	for _, i18n := range i18ns {
		ret[i18n.Lang] = i18n.Value
	}
	return ret, nil
}

func (c CommonService) PatchAppOAuthConfig(ctx context.Context, auth2 *models.AppOAuth2) error {
	query := c.Session(ctx).Model(&models.AppOAuth2{}).Where("app_id = ?", auth2.AppId)
	var oriOAuthConfig models.AppOAuth2
	if err := query.First(&oriOAuthConfig).Error; err != nil {
		query.Error = nil
		if errors.Is(err, gogorm.ErrRecordNotFound) {
			return query.Create(&auth2).Error
		}
		return err
	}
	if auth2.JwtSignatureMethod == models.AppMeta_default {
		auth2.JwtSignatureKey = &safe.String{}
	} else if oriOAuthConfig.JwtSignatureMethod != auth2.JwtSignatureMethod && (auth2.JwtSignatureKey == nil || auth2.JwtSignatureKey.String() == "") {
		auth2.JwtSignatureKey = safe.NewEncryptedString("", config.Get().GetSecurity().GetSecret())
		if oriOAuthConfig.JwtSignatureKey != nil && oriOAuthConfig.JwtSignatureKey.String() != "" {
			if privKey, err := oriOAuthConfig.JwtSignatureKey.UnsafeString(); err == nil {
				if _, err = jwt.NewJWTIssuer(auth2.AppId, auth2.JwtSignatureMethod.String(), privKey); err == nil {
					auth2.JwtSignatureKey = nil
				}
			}
		}
		if auth2.JwtSignatureKey != nil {
			key, err := jwt.NewRandomKey(auth2.JwtSignatureMethod.String())
			if err != nil {
				return err
			}
			if err := auth2.JwtSignatureKey.SetValue(key); err != nil {
				return err
			}
		}
	} else if oriOAuthConfig.JwtSignatureMethod == auth2.JwtSignatureMethod && (auth2.JwtSignatureKey == nil || auth2.JwtSignatureKey.String() == "") {
		auth2.JwtSignatureKey = nil
	}
	return query.Omit("app_id").Updates(auth2).Error
}

func (c CommonService) GetAppOAuthConfig(ctx context.Context, id string) (*models.AppOAuth2, error) {
	var auth2 models.AppOAuth2
	if err := c.Session(ctx).Model(&models.AppOAuth2{}).Where("app_id = ?", id).First(&auth2).Error; err != nil && !errors.Is(err, gogorm.ErrRecordNotFound) {
		return nil, err
	}
	if auth2.JwtSignatureKey != nil {
		auth2.JwtSignatureKey.SetSecret(config.Get().GetSecurity().GetSecret())
	}
	return &auth2, nil
}
