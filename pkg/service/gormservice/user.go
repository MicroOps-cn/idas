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

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/log/level"
	uuid "github.com/satori/go.uuid"
	gogorm "gorm.io/gorm"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

// ResetPassword
//
//	@Description[en-US]: Reset User Password.
//	@Description[zh-CN]: 重置用户密码。
//	@param ctx       context.Context
//	@param id        string
//	@param password  string           : New password.
//	@return err      error
func (s UserAndAppService) ResetPassword(ctx context.Context, ids string, password string) error {
	conn := s.Session(ctx).Begin()
	defer conn.Callback()
	for _, id := range strings.Split(ids, ",") {
		u := models.User{Model: models.Model{Id: id}, Salt: uuid.NewV4().Bytes(), Status: models.UserMeta_normal}
		u.Password = u.GenSecret(password)
		if err := conn.Select("password", "salt", "status").Where("status not in ?", []models.UserMeta_UserStatus{
			models.UserMeta_disabled,
		}).Updates(&u).Error; err != nil {
			return err
		}
	}

	return conn.Commit().Error
}

// UpdateLoginTime
//
//	@Description[en-US]: Update the user's last login time.
//	@Description[zh-CN]: 更新用户最后一次登陆时间。
//	@param ctx 	context.Context
//	@param id 	string
//	@return error
func (s UserAndAppService) UpdateLoginTime(ctx context.Context, id string) error {
	return s.Session(ctx).Model(&models.User{Model: models.Model{Id: id}}).UpdateColumn("login_time", time.Now().UTC()).Error
}

func (s UserAndAppService) Name() string {
	return s.name
}

const sqlGetUserAndRoleInfoById = `
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
    T1.id = ?
    AND T3.name = 'IDAS'
`

// VerifyPasswordById
//
//	@Description[en-US]: Verify the user's password through ID.
//	@Description[zh-CN]: 通过ID验证用户密码。
//	@param ctx 	context.Context
//	@param id 	string
//	@param password 	string
//	@return users	[]*models.User
func (s UserAndAppService) VerifyPasswordById(ctx context.Context, id, password string) *models.User {
	logger := logs.GetContextLogger(ctx)
	var user models.User
	if err := s.Session(ctx).Raw(sqlGetUserAndRoleInfoById, id).First(&user).Error; err != nil {
		if err == gogorm.ErrRecordNotFound {
			level.Debug(logger).Log("msg", "incorrect username", "id", id)
		} else {
			level.Error(logger).Log("msg", "unknown error", "id", id, "err", err)
		}
		return nil
	}
	if !bytes.Equal(user.GenSecret(password), user.Password) {
		level.Debug(logger).Log("msg", "incorrect password", "id", id)
		return nil
	}
	return &user
}

// VerifyPassword
//
//	@Description[en-US]: Verify password for user.
//	@Description[zh-CN]: 验证用户密码。
//	@param ctx 	context.Context
//	@param username 	string
//	@param password 	string
//	@return users	[]*models.User
func (s UserAndAppService) VerifyPassword(ctx context.Context, username string, password string) *models.User {
	logger := logs.GetContextLogger(ctx)
	var user models.User
	if err := s.Session(ctx).Where("(username = ? or email = ?) and delete_time is NULL", username, username).First(&user).Error; err != nil {
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
	return &user
}

// GetUserInfoByUsernameAndEmail
//
//	@Description[en-US]: Use username or email to obtain user information.
//	@Description[zh-CN]: 使用用户名或email获取用户信息。
//	@param ctx           context.Context
//	@param username      string
//	@param email         string
//	@return user   *models.User
//	@return err          error
func (s UserAndAppService) GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (user *models.User, err error) {
	user = new(models.User)
	query := s.Session(ctx).Where("username = ? and email = ? and delete_time is NULL", username, email)
	if err = query.First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUsers
//
//	@Description[en-US]: Get user list.
//	@Description[zh-CN]: 获取用户列表。
//	@param ctx       context.Context
//	@param keywords  string
//	@param status    models.UserMeta_UserStatus
//	@param appId     string
//	@param current   int64
//	@param pageSize  int64
//	@return total    int64
//	@return users    []*models.User
//	@return err      error
func (s UserAndAppService) GetUsers(ctx context.Context, keywords string, status models.UserMeta_UserStatus, appId string, current, pageSize int64) (total int64, users []*models.User, err error) {
	query := s.Session(ctx).Where("t_user.delete_time is NULL").Model(&models.User{})
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
	if status != models.UserMetaStatusAll {
		query = query.Where("status", status)
	}
	if err = query.Count(&total).Error; err != nil || total == 0 {
		return 0, nil, err
	} else if err = query.Order("username,id").Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&users).Error; err != nil {
		return 0, nil, err
	}
	return total, users, nil
}

// PatchUsers
//
//	@Description[en-US]: Incrementally update information of multiple users.
//	@Description[zh-CN]: 增量更新多个用户的信息。
//	@param ctx 		context.Context
//	@param patch 	[]map[string]interface{}
//	@return count	int64
//	@return err		error
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

// DeleteUsers
//
//	@Description[en-US]: Delete users in batch.
//	@Description[zh-CN]: 批量删除用户。
//	@param ctx 		context.Context
//	@param ids 		[]string
//	@return count	int64
//	@return err		error
func (s UserAndAppService) DeleteUsers(ctx context.Context, id []string) (int64, error) {
	deleted := s.Session(ctx).Model(&models.User{}).Where("id in ?", id).Update("delete_time", time.Now())
	if err := deleted.Error; err != nil {
		return deleted.RowsAffected, err
	}
	return deleted.RowsAffected, nil
}

// UpdateUser
//
//	@Description[en-US]: Update user information.
//	@Description[zh-CN]: 更新用户信息.
//	@param ctx	context.Context
//	@param user	*models.User
//	@param updateColumns	...string
//	@return err	error
func (s UserAndAppService) UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (err error) {
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()
	q := tx.Omit("create_time")
	if len(updateColumns) != 0 {
		q = q.Select(updateColumns)
	} else {
		q = q.Select("email", "phone_number", "full_name", "avatar", "status")
	}

	if err = q.Updates(&user).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// GetUserInfo
//
//	@Description[en-US]: Obtain user information through ID or username.
//	@Description[zh-CN]: 通过ID或用户名获取用户信息。
//	@param ctx 	context.Context
//	@param id 	string
//	@param username 	string
//	@return userDetail	*models.User
//	@return err	error
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

func (s UserAndAppService) GetUsersById(ctx context.Context, id []string) (users models.Users, err error) {
	conn := s.Session(ctx)
	query := conn.Model(&models.User{}).Where("id in ?", id)
	if err = query.Find(&users).Error; err != nil {
		if err == gogorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return
}

// CreateUser
//
//	@Description[en-US]: Create a user.
//	@Description[zh-CN]: 创建用户。
//	@param ctx 	context.Context
//	@param user 	*models.User
//	@return err	error
func (s UserAndAppService) CreateUser(ctx context.Context, user *models.User) (err error) {
	conn := s.Session(ctx)
	if len(user.Password) != 0 {
		user.Salt = uuid.NewV4().Bytes()
		user.Password = user.GenSecret()
	}
	return conn.Omit("role", "role_id").Create(user).Error
}

// PatchUser
//
//	@Description[en-US]: Incremental update user.
//	@Description[zh-CN]: 增量更新用户。
//	@param ctx 	context.Context
//	@param user 	map[string]interface{}
//	@return err	error
func (s UserAndAppService) PatchUser(ctx context.Context, patch map[string]interface{}) (err error) {
	if id, ok := patch["id"].(string); ok {
		tx := s.Session(ctx).Begin()
		delete(patch, "username")
		if err = tx.Model(&models.User{}).Where("id = ?", id).Updates(patch).Error; err != nil {
			return err
		}
		tx.Commit()
		return nil
	}
	return errors.ParameterError("id is null")
}

// DeleteUser
//
//	@Description[en-US]: Delete a user.
//	@Description[zh-CN]: 删除用户。
//	@param ctx 	context.Context
//	@param id 	string
//	@return error
func (s UserAndAppService) DeleteUser(ctx context.Context, id string) (err error) {
	_, err = s.DeleteUsers(ctx, []string{id})
	return err
}

func (c *CommonService) GetUserExtendedData(ctx context.Context, id string) (*models.UserExt, error) {
	conn := c.Session(ctx)
	var ext models.UserExt
	err := conn.Where("user_id = ?", id).First(&ext).Error
	if err == gogorm.ErrRecordNotFound {
		return nil, nil
	}
	if ext.EmailAsMFA || ext.SmsAsMFA || ext.TOTPAsMFA {
		ext.ForceMFA = true
	}
	return &ext, err
}
