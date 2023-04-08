/*
 Copyright © 2023 MicroOps-cn.

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
	"encoding/json"
	"time"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

func (c CommonService) GetPages(ctx context.Context, filter map[string]interface{}, keywords string, current int64, pageSize int64) (count int64, pages []*models.PageConfig, err error) {
	conn := c.Session(ctx)
	tb := conn.Model(&models.PageConfig{}).Omit("fields")
	if keywords != "" {
		tb = tb.Where("Name LIKE ?", "%"+keywords+"%")
	}
	if len(filter) > 0 {
		tb = tb.Where(filter)
	}
	if err = tb.Count(&count).Error; err != nil {
		return 0, nil, err
	} else if count == 0 {
		return 0, nil, nil
	}

	if err = tb.Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&pages).Error; err != nil {
		return 0, nil, err
	}
	return count, pages, err
}

func (c *CommonService) GetPage(ctx context.Context, id string) (cfg *models.PageConfig, err error) {
	cfg = new(models.PageConfig)
	err = c.Session(ctx).Where("id = ?", id).First(cfg).Error

	return cfg, err
}

func (c CommonService) CreatePage(ctx context.Context, page *models.PageConfig) error {
	return c.Session(ctx).Create(page).Error
}

func (c CommonService) UpdatePage(ctx context.Context, page *models.PageConfig) error {
	return c.Session(ctx).Select("name", "description", "fields", "icon").Updates(page).Error
}

func (c CommonService) DeletePages(ctx context.Context, ids []string) error {
	return c.Session(ctx).Model(models.PageConfig{}).Where("id in (?)", ids).Update("delete_time", time.Now()).Error
}

func (c *CommonService) PatchPages(ctx context.Context, patch []map[string]interface{}) error {
	tx := c.Session(ctx).Begin()
	defer tx.Rollback()
	for _, m := range patch {
		id, ok := m["id"].(string)
		if !ok || len(id) == 0 {
			return errors.ParameterError("id")
		}
		delete(m, "id")
		if err := tx.Model(&models.PageConfig{}).Where("id = ?", id).Updates(m).Error; err != nil {
			return err
		}
	}
	return tx.Commit().Error
}

func (c *CommonService) GetPageDatas(ctx context.Context, filters map[string]string, keywords string, current int64, pageSize int64) (count int64, datas []*models.PageData, err error) {
	conn := c.Session(ctx)
	tb := conn.Model(&models.PageData{}).Omit("fields")
	if keywords != "" {
		tb = tb.Where("data LIKE ?", "%"+keywords+"%")
	}
	if len(filters) > 0 {
		tb = tb.Where(filters)
	}
	if err = tb.Count(&count).Error; err != nil {
		return 0, nil, err
	} else if count == 0 {
		return 0, nil, nil
	}
	if err = tb.Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&datas).Error; err != nil {
		return 0, nil, err
	}
	return count, datas, err
}

func (c *CommonService) GetPageData(ctx context.Context, pageId string, id string) (*models.PageData, error) {
	var data models.PageData
	if err := c.Session(ctx).Where("id = ? and page_id = ?", id, pageId).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *CommonService) CreatePageData(ctx context.Context, pageId string, data *json.RawMessage) error {
	return c.Session(ctx).Create(&models.PageData{
		PageId: pageId,
		Data:   (*models.JSON)(data),
	}).Error
}

func (c *CommonService) UpdatePageData(ctx context.Context, pageId string, id string, data *json.RawMessage) error {
	return c.Session(ctx).Model(&models.PageData{}).Where("id = ? and page_id = ?", id, pageId).Update("data", data).Error
}

func (c *CommonService) PatchPageDatas(ctx context.Context, patch []models.PageData) error {
	tx := c.Session(ctx).Begin()
	defer tx.Rollback()
	for _, data := range patch {
		if err := tx.Where("id = ? and page_id = ?", data.Id, data.PageId).Updates(&data).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}
