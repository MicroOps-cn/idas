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
	"time"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/log/level"
	uuid "github.com/satori/go.uuid"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
	"github.com/MicroOps-cn/idas/pkg/utils/sign"
)

// GetUserSource
//
//	@Description[en-US]: Get user storage source.
//	@Description[zh-CN]: 获取用户存储源。
//	@param _ 	context.Context
//	@return total	int64
//	@return data	map[string]string
//	@return err	error
func (s Set) GetUserSource(_ context.Context) (total int64, data map[string]string, err error) {
	data = map[string]string{}
	for _, userService := range s.userAndAppService {
		data[userService.Name()] = userService.Name()
	}
	return
}

// GetUsers
//
//	@Description[en-US]: Get user list.
//	@Description[zh-CN]: 获取用户列表。
//	@param ctx       context.Context
//	@param storage 	string
//	@param keywords  string
//	@param status    models.UserMeta_UserStatus
//	@param appId     string
//	@param current   int64
//	@param pageSize  int64
//	@return total    int64
//	@return users    []*models.User
//	@return err      error
func (s Set) GetUsers(ctx context.Context, storage string, keywords string, status models.UserMeta_UserStatus, appId string, current, pageSize int64) (total int64, users models.Users, err error) {
	return s.SafeGetUserAndAppService(storage).GetUsers(ctx, keywords, status, appId, current, pageSize)
}

// PatchUsers
//
//	@Description[en-US]: Incrementally update information of multiple users.
//	@Description[zh-CN]: 增量更新多个用户的信息。
//	@param ctx 		context.Context
//	@param storage 	string
//	@param patch 	[]map[string]interface{}
//	@return count	int64
//	@return err		error
func (s Set) PatchUsers(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, err error) {
	return s.GetUserAndAppService(storage).PatchUsers(ctx, patch)
}

// DeleteUsers
//
//	@Description[en-US]: Delete users in batch.
//	@Description[zh-CN]: 批量删除用户。
//	@param ctx 		context.Context
//	@param storage 	string
//	@param ids 		[]string
//	@return count	int64
//	@return err		error
func (s Set) DeleteUsers(ctx context.Context, storage string, id []string) (total int64, err error) {
	return s.GetUserAndAppService(storage).DeleteUsers(ctx, id)
}

// UpdateUser
//
//	@Description[en-US]: Update user information.
//	@Description[zh-CN]: 更新用户信息.
//	@param ctx	context.Context
//	@param storage 	string
//	@param user	*models.User
//	@param updateColumns	...string
//	@return err	error
func (s Set) UpdateUser(ctx context.Context, storage string, user *models.User, updateColumns ...string) (err error) {
	return s.GetUserAndAppService(storage).UpdateUser(ctx, user, updateColumns...)
}

// GetUserInfo
//
//	@Description[en-US]: Obtain user information through ID or username.
//	@Description[zh-CN]: 通过ID或用户名获取用户信息。
//	@param ctx 	context.Context
//	@param id 	string
//	@param storage 	string
//	@param username 	string
//	@return userDetail	*models.User
//	@return err	error
func (s Set) GetUserInfo(ctx context.Context, storage string, id string, username string) (user *models.User, err error) {
	return s.GetUserAndAppService(storage).GetUserInfo(ctx, id, username)
}

// CreateUser
//
//	@Description[en-US]: Create a user.
//	@Description[zh-CN]: 创建用户。
//	@param ctx 	context.Context
//	@param storage 	string
//	@param user 	*models.User
//	@return err	error
func (s Set) CreateUser(ctx context.Context, storage string, user *models.User) (err error) {
	if len(user.Username) == 0 {
		return errors.ParameterError("username is null")
	}
	return s.GetUserAndAppService(storage).CreateUser(ctx, user)
}

// CreateUserKey
//
//	@Description[en-US]: Create a user key-pair.
//	@Description[zh-CN]: 创建用户密钥对。
//	@param ctx 	context.Context
//	@param userId 	string
//	@param name 	string
//	@return keyPair	*models.UserKey
//	@return err	error
func (s Set) CreateUserKey(ctx context.Context, userId, name string) (keyPair *models.UserKey, err error) {
	return s.commonService.CreateUserKeyWithId(ctx, userId, name)
}

// DeleteUserKey
//
//	@Description[en-US]: Delete a user key-pair.
//	@Description[zh-CN]: 删除一个用户密钥对。
//	@param ctx 	context.Context
//	@param userId 	string
//	@param id 	string
//	@return error
func (s Set) DeleteUserKey(ctx context.Context, userId string, id string) (err error) {
	_, err = s.commonService.DeleteUserKeys(ctx, userId, []string{id})
	return err
}

func (s Set) GetUserKeys(ctx context.Context, userId string, current, pageSize int64) (count int64, keyPairs []*models.UserKey, err error) {
	return s.commonService.GetUserKeys(ctx, userId, current, pageSize)
}

// PatchUser
//
//	@Description[en-US]: Incremental update user.
//	@Description[zh-CN]: 增量更新用户。
//	@param ctx 	context.Context
//	@param storage 	string
//	@param user 	map[string]interface{}
//	@return err	error
func (s Set) PatchUser(ctx context.Context, storage string, user map[string]interface{}) (err error) {
	return s.GetUserAndAppService(storage).PatchUser(ctx, user)
}

// DeleteUser
//
//	@Description[en-US]: Delete a user.
//	@Description[zh-CN]: 删除用户。
//	@param ctx 	context.Context
//	@param storage 	string
//	@param id 	string
//	@return error
func (s Set) DeleteUser(ctx context.Context, storage string, id string) (err error) {
	return s.GetUserAndAppService(storage).DeleteUser(ctx, id)
}

// VerifyPasswordById
//
//	@Description[en-US]: Verify the user's password through ID.
//	@Description[zh-CN]: 通过ID验证用户密码。
//	@param ctx 	context.Context
//	@param id 	string
//	@param password 	string
//	@return users	[]*models.User
func (s Set) VerifyPasswordById(ctx context.Context, storage, userId, password string) (users models.Users) {
	return s.GetUserAndAppService(storage).VerifyPasswordById(ctx, userId, password)
}

// VerifyPassword
//
//	@Description[en-US]: Verify password for user.
//	@Description[zh-CN]: 验证用户密码。
//	@param ctx 	context.Context
//	@param username 	string
//	@param password 	string
//	@return users	[]*models.User
func (s Set) VerifyPassword(ctx context.Context, username string, password string) (users models.Users, err error) {
	logger := logs.GetContextLogger(ctx)
	for _, userService := range s.userAndAppService {
		for _, user := range userService.VerifyPassword(ctx, username, password) {
			if user.Status != models.UserMeta_normal {
				continue
			}
			user.ExtendedData, err = s.commonService.GetUserExtendedData(ctx, user.Id)
			if err != nil {
				return nil, err
			}
			user.LoginTime = new(time.Time)
			*user.LoginTime = time.Now().UTC()
			if err = userService.UpdateLoginTime(ctx, user.Id); err != nil {
				level.Error(logger).Log("msg", "failed to update user login time", "err", err)
			}
			user.Storage = userService.Name()
			users = append(users, user)
			//if global.UserInactiveTime > 0 && time.Since(*user.LoginTime) > global.UserInactiveTime {
			//	return nil, errors.NewServerError(http.StatusForbidden, "The user is disabled due to inactivity. Please change the password and log in again.", errors.CodeUserNeedResetPassword)
			//}
		}
	}
	return users, nil
}

// Authentication
//
//	@Description[en-US]: Authenticate the user's request.
//	@Description[zh-CN]: 对用户请求进行身份认证。
//	@param ctx 	context.Context
//	@param method 	models.AuthMeta_Method
//	@param algorithm 	sign.AuthAlgorithm
//	@param key 	string
//	@param secret 	string
//	@param payload 	string
//	@param signStr 	string
//	@return ${ret_name}	[]*models.User
//	@return ${ret_name}	error
func (s Set) Authentication(ctx context.Context, method models.AuthMeta_Method, algorithm sign.AuthAlgorithm, key, secret, payload, signStr string) (models.Users, error) {
	if method == models.AuthMeta_basic {
		if _, err := uuid.FromString(key); err != nil {
			return s.VerifyPassword(ctx, key, secret)
		}
	}
	var user *models.User
	userKey, err := s.commonService.GetUserKey(ctx, key)
	if err != nil {
		return nil, err
	} else if userKey == nil {
		return nil, nil
	} else {
		for _, service := range s.userAndAppService {
			if info, err := service.GetUserInfo(ctx, userKey.UserId, ""); err != nil {
				continue
			} else {
				user = info
				break
			}
		}
	}
	if user == nil {
		return nil, errors.StatusNotFound("user")
	}
	switch method {
	case models.AuthMeta_basic:
		if userKey.Secret == secret {
			return []*models.User{user}, nil
		}
	case models.AuthMeta_signature:
		if sign.Verify(userKey.Key, userKey.Secret, userKey.Private, algorithm, signStr, payload) {
			return []*models.User{user}, nil
		}
		return nil, errors.ParameterError("Failed to verify the signature")
	default:
		return nil, errors.ParameterError("unknown auth method")
	}
	return nil, errors.ParameterError("unknown auth request")
}

// GetAuthCodeByAppId
//
//	@Description[en-US]: Get auth code by client id.
//	@Description[zh-CN]: 通过客户端id获取授权代。
//	@param ctx 	context.Context
//	@param clientId 	string
//	@param user 	*models.User
//	@param sessionId 	string
//	@param storage 	string
//	@return code	string
//	@return err	error
func (s Set) GetAuthCodeByAppId(ctx context.Context, clientId string, user *models.User, sessionId, storage string) (code string, err error) {
	svc := s.GetUserAndAppService(storage)
	if svc == nil {
		return "", errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
	}
	app := models.App{Model: models.Model{Id: clientId}}
	err = s.GetAppAccessControl(ctx, &app, opts.WithoutUsers, opts.WithUsers(user.Id))
	if err != nil {
		return "", err
	}
	if len(app.Roles) == 0 || len(app.Users) == 0 {
		return "", errors.NotFoundError()
	}
	token, err := s.CreateToken(ctx, models.TokenTypeCode, user)
	if err != nil {
		return "", err
	}
	return token.Id, nil
}

func (s Set) GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (users models.Users) {
	for _, service := range s.userAndAppService {
		if info, err := service.GetUserInfoByUsernameAndEmail(ctx, username, email); err == nil {
			info.Storage = service.Name()
			users = append(users, info)
		} else {
			level.Error(logs.GetContextLogger(ctx)).Log("err", err, "msg", "Failed to get user from username and email")
		}
	}
	return users
}
