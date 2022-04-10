package gormservice

import (
	"bytes"
	"context"
	"fmt"
	"idas/pkg/client/gorm"
	"reflect"

	gogorm "gorm.io/gorm"

	"idas/pkg/errors"
	"idas/pkg/service/models"
)

type UserService struct {
	*gorm.Client
	name string
}

func (s UserService) Name() string {
	return s.name
}

func (s UserService) AutoMigrate(ctx context.Context) error {
	return s.Session(ctx).AutoMigrate(&models.User{})
}

func (s UserService) VerifyPassword(ctx context.Context, username string, password string) (*models.User, error) {
	var user models.User
	if err := s.Session(ctx).Where("username = ? or email = ?", username, username).First(&user).Error; err != nil {
		return nil, err
	}

	if bytes.Equal(user.GenSecret(password), user.Password) {
		return nil, fmt.Errorf("invalid username or password")
	}
	return &user, nil
}

func NewUserService(name string, client *gorm.Client) *UserService {
	return &UserService{name: name, Client: client}
}

func (s UserService) GetUsers(ctx context.Context, keyword string, status models.UserStatus, current int64, pageSize int64) (users []*models.User, total int64, err error) {

	query := s.Session(ctx).Where("is_delete = 0")
	if len(keyword) > 0 {
		keyword = fmt.Sprintf("%%%s%%", keyword)
		query = query.Where(
			query.Where("username like ?", keyword).
				Or("email like ?", keyword).
				Or("phone_number like ?", keyword).
				Or("fullname like ?", keyword),
		)
	}
	if status != models.UserStatusUnknown {
		query = query.Where("status", status)
	}
	if err = query.Order("username,id").Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&users).Error; err != nil {
		return nil, 0, err
	} else if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	} else {
		for _, user := range users {
			user.Storage = s.name
		}
		return users, total, nil
	}
}

func (s UserService) PatchUsers(ctx context.Context, patch []map[string]interface{}) (int64, error) {
	var patchCount int64
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()
	updateQuery := tx.Model(&models.User{}).Select("is_delete", "status")
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
			if err := patched.Error; err != nil {
				return 0, err
			}
			patchCount = patched.RowsAffected
			newPatchIds = []string{}
			newPatch = map[string]interface{}{}
		}
	}
	if len(newPatchIds) > 0 {
		patched := updateQuery.Where("id in ?", newPatchIds).Updates(newPatch)
		if err := patched.Error; err != nil {
			return 0, err
		}
		patchCount = patched.RowsAffected
	}
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	return patchCount, nil
}

func (s UserService) DeleteUsers(ctx context.Context, id []string) (int64, error) {
	deleted := s.Session(ctx).Model(&models.User{}).Where("id in ?", id).Update("is_delete", true)
	if err := deleted.Error; err != nil {
		return deleted.RowsAffected, err
	}
	return deleted.RowsAffected, nil
}

func (s UserService) UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (*models.User, error) {
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()
	q := tx.Omit("create_time")
	if len(updateColumns) != 0 {
		q = q.Select(updateColumns)
	} else {
		q = q.Select("username", "email", "phone_number", "full_name", "avatar", "status")
	}

	if err := q.Updates(&user).Error; err != nil {
		return nil, err
	}
	if err := tx.Find(&user).Error; err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s UserService) GetUserInfo(ctx context.Context, id string, username string) (*models.User, error) {
	conn := s.Session(ctx)
	var user models.User
	query := conn.Model(&models.User{})
	if len(id) != 0 && len(username) != 0 {
		subQuery := query.Where("id = ?", id).Or("username = ?", username)
		query = query.Where(subQuery)
	} else if len(id) != 0 {
		query = query.Where("id = ?", id)
	} else if len(username) != 0 {
		query = query.Where("username = ?", username)
	} else {
		return nil, errors.ParameterError("require id or username")
	}
	if err := query.First(&user).Error; err != nil {
		if err == gogorm.ErrRecordNotFound {
			return nil, errors.StatusNotFound("user")
		}
		return nil, err
	}
	return &user, nil
}

func (s UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	conn := s.Session(ctx)
	if err := conn.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s UserService) PatchUser(ctx context.Context, patch map[string]interface{}) (*models.User, error) {
	if id, ok := patch["id"].(string); ok {
		tx := s.Session(ctx).Begin()
		user := models.User{Model: models.Model{Id: id}}
		if err := tx.Model(&models.User{}).Where("id = ?", id).Updates(patch).Error; err != nil {
			return nil, err
		} else if err = tx.First(&user).Error; err != nil {
			return nil, err
		}
		tx.Commit()
		return &user, nil
	}
	return nil, errors.ParameterError("id is null")
}

func (s UserService) DeleteUser(ctx context.Context, id string) (err error) {
	_, err = s.DeleteUsers(ctx, []string{id})
	return err
}
