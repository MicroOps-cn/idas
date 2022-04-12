package service

import (
	"context"
	"fmt"
	"github.com/go-kit/log/level"
	"idas/pkg/errors"
	"idas/pkg/logs"
	"idas/pkg/service/models"
	"idas/pkg/utils/image"
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

func (s Set) GetApps(ctx context.Context, storage string, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error) {
	return s.SafeGetUserAndAppService(storage).GetApps(ctx, keywords, current, pageSize)
}

func (s Set) GetAppSource(ctx context.Context) (data map[string]string, total int64, err error) {
	data = map[string]string{}
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
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.GetAppInfo(ctx, id, "")
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
