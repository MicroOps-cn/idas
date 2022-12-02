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
	"sort"
	"strings"

	gogorm "gorm.io/gorm"

	"github.com/MicroOps-cn/idas/pkg/client/gorm"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/sign"
)

func NewCommonService(name string, client *gorm.Client) *CommonService {
	return &CommonService{name: name, Client: client}
}

type CommonService struct {
	*gorm.Client
	name string
}

func (c CommonService) Name() string {
	return c.name
}

func (c CommonService) AutoMigrate(ctx context.Context) error {
	err := c.Session(ctx).AutoMigrate(
		&models.File{},
		&models.Permission{},
		&models.Role{},
		&models.UserKey{},
		&models.AppProxy{},
		&models.AppProxyUrl{},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c CommonService) RecordUploadFile(ctx context.Context, name string, path string, contentType string, size int64) (id string, err error) {
	file := &models.File{MimiType: contentType, Name: name, Path: path, Size: size}
	if err = c.Session(ctx).Create(file).Error; err != nil {
		return
	}
	return file.Id, err
}

func (c CommonService) GetFileInfoFromId(ctx context.Context, id string) (fileName, mimiType, filePath string, err error) {
	file := &models.File{Model: models.Model{Id: id}}
	if err = c.Session(ctx).First(file).Error; err != nil {
		return "", "", "", err
	}
	return file.Name, file.MimiType, file.Path, nil
}

func (c CommonService) CreateUserKeyWithId(ctx context.Context, userId string, name string) (userKey *models.UserKey, err error) {
	conn := c.Session(ctx)
	if err != nil {
		return nil, err
		//} else if id := strings.ReplaceAll(uuid.NewV4().String(), "-", ""); len(id) != 32 {
		//	return nil, errors.New("生成ID失败: " + id)
	} else if pub1, pub2, privateKey, err := sign.GenerateECDSAKeyPair(); err != nil {
		return nil, err
	} else {
		userKey = &models.UserKey{
			Name:   name,
			UserId: userId,
			Key:    pub1,
			Secret: pub2,
		}
		if err = conn.Create(&userKey).Error; err != nil {
			return nil, err
		}
		return &models.UserKey{
			Model:   userKey.Model,
			UserId:  userId,
			Key:     pub1,
			Secret:  pub2,
			Private: privateKey,
		}, nil
	}
}

func (c CommonService) GetUserKey(ctx context.Context, key string) (*models.UserKey, error) {
	userKey := &models.UserKey{Key: key}
	conn := c.Session(ctx)
	if err := conn.Where("`key` = ?", key).First(&userKey).Error; err != nil && err != gogorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to query user key, database exception: %s", err)
	} else if err != nil {
		return nil, nil
	}
	return userKey, nil
}

func (c CommonService) GetUserKeys(ctx context.Context, userId string, current, pageSize int64) (count int64, keyPairs []*models.UserKey, err error) {
	query := c.Session(ctx).Model(&models.UserKey{}).Where("user_id = ? and is_delete = 0", userId)
	if err = query.Select("id", "name", "create_time", "key").Order("id").Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).
		Find(&keyPairs).Error; err != nil {
		return 0, nil, err
	} else if err = query.Count(&count).Error; err != nil {
		return 0, nil, err
	} else {
		for _, keyPair := range keyPairs {
			keyPair.UserId = userId
		}
		return count, keyPairs, nil
	}
}

func (c CommonService) DeleteUserKeys(ctx context.Context, userId string, id []string) (affected int64, err error) {
	deleted := c.Session(ctx).Model(&models.UserKey{}).Where("id in ? and user_id = ?", id, userId).Update("is_delete", true)
	if err = deleted.Error; err != nil {
		return deleted.RowsAffected, err
	}
	return deleted.RowsAffected, nil
}

func (c CommonService) GetProxyConfig(ctx context.Context, host string, method string, path string) (*models.AppProxyConfig, error) {
	conn := c.Session(ctx)
	var proxy []models.AppProxy

	if err := conn.Find(&proxy, "host = ?", host).Error; err != nil {
		return nil, err
	}
	var dftURL *models.AppProxyConfig
	for _, appProxy := range proxy {
		sort.Sort(appProxy.Urls)
		for _, url := range appProxy.Urls {
			if url.Method == "*" && url.Url == "/" {
				dftURL = &models.AppProxyConfig{
					AppProxyUrl: url,
					Domain:      appProxy.Domain,
					Upstream:    appProxy.Domain,
				}
			} else if url.Method == "*" || url.Method == method {
				if strings.HasPrefix(path, url.Url) {
					return &models.AppProxyConfig{
						AppProxyUrl: url,
						Domain:      appProxy.Domain,
						Upstream:    appProxy.Domain,
					}, nil
				}
			}
		}
	}
	if dftURL != nil {
		return dftURL, nil
	}
	return nil, errors.NotFoundError
}

func (c CommonService) UpdateProxyConfig(ctx context.Context, proxy *models.AppProxy) (newProxy *models.AppProxy, err error) {
	conn := c.Session(ctx)
	var model models.Model
	for i, url := range proxy.Urls {
		url.Index = uint32(i)
	}
	if err = conn.Model(&models.AppProxy{}).Select("id").Where("app_id = ?", proxy.AppId).First(&model).Error; err != nil && err != gogorm.ErrRecordNotFound {
		return nil, err
	}
	proxy.Id = model.Id
	if len(proxy.Id) > 0 {
		if err = conn.Updates(proxy).Error; err != nil {
			return nil, err
		}
	} else {
		if err = conn.Create(proxy).Error; err != nil {
			return nil, err
		}
	}
	newProxy = new(models.AppProxy)
	if err = conn.Where("id = ?", proxy.Id).First(&newProxy).Error; err != nil {
		return nil, err
	}
	return newProxy, nil
}

func (c CommonService) GetAppProxyConfig(ctx context.Context, appId string) (proxy *models.AppProxy, err error) {
	if err = c.Session(ctx).Where("app_id = ?", appId).Preload("Urls").First(&proxy).Error; err != nil {
		return nil, err
	}
	return proxy, nil
}
