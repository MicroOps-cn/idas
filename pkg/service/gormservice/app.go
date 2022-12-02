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
	"reflect"

	gogorm "gorm.io/gorm"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

// PatchApps
//
//	@Description[en-US]: Incrementally update information of multiple applications.
//	@Description[zh-CN]: 增量更新多个应用的信息。
//	@param ctx     context.Context
//	@param patch   []map[string]interface{}
//	@return total  int64
//	@return err    error
func (s UserAndAppService) PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error) {
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()
	updateQuery := tx.Model(&models.App{}).Select("is_delete", "status")
	var newPatch map[string]interface{}
	var newPatchIds []string
	for _, patchInfo := range patch {
		tmpPatch := map[string]interface{}{}
		var tmpPatchId string
		for name, value := range patchInfo {
			if name != "id" {
				tmpPatch[name] = value
			} else {
				tmpPatchId, _ = value.(string)
			}
		}
		if tmpPatchId == "" {
			return 0, fmt.Errorf("parameter exception: invalid id")
		} else if len(tmpPatch) == 0 {
			return 0, fmt.Errorf("parameter exception: update content is empty")
		}
		if len(newPatchIds) == 0 {
			newPatchIds = append(newPatchIds, tmpPatchId)
			newPatch = tmpPatch
		} else if reflect.DeepEqual(tmpPatch, newPatch) {
			newPatchIds = append(newPatchIds, tmpPatchId)
		} else {
			patched := updateQuery.Where("id in ?", newPatchIds).Updates(newPatch)
			if err = patched.Error; err != nil {
				return 0, err
			}
			total = patched.RowsAffected
			newPatchIds = []string{}
			newPatch = map[string]interface{}{}
		}
	}
	if len(newPatchIds) > 0 {
		patched := updateQuery.Where("id in ?", newPatchIds).Updates(newPatch)
		if err = patched.Error; err != nil {
			return 0, err
		}
		total = patched.RowsAffected
	}
	if err = tx.Commit().Error; err != nil {
		return 0, err
	}
	return total, nil
}

// DeleteApps
//
//	@Description[en-US]: Delete apps in batch.
//	@Description[zh-CN]: 批量删除应用。
//	@param ctx     context.Context
//	@param ids     []string         : ID List
//	@return total  int64
//	@return err    error
func (s UserAndAppService) DeleteApps(ctx context.Context, id []string) (total int64, err error) {
	deleted := s.Session(ctx).Model(&models.App{}).Where("id in ?", id).Update("is_delete", true)
	if err = deleted.Error; err != nil {
		return deleted.RowsAffected, err
	}
	return deleted.RowsAffected, nil
}

// PatchAppRole
//
//	@Description[en-US]: Update App Role.
//	@Description[zh-CN]: 更新应用角色。
//	@param ctx     context.Context
//	@param dn      string
//	@param patch   *models.AppRole
//	@return err    error
func (s UserAndAppService) PatchAppRole(ctx context.Context, role *models.AppRole) error {
	conn := s.Session(ctx)

	var r models.AppRole
	if len(role.Id) != 0 {
		if err := conn.Where("id = ? and app_id = ?", role.Id, role.AppId).First(&r).Error; err == gogorm.ErrRecordNotFound {
			role.Id = ""
			if err = conn.Create(&role).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if role.Name != r.Name || role.IsDefault != r.IsDefault || role.IsDelete != r.IsDelete {
			if err = conn.Select("name", "config", "is_delete", "is_default").Updates(&role).Error; err != nil {
				return err
			}
		}
	} else {
		if err := conn.Create(&role).Error; err != nil {
			return err
		}
	}
	var userIds []string
	for _, user := range role.Users {
		user.RoleId = role.Id
		appUser := models.AppUser{AppId: role.AppId, UserId: user.Id, RoleId: role.Id}
		var oldAppUser models.AppUser
		if err := conn.Where("app_id = ? and user_id = ?", role.AppId, user.Id).First(&oldAppUser).Error; err == gogorm.ErrRecordNotFound {
			if err = conn.Create(&appUser).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if oldAppUser.RoleId != role.Id {
			if err = conn.Model(&models.AppUser{}).Where("app_id = ? and user_id = ?", role.AppId, user.Id).Update("role_id", role.Id).Error; err != nil {
				return err
			}
		}
		userIds = append(userIds, user.Id)
	}
	return conn.Delete(&models.AppUser{}, "app_id = ? and role_id = ? and user_id not in ? ", role.AppId, role.Id, userIds).Error
}

// UpdateApp
//
//	@Description[en-US]: Update applies the value of the specified column. If no column is specified, all column information is updated.
//	@Description[zh-CN]: 更新应用指定列的值，如果未指定列，则表示更新所有列信息。
//	@param ctx           context.Context
//	@param app           *models.App
//	@param updateColumns ...string
//	@return newApp       *models.App
//	@return err          error
func (s UserAndAppService) UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (*models.App, error) {
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()
	q := tx.Omit("create_time")

	if app.GrantType != models.AppMeta_proxy {
		app.Proxy = nil
	}

	if len(updateColumns) != 0 {
		q = q.Select(updateColumns)
	} else {
		q = q.Select("name", "description", "avatar", "grant_type", "grant_mode", "status")
	}

	if err := q.Updates(&app).Error; err != nil {
		return nil, err
	}

	if len(app.Roles) > 0 {
		var roleIds []string
		for _, role := range app.Roles {
			for _, user := range app.Users {
				if len(role.Id) != 0 {
					if user.RoleId == role.Id || (user.RoleId == "" && role.IsDefault) {
						role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
					}
				} else if len(role.Name) != 0 && string(user.Role) == role.Name {
					role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
				}
			}
			role.AppId = app.Id
			if err := s.PatchAppRole(context.WithValue(ctx, global.GormConnName, tx), role); err != nil {
				return nil, err
			}
			roleIds = append(roleIds, role.Id)
		}
		if err := tx.Delete(&models.AppRole{}, "app_id = ? and id not in ? ", app.Id, roleIds).Error; err != nil {
			return nil, err
		}
	}

	if err := tx.Find(&app).Error; err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return app, nil
}

// GetAppInfo
//
//	@Description[en-US]: Use the ID or application name to get app info.
//	@Description[zh-CN]: 使用ID或应用名称获取应用信息。
//	@param ctx  context.Context
//	@param id   string          : App ID
//	@param name string          : App Name
//	@return app *models.App     : App Details
//	@return err error
func (s UserAndAppService) GetAppInfo(ctx context.Context, id string, name string) (app *models.App, err error) {
	conn := s.Session(ctx)
	app = new(models.App)
	query := conn.Model(&models.App{})
	if len(id) != 0 && len(name) != 0 {
		subQuery := query.Where("id = ?", id).Or("name = ?", name)
		query = query.Where(subQuery)
	} else if len(id) != 0 {
		query = query.Where("id = ?", id)
	} else if len(name) != 0 {
		query = query.Where("name = ?", name)
	} else {
		return nil, errors.ParameterError("require id or name")
	}
	if err = query.First(&app).Error; err != nil {
		return nil, err
	}
	//if err = conn.Model(&app).Association("User").Find(&app.User); err != nil {
	//	return nil, err
	//}
	if err = conn.Model(&models.User{}).Select("`t_user`.`id`,`t_user`.`create_time`,`t_user`.`update_time`,`t_user`.`is_delete`,`t_user`.`username`,`t_user`.`salt`,`t_user`.`password`,`t_user`.`email`,`t_user`.`phone_number`,`t_user`.`full_name`,`t_user`.`avatar`,`t_user`.`status`,`t_user`.`login_time`, `t_app_role`.`name` as role, `t_app_role`.`id` as role_id").
		Joins("JOIN `t_app_user` ON `t_app_user`.`user_id` = `t_user`.`id` AND `t_app_user`.`app_id` = ?", app.Id).
		Joins("JOIN `t_app_role` ON `t_app_user`.`role_id` = `t_app_role`.`id`").Find(&app.Users).Error; err != nil {
		return nil, err
	}
	if err = conn.Model(&app).Association("Roles").Find(&app.Roles); err != nil {
		return nil, err
	}
	return
}

// CreateApp
//
//	@Description[en-US]: Create an app.
//	@Description[zh-CN]: 创建应用
//	@param ctx        context.Context
//	@param app        *models.App
//	@return appDetail *models.App
//	@return error
func (s UserAndAppService) CreateApp(ctx context.Context, app *models.App) (*models.App, error) {
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()

	if app.GrantType != models.AppMeta_proxy {
		app.Proxy = nil
	}

	if err := tx.Omit("Users").Create(app).Error; err != nil {
		return nil, err
	}
	if len(app.Roles) > 0 {
		for _, role := range app.Roles {
			for _, user := range app.Users {
				if len(role.Id) != 0 {
					if user.RoleId == role.Id || (user.RoleId == "" && role.IsDefault) {
						role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
					}
				} else if len(role.Name) != 0 && user.Role == role.Name {
					role.Users = append(role.Users, &models.User{Model: models.Model{Id: user.Id}})
				}
			}
			role.AppId = app.Id
			if err := s.PatchAppRole(context.WithValue(ctx, global.GormConnName, tx), role); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Find(&app).Error; err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return app, nil
}

// PatchApp
//
//	@Description[en-US]: Incremental update application.
//	@Description[zh-CN]: 增量更新应用。
//	@param ctx        context.Context
//	@param fields     map[string]interface{}
//	@return appDetail app *models.App
//	@return err       error
func (s UserAndAppService) PatchApp(ctx context.Context, fields map[string]interface{}) (app *models.App, err error) {
	if id, ok := fields["id"].(string); ok {
		tx := s.Session(ctx).Begin()
		app = &models.App{Model: models.Model{Id: id}}
		if err = tx.Model(&models.User{}).Where("id = ?", id).Updates(fields).Error; err != nil {
			return nil, err
		} else if err = tx.First(app).Error; err != nil {
			return nil, err
		}
		tx.Commit()
		return app, nil
	}
	return nil, errors.ParameterError("id is null")
}

// DeleteApp
//
//	@Description[en-US]: Delete an app.
//	@Description[zh-CN]: 删除应用。
//	@param ctx 	context.Context
//	@param id 	string
//	@return err	error
func (s UserAndAppService) DeleteApp(ctx context.Context, id string) (err error) {
	_, err = s.DeleteApps(ctx, []string{id})
	return err
}

// GetApps
//
//	@Description[en-US]: Get the application list. The application information does not include agent, role, user and other information.
//	@Description[zh-CN]: 获取应用列表，应用信息中不包含代理、角色、用户等信息。
//	@param ctx       context.Context
//	@param keywords  string
//	@param current   int64
//	@param pageSize  int64
//	@return total    int64
//	@return apps     []*models.App
//	@return err      error
func (s UserAndAppService) GetApps(ctx context.Context, keywords string, current, pageSize int64) (total int64, apps []*models.App, err error) {
	query := s.Session(ctx).Model(&models.App{}).Where("is_delete = 0")
	if len(keywords) > 0 {
		keywords = fmt.Sprintf("%%%s%%", keywords)
		query = query.Where("name like ? or description like ?", keywords, keywords)
	}
	if err = query.Order("name,id").Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&apps).Error; err != nil {
		return 0, nil, err
	} else if err = query.Count(&total).Error; err != nil {
		return 0, nil, err
	} else {
		for _, app := range apps {
			app.Storage = s.name
		}
		return total, apps, nil
	}
}

type RoleResult struct {
	Role string
}

// VerifyUserAuthorizationForApp
//
//	@Description[en-US]: Verify user authorization for the application.
//	@Description[zh-CN]: 验证应用程序的用户授权
//	@param ctx    context.Context
//	@param appId  string
//	@param userId string
//	@return role  string   :Role name, such as admin, viewer, editor ...
//	@return err   error
func (s UserAndAppService) VerifyUserAuthorizationForApp(ctx context.Context, appId string, userId string) (role string, err error) {
	var result RoleResult
	if err = s.Session(ctx).Model(&models.AppUser{}).Select("t_app_role.name as `role`").
		Joins("JOIN `t_app_role` ON `t_app_role`.`id` = `t_app_user`.`role_id`").
		Where("t_app_user.app_id = ? AND t_app_user.user_id = ?", appId, userId).First(&result).Error; err == gogorm.ErrRecordNotFound {
		if err = s.Session(ctx).Model(&models.AppUser{}).Select("t_app_role.name as `role`").
			Joins(" LEFT JOIN `t_app_role` ON (`t_app_user`.`app_id` = `t_app_role`.`app_id`)").
			Where("`t_app_user`.`app_id` = ? and (`t_app_role`.`is_default` = 1 or `t_app_role`.`is_default` is null )", appId).First(&result).Error; err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	return result.Role, nil
}
