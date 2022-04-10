package gormservice

import (
	"context"
	"fmt"
	"idas/pkg/client/gorm"
	"reflect"

	"idas/pkg/errors"
	"idas/pkg/service/models"
)

type AppService struct {
	*gorm.Client
	name string
}

func (a AppService) PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error) {
	tx := a.Session(ctx).Begin()
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

func (a AppService) DeleteApps(ctx context.Context, id []string) (total int64, err error) {
	deleted := a.Session(ctx).Model(&models.App{}).Where("id in ?", id).Update("is_delete", true)
	if err = deleted.Error; err != nil {
		return deleted.RowsAffected, err
	}
	return deleted.RowsAffected, nil
}

func (a AppService) UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (*models.App, error) {
	tx := a.Session(ctx).Begin()
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
	if err := tx.Find(&app).Error; err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return app, nil
}

func (a AppService) GetAppInfo(ctx context.Context, id string, name string) (app *models.App, err error) {
	conn := a.Session(ctx)
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
	return
}

func (a AppService) CreateApp(ctx context.Context, app *models.App) (*models.App, error) {
	conn := a.Session(ctx)
	if err := conn.Create(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

func (a AppService) PatchApp(ctx context.Context, fields map[string]interface{}) (app *models.App, err error) {
	if id, ok := fields["id"].(string); ok {
		tx := a.Session(ctx).Begin()
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

func (a AppService) DeleteApp(ctx context.Context, id string) (err error) {
	_, err = a.DeleteApps(ctx, []string{id})
	return err
}

func (a AppService) Name() string {
	return a.name
}

func (a AppService) GetApps(ctx context.Context, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error) {
	query := a.Session(ctx).Model(&models.App{})
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
			app.Storage = a.name
		}
		return apps, total, nil
	}
}

func NewAppService(name string, client *gorm.Client) *AppService {
	if err := client.Session(context.Background()).SetupJoinTable(&models.App{}, "User", models.AppUser{}); err != nil {
		panic(err)
	}
	return &AppService{name: name, Client: client}
}

func (a AppService) AutoMigrate(ctx context.Context) error {
	return a.Session(ctx).AutoMigrate(&models.App{}, &models.AppUser{})
}
