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
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

func newNullService(serviceType string, serviceName string) *nullService {
	return &nullService{n: serviceName, t: serviceType}
}

type nullService struct {
	t, n string
}

func (n nullService) VerifyPasswordById(ctx context.Context, id, password string) (users []*models.User) {
	return nil
}

func (n nullService) Error() string {
	if n.t == "" {
		return fmt.Sprintf("service not foundL %s", n.n)
	}
	return fmt.Sprintf("%s service not foundL %s", n.t, n.n)
}

func (n nullService) AutoMigrate(ctx context.Context) error {
	return n
}

func (n nullService) Name() string {
	return n.n
}

func (n nullService) GetUsers(ctx context.Context, keywords string, status models.UserMeta_UserStatus, appId string, current, pageSize int64) (total int64, users []*models.User, err error) {
	return 0, nil, n
}

func (n nullService) PatchUsers(ctx context.Context, patch []map[string]interface{}) (count int64, err error) {
	return 0, n
}

func (n nullService) DeleteUsers(ctx context.Context, id []string) (count int64, err error) {
	return 0, n
}

func (n nullService) UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (*models.User, error) {
	return nil, n
}

func (n nullService) UpdateLoginTime(ctx context.Context, id string) error {
	return n
}

func (n nullService) GetUserInfo(ctx context.Context, id string, username string) (*models.User, error) {
	return nil, n
}

func (n nullService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, n
}

func (n nullService) PatchUser(ctx context.Context, user map[string]interface{}) (*models.User, error) {
	return nil, n
}

func (n nullService) DeleteUser(ctx context.Context, id string) error {
	return n
}

func (n nullService) VerifyPassword(ctx context.Context, username string, password string) []*models.User {
	return nil
}

func (n nullService) GetApps(ctx context.Context, keywords string, current, pageSize int64) (total int64, apps []*models.App, err error) {
	return 0, nil, n
}

func (n nullService) PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error) {
	return 0, n
}

func (n nullService) DeleteApps(ctx context.Context, id []string) (total int64, err error) {
	return 0, n
}

func (n nullService) UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (*models.App, error) {
	return nil, n
}

func (n nullService) GetAppInfo(ctx context.Context, id string, name string) (app *models.App, err error) {
	return nil, n
}

func (n nullService) CreateApp(ctx context.Context, app *models.App) (*models.App, error) {
	return nil, n
}

func (n nullService) PatchApp(ctx context.Context, fields map[string]interface{}) (app *models.App, err error) {
	return nil, n
}

func (n nullService) DeleteApp(ctx context.Context, id string) (err error) {
	return n
}

func (n nullService) VerifyUserAuthorizationForApp(ctx context.Context, appId string, userId string) (role string, err error) {
	return "", n
}

func (n nullService) ResetPassword(ctx context.Context, id string, password string) error {
	return n
}

func (n nullService) GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (*models.User, error) {
	return nil, n
}
