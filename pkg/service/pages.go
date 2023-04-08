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

package service

import (
	"context"
	"encoding/json"

	"github.com/MicroOps-cn/idas/pkg/service/models"
)

func (s Set) DeletePages(ctx context.Context, ids []string) error {
	return s.commonService.DeletePages(ctx, ids)
}

func (s Set) UpdatePage(ctx context.Context, page *models.PageConfig) error {
	return s.commonService.UpdatePage(ctx, page)
}

func (s Set) CreatePage(ctx context.Context, page *models.PageConfig) error {
	return s.commonService.CreatePage(ctx, page)
}

func (s Set) GetPages(ctx context.Context, filter map[string]interface{}, keywords string, current int64, size int64) (int64, []*models.PageConfig, error) {
	return s.commonService.GetPages(ctx, filter, keywords, current, size)
}

func (s Set) GetPage(ctx context.Context, id string) (*models.PageConfig, error) {
	return s.commonService.GetPage(ctx, id)
}

func (s Set) PatchPages(ctx context.Context, patch []map[string]interface{}) error {
	return s.commonService.PatchPages(ctx, patch)
}

func (s Set) PatchPageDatas(ctx context.Context, patch []models.PageData) error {
	return s.commonService.PatchPageDatas(ctx, patch)
}

func (s Set) UpdatePageData(ctx context.Context, pageId string, id string, data *json.RawMessage) error {
	return s.commonService.UpdatePageData(ctx, pageId, id, data)
}

func (s Set) CreatePageData(ctx context.Context, pageId string, data *json.RawMessage) error {
	return s.commonService.CreatePageData(ctx, pageId, data)
}

func (s Set) GetPageData(ctx context.Context, pageId string, id string) (*models.PageData, error) {
	return s.commonService.GetPageData(ctx, pageId, id)
}

func (s Set) GetPageDatas(ctx context.Context, filters map[string]string, keywords string, current int64, size int64) (int64, []*models.PageData, error) {
	return s.commonService.GetPageDatas(ctx, filters, keywords, current, size)
}
