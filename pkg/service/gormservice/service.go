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

	"github.com/MicroOps-cn/fuck/clients/gorm"

	"github.com/MicroOps-cn/idas/pkg/service/models"
)

func NewUserAndAppService(_ context.Context, name string, client *gorm.Client) *UserAndAppService {
	set := &UserAndAppService{name: name, Client: client}
	return set
}

type UserAndAppService struct {
	*gorm.Client
	name string
}

func (s UserAndAppService) AutoMigrate(ctx context.Context) error {
	err := s.Session(ctx).AutoMigrate(&models.App{}, &models.User{})
	if err != nil {
		return err
	}

	return nil
}
