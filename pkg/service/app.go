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

	"github.com/go-kit/log/level"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/logs"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/image"
)

func (s Set) GetUserAndAppService(name string) UserAndAppService {
	for _, svc := range s.userAndAppService {
		if svc.Name() == name /*|| len(name) == 0 */ {
			return svc
		}
	}
	return nil
}

func (s Set) SafeGetUserAndAppService(name string) UserAndAppService {
	for _, svc := range s.userAndAppService {
		if svc.Name() == name {
			return svc
		}
	}
	for _, svc := range s.userAndAppService {
		return svc
	}
	return nil
}

func (s Set) GetApps(ctx context.Context, storage string, keywords string, current, pageSize int64) (total int64, apps []*models.App, err error) {
	return s.SafeGetUserAndAppService(storage).GetApps(ctx, keywords, current, pageSize)
}

func (s Set) GetAppSource(ctx context.Context) (total int64, data map[string]string, err error) {
	total = int64(len(s.userAndAppService))
	data = make(map[string]string, total)
	for _, appService := range s.userAndAppService {
		data[appService.Name()] = appService.Name()
	}

	return
}

func (s Set) PatchApps(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.PatchApps(ctx, patch)
}

func (s Set) DeleteApps(ctx context.Context, storage string, id []string) (total int64, err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.DeleteApps(ctx, id)
}

func (s Set) UpdateApp(ctx context.Context, storage string, app *models.App, updateColumns ...string) (a *models.App, err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.UpdateApp(ctx, app, updateColumns...)
}

func (s Set) GetAppInfo(ctx context.Context, storage string, id string) (app *models.App, err error) {
	if len(storage) == 0 {
		for _, service := range s.userAndAppService {
			info, err := service.GetAppInfo(ctx, id, "")
			if err != nil {
				continue
			}
			return info, nil
		}
	} else {
		service := s.GetUserAndAppService(storage)
		if service == nil {
			err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
			return
		}
		return service.GetAppInfo(ctx, id, "")
	}
	return nil, errors.StatusNotFound("app")
}

func (s Set) CreateApp(ctx context.Context, storage string, app *models.App) (a *models.App, err error) {
	logger := logs.GetContextLogger(ctx)
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	if avatar, err := image.GenerateAvatar(ctx, app.Name); err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to generate avatar")
	} else if fileKey, err := s.UploadFile(ctx, app.Name+".png", "image/png", avatar); err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to save avatar")
	} else {
		app.Avatar = fileKey
	}
	return service.CreateApp(ctx, app)
}

func (s Set) PatchApp(ctx context.Context, storage string, fields map[string]interface{}) (app *models.App, err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.PatchApp(ctx, fields)
}

func (s Set) DeleteApp(ctx context.Context, storage string, id string) (err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.DeleteApp(ctx, id)
}
