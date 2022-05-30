package gormservice

import (
	"context"
	"fmt"
	gogorm "gorm.io/gorm"
	"idas/pkg/client/gorm"
	"idas/pkg/global"
	"reflect"

	"idas/pkg/errors"
	"idas/pkg/service/models"
)

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

func (s UserAndAppService) DeleteApps(ctx context.Context, id []string) (total int64, err error) {
	deleted := s.Session(ctx).Model(&models.App{}).Where("id in ?", id).Update("is_delete", true)
	if err = deleted.Error; err != nil {
		return deleted.RowsAffected, err
	}
	return deleted.RowsAffected, nil
}

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
		} else if role.Name != r.Name || role.Config != r.Config || role.IsDefault != r.IsDefault || role.IsDelete != r.IsDelete {
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
	for _, user := range role.User {
		user.RoleId = role.Id
		var appUser = models.AppUser{AppId: role.AppId, UserId: user.Id, RoleId: role.Id}
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

func (s UserAndAppService) UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (*models.App, error) {
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()
	q := tx.Omit("create_time")
	if len(updateColumns) != 0 {
		q = q.Select(updateColumns)
	} else {
		q = q.Select("name", "description", "avatar", "grant_type", "grant_mode", "status")
	}

	if err := q.Updates(&app).Error; err != nil {
		return nil, err
	}

	if len(app.Role) > 0 {
		var roleIds []string
		for _, role := range app.Role {
			for _, user := range app.User {
				if user.RoleId == role.Id || (user.RoleId == "" && role.IsDefault) {
					role.User = append(role.User, &models.User{Model: models.Model{Id: user.Id}})
				}
			}
			role.AppId = app.Id
			if err := s.PatchAppRole(context.WithValue(ctx, global.MySQLConnName, tx), role); err != nil {
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

func (s UserAndAppService) GetAppInfo(ctx context.Context, id string, name string) (app *models.App, err error) {
	conn := s.Session(ctx)
	app = new(models.App)
	query := conn.Model(&models.User{})
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
	if err = conn.Where("id = ?", id).First(&app).Error; err != nil {
		return nil, err
	}
	//if err = conn.Model(&app).Association("User").Find(&app.User); err != nil {
	//	return nil, err
	//}
	if err = conn.Model(&models.User{}).Select("`t_user`.`id`,`t_user`.`create_time`,`t_user`.`update_time`,`t_user`.`is_delete`,`t_user`.`username`,`t_user`.`salt`,`t_user`.`password`,`t_user`.`email`,`t_user`.`phone_number`,`t_user`.`full_name`,`t_user`.`avatar`,`t_user`.`status`,`t_user`.`login_time`, `t_app_role`.`name` as role, `t_app_role`.`id` as role_id").
		Joins("JOIN `t_app_user` ON `t_app_user`.`user_id` = `t_user`.`id` AND `t_app_user`.`app_id` = ?", app.Id).
		Joins("JOIN `t_app_role` ON `t_app_user`.`role_id` = `t_app_role`.`id`").Find(&app.User).Error; err != nil {
		return nil, err
	}
	if err = conn.Model(&app).Association("Role").Find(&app.Role); err != nil {
		return nil, err
	}
	return
}

func (s UserAndAppService) CreateApp(ctx context.Context, app *models.App) (*models.App, error) {
	conn := s.Session(ctx)
	if err := conn.Create(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

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

func (s UserAndAppService) DeleteApp(ctx context.Context, id string) (err error) {
	_, err = s.DeleteApps(ctx, []string{id})
	return err
}

func (s UserAndAppService) GetApps(ctx context.Context, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error) {
	query := s.Session(ctx).Model(&models.App{})
	if len(keywords) > 0 {
		keywords = fmt.Sprintf("%%%s%%", keywords)
		query = query.Where("name like ? or description like ?", keywords, keywords)
	}
	if err = query.Order("name,id").Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&apps).Error; err != nil {
		return nil, 0, err
	} else if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	} else {
		for _, app := range apps {
			app.Storage = s.name
		}
		return apps, total, nil
	}
}

type Scope struct {
	Scope string
}

func (s UserAndAppService) VerifyUserAuthorizationForApp(ctx context.Context, appId string, userId string) (scope string, err error) {
	var result Scope
	if err = s.Session(ctx).Model(&models.AppUser{}).Select("t_app_role.name as `scope`").
		Joins("JOIN `t_app_role` ON `t_app_role`.`id` = `t_app_user`.`role_id`").
		Where("t_app_user.app_id = ? AND t_app_user.user_id = ?", appId, userId).First(&result).Error; err == gogorm.ErrRecordNotFound {
		if err = s.Session(ctx).Model(&models.AppUser{}).Select("t_app_role.name as `scope`").
			Joins(" LEFT JOIN `t_app_role` ON (`t_app_user`.`app_id` = `t_app_role`.`app_id`)").
			Where("`t_app_user`.`app_id` = ? and (`t_app_role`.`is_default` = 1 or `t_app_role`.`is_default` is null )", appId).First(&result).Error; err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	return result.Scope, nil
}

func NewUserAndAppService(name string, client *gorm.Client) *UserAndAppService {
	conn := client.Session(context.Background())
	if err := conn.SetupJoinTable(&models.App{}, "User", models.AppUser{}); err != nil {
		panic(err)
	}
	return &UserAndAppService{name: name, Client: client}
}

func (s UserAndAppService) AutoMigrate(ctx context.Context) error {
	return s.Session(ctx).AutoMigrate(&models.AppUser{}, &models.AppRole{}, &models.User{}, &models.AppAuthCode{}, &models.App{})
}
