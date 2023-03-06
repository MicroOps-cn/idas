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
	"github.com/MicroOps-cn/idas/pkg/service/opts"
)

func newNullService(serviceType string, serviceName string) *nullService {
	return &nullService{n: serviceName, t: serviceType}
}

type nullService struct {
	t, n string
}

func (n nullService) GetUsersById(_ context.Context, _ []string) (models.Users, error) {
	return nil, n
}

func (n nullService) VerifyPasswordById(_ context.Context, _, _ string) (users []*models.User) {
	return nil
}

func (n nullService) Error() string {
	if n.t == "" {
		return fmt.Sprintf("service not foundL %s", n.n)
	}
	return fmt.Sprintf("%s service not foundL %s", n.t, n.n)
}

func (n nullService) AutoMigrate(_ context.Context) error {
	return n
}

func (n nullService) Name() string {
	return n.n
}

func (n nullService) GetUsers(_ context.Context, _ string, _ models.UserMeta_UserStatus, _ string, _, _ int64) (total int64, users []*models.User, err error) {
	return 0, nil, n
}

func (n nullService) PatchUsers(_ context.Context, _ []map[string]interface{}) (count int64, err error) {
	return 0, n
}

func (n nullService) DeleteUsers(_ context.Context, _ []string) (count int64, err error) {
	return 0, n
}

func (n nullService) UpdateUser(_ context.Context, _ *models.User, _ ...string) error {
	return n
}

func (n nullService) UpdateLoginTime(_ context.Context, _ string) error {
	return n
}

func (n nullService) GetUserInfo(_ context.Context, _ string, _ string) (*models.User, error) {
	return nil, n
}

func (n nullService) CreateUser(_ context.Context, _ *models.User) (err error) {
	return n
}

func (n nullService) PatchUser(_ context.Context, _ map[string]interface{}) (err error) {
	return n
}

func (n nullService) DeleteUser(_ context.Context, _ string) error {
	return n
}

func (n nullService) VerifyPassword(_ context.Context, _ string, _ string) []*models.User {
	return nil
}

func (n nullService) GetApps(_ context.Context, _ string, _ map[string]interface{}, _, _ int64) (total int64, apps []*models.App, err error) {
	return 0, nil, n
}

func (n nullService) PatchApps(_ context.Context, _ []map[string]interface{}) (total int64, err error) {
	return 0, n
}

func (n nullService) DeleteApps(_ context.Context, _ []string) (total int64, err error) {
	return 0, n
}

func (n nullService) UpdateApp(_ context.Context, _ *models.App, _ ...string) (err error) {
	return n
}

func (n nullService) GetAppInfo(_ context.Context, _ ...opts.WithGetAppOptions) (app *models.App, err error) {
	return nil, n
}

func (n nullService) CreateApp(_ context.Context, _ *models.App) (err error) {
	return n
}

func (n nullService) PatchApp(_ context.Context, _ map[string]interface{}) (err error) {
	return n
}

func (n nullService) DeleteApp(_ context.Context, _ string) (err error) {
	return n
}

func (n nullService) ResetPassword(_ context.Context, _ string, _ string) error {
	return n
}

func (n nullService) GetUserInfoByUsernameAndEmail(_ context.Context, _, _ string) (*models.User, error) {
	return nil, n
}
