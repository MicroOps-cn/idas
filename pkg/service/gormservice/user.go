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
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-kit/log/level"
	uuid "github.com/satori/go.uuid"
	gogorm "gorm.io/gorm"

	"github.com/MicroOps-cn/idas/pkg/client/gorm"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/logs"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

type UserAndAppService struct {
	*gorm.Client
	name string
}

func (s UserAndAppService) ResetPassword(ctx context.Context, ids string, password string) error {
	conn := s.Session(ctx).Begin()
	defer conn.Callback()
	for _, id := range strings.Split(ids, ",") {
		u := models.User{Model: models.Model{Id: id}, Salt: uuid.NewV4().Bytes(), Status: models.UserMeta_normal}
		u.Password = u.GenSecret(password)
		if err := conn.Select("password", "salt", "status").Where("status in ?", []models.UserMeta_UserStatus{
			models.UserMeta_inactive, models.UserMeta_inactive, models.UserMeta_unknown,
		}).Updates(&u).Error; err != nil {
			return err
		}
	}

	return conn.Commit().Error
}

func (s UserAndAppService) UpdateLoginTime(ctx context.Context, id string) error {
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()
	return tx.Model(&models.User{Model: models.Model{Id: id}}).UpdateColumn("login_time", time.Now().UTC()).Error
}

func (s UserAndAppService) Name() string {
	return s.name
}

const sqlGetUserAndRoleInfo = `
SELECT 
    T4.id AS role_id, T4.name AS role, T1.*
FROM
    t_user T1
        LEFT JOIN
    t_app_user T2 ON T2.user_id = T1.id
        LEFT JOIN
    t_app T3 ON T3.id = T2.app_id 
        LEFT JOIN
    t_app_role T4 ON T2.role_id = T4.id
WHERE
    (T1.username = ? or T1.email = ?)
    AND T3.name = 'IDAS'
`

func (s UserAndAppService) VerifyPassword(ctx context.Context, username string, password string) []*models.User {
	logger := logs.GetContextLogger(ctx)
	var user models.User
	if err := s.Session(ctx).Raw(sqlGetUserAndRoleInfo, username, username).First(&user).Error; err != nil {
		if err == gogorm.ErrRecordNotFound {
			level.Debug(logger).Log("msg", "incorrect username", "username", username)
		} else {
			level.Error(logger).Log("msg", "unknown error", "username", username, "err", err)
		}
		return nil
	}
	if !bytes.Equal(user.GenSecret(password), user.Password) {
		level.Debug(logger).Log("msg", "incorrect password", "username", username)
		return nil
	}
	return []*models.User{&user}
}

func (s UserAndAppService) GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (user *models.User, err error) {
	user = new(models.User)
	query := s.Session(ctx).Where("username = ? and email = ? and is_delete = 0", username, email)
	if err = query.First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s UserAndAppService) GetUsers(ctx context.Context, keywords string, status models.UserMeta_UserStatus, appId string, current, pageSize int64) (total int64, users []*models.User, err error) {
	query := s.Session(ctx).Where("t_user.is_delete = 0")
	if len(keywords) > 0 {
		keywords = fmt.Sprintf("%%%s%%", keywords)
		query = query.Where(
			query.Where("username like ?", keywords).
				Or("email like ?", keywords).
				Or("phone_number like ?", keywords).
				Or("full_name like ?", keywords),
		)
	}
	if len(appId) != 0 {
		query = query.
			Joins("LEFT JOIN t_app_user ON t_app_user.user_id = t_user.id").
			Where("t_app_user.app_id = ?", appId)
	}
	if status != models.UserMeta_unknown {
		query = query.Where("status", status)
	}
	if err = query.Order("username,id").Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&users).Error; err != nil {
		return 0, nil, err
	} else if err = query.Count(&total).Error; err != nil {
		return 0, nil, err
	} else {
		for _, user := range users {
			user.Storage = s.name
		}
		return total, users, nil
	}
}

func (s UserAndAppService) PatchUsers(ctx context.Context, patch []map[string]interface{}) (int64, error) {
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

func (s UserAndAppService) DeleteUsers(ctx context.Context, id []string) (int64, error) {
	deleted := s.Session(ctx).Model(&models.User{}).Where("id in ?", id).Update("is_delete", true)
	if err := deleted.Error; err != nil {
		return deleted.RowsAffected, err
	}
	return deleted.RowsAffected, nil
}

func (s UserAndAppService) UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (*models.User, error) {
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

func (s UserAndAppService) GetUserInfo(ctx context.Context, id string, username string) (*models.User, error) {
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

func (s UserAndAppService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	conn := s.Session(ctx)
	if len(user.Password) != 0 {
		user.Salt = uuid.NewV4().Bytes()
		user.Password = user.GenSecret()
	}
	if err := conn.Omit("role", "role_id").Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s UserAndAppService) PatchUser(ctx context.Context, patch map[string]interface{}) (*models.User, error) {
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

func (s UserAndAppService) DeleteUser(ctx context.Context, id string) (err error) {
	_, err = s.DeleteUsers(ctx, []string{id})
	return err
}
