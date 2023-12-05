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
	"time"

	"gorm.io/gorm/clause"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
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
			return 0, errors.ParameterError("invalid id")
		} else if len(tmpPatch) == 0 {
			return 0, errors.ParameterError("update content is empty")
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
func (s UserAndAppService) DeleteApps(ctx context.Context, id ...string) (total int64, err error) {
	var idasCount int64
	if err = s.Session(ctx).Model(&models.App{}).Where("`id` in ? and `name` = 'IDAS'", id).Count(&idasCount).Error; err != nil {
		return 0, err
	}
	if idasCount > 0 {
		return 0, errors.NewServerError(400, "can't delete the idas app", errors.CodeAppCannotBeDelete)
	}
	deleted := s.Session(ctx).Model(&models.App{}).Where("`id` in ? and `name` != 'IDAS'", id).Update("delete_time", time.Now().UTC())
	if err = deleted.Error; err != nil {
		return deleted.RowsAffected, err
	}
	return deleted.RowsAffected, nil
}

// UpdateApp
//
//	@Description[en-US]: Update applies the value of the specified column. If no column is specified, all column information is updated.
//	@Description[zh-CN]: 更新应用指定列的值，如果未指定列，则表示更新所有列信息。
//	@param ctx           context.Context
//	@param app           *models.App
//	@param updateColumns ...string
//	@return err          error
func (s UserAndAppService) UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (err error) {
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()
	q := tx.Omit("create_time")

	if app.GrantType != models.AppMeta_proxy {
		app.Proxy = nil
	}

	if len(updateColumns) != 0 {
		q = q.Select(updateColumns)
	} else {
		q = q.Select("description", "display_name", "avatar", "grant_type", "grant_mode", "status", "url")
	}

	if err = q.Omit("name").Updates(&app).Error; err != nil {
		return err
	}

	return tx.Commit().Error
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
func (s UserAndAppService) GetAppInfo(ctx context.Context, options ...opts.WithGetAppOptions) (app *models.App, err error) {
	conn := s.Session(ctx)
	app = new(models.App)
	o := opts.NewAppOptions(options...)
	query := conn.Model(&models.App{})
	if len(o.Id) != 0 && len(o.Name) != 0 {
		subQuery := query.Where("id = ?", o.Id).Or("name = ?", o.Name)
		query = query.Where(subQuery)
	} else if len(o.Id) != 0 {
		query = query.Where("id = ?", o.Id)
	} else if len(o.Name) != 0 {
		query = query.Where("name = ?", o.Name)
	} else {
		return nil, errors.ParameterError("require id or name")
	}
	if err = query.First(&app).Error; err != nil {
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
//	@return error
func (s UserAndAppService) CreateApp(ctx context.Context, app *models.App) (err error) {
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()

	if app.GrantType != models.AppMeta_proxy {
		app.Proxy = nil
	}
	if err = tx.Omit("Users").Create(app).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

// PatchApp
//
//	@Description[en-US]: Incremental update application.
//	@Description[zh-CN]: 增量更新应用。
//	@param ctx        context.Context
//	@param fields     map[string]interface{}
//	@return err       error
func (s UserAndAppService) PatchApp(ctx context.Context, fields map[string]interface{}) (err error) {
	if id, ok := fields["id"].(string); ok {
		tx := s.Session(ctx).Begin()
		if err = tx.Model(&models.User{}).Omit("create_time", "name").Where("id = ?", id).Updates(fields).Error; err != nil {
			return err
		}
		return tx.Commit().Error
	}
	return errors.ParameterError("id is null")
}

// DeleteApp
//
//	@Description[en-US]: Delete an app.
//	@Description[zh-CN]: 删除应用。
//	@param ctx 	context.Context
//	@param id 	string
//	@return err	error
func (s UserAndAppService) DeleteApp(ctx context.Context, id string) (err error) {
	_, err = s.DeleteApps(ctx, id)
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
func (s UserAndAppService) GetApps(ctx context.Context, keywords string, filters map[string]interface{}, current, pageSize int64) (total int64, apps []*models.App, err error) {
	query := s.Session(ctx).Model(&models.App{})
	if len(keywords) > 0 {
		keywords = fmt.Sprintf("%%%s%%", keywords)
		query = query.Where("`t_app`.name like ? or `t_app`.description like ? or `t_app`.display_name like ?", keywords, keywords, keywords)
	}
	for name, val := range filters {
		switch name {
		case "user_id":
			query = query.Joins("LEFT JOIN `t_app_user` ON `t_app`.`id` = `t_app_user`.`app_id`").Where("`t_app_user`.`user_id` = ? AND `t_app_user`.`delete_time` IS NULL", val)
		case "url":
			if val == "*" {
				query = query.Where("`t_app`.url <> '' and `t_app`.url IS NOT NULL")
			}
		}
		query = query.Where(clause.Eq{Column: name, Value: val})
	}
	if err = query.Count(&total).Error; err != nil {
		return 0, nil, err
	} else if total > 0 {
		if err = query.Select("t_app.*").Order("`t_app`.name,`t_app`.id").Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&apps).Error; err != nil {
			return 0, nil, err
		}
	}
	return total, apps, nil
}
