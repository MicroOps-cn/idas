package service

import (
	"context"
	"fmt"
	"idas/config"
	"idas/pkg/client/mysql"
	"idas/pkg/errors"
	"idas/pkg/service/models"
	"idas/pkg/service/mysqlservice"
)

type AppService interface {
	baseService
	Name() string
	GetApps(ctx context.Context, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error)
	PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error)
	DeleteApps(ctx context.Context, id []string) (total int64, err error)
	UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (*models.App, error)
	GetAppInfo(ctx context.Context, id string) (app *models.App, err error)
	CreateApp(ctx context.Context, app *models.App) (*models.App, error)
	PatchApp(ctx context.Context, fields map[string]interface{}) (app *models.App, err error)
	DeleteApp(ctx context.Context, id string) (err error)
}

type AppServices []AppService

func (s AppServices) Include(name string) bool {
	for _, service := range s {
		if service.Name() == name {
			return true
		}
	}
	return false
}
func NewAppService(ctx context.Context) AppServices {
	var appServices AppServices
	if len(config.Get().GetStorage().GetApp()) > 0 {
		for _, appStorage := range config.Get().GetStorage().GetApp() {

			if appServices.Include(appStorage.GetName()) {
				panic(any(fmt.Errorf("Failed to init AppService: duplicate datasource: %T ", appStorage.Name)))
			}
			switch appSource := appStorage.GetStorageSource().(type) {
			case *config.Storage_Mysql:
				if client, err := mysql.NewMySQLClient(ctx, appSource.Mysql); err != nil {
					panic(any(fmt.Errorf("初始化AppService失败: MySQL数据库连接失败: %s", err)))
				} else {
					appServices = append(appServices, mysqlservice.NewAppService(appStorage.GetName(), client))
				}
			default:
				panic(any(fmt.Errorf("Failed to init AppService: Unknown datasource: %T ", appSource)))
			}
		}
	}
	return appServices
}

func (s Set) GetAppService(name string) AppService {
	for _, appService := range s.appService {
		if appService.Name() == name /*|| len(name) == 0 */ {
			return appService
		}
	}
	return nil
}

func (s Set) SafeGetAppService(name string) AppService {
	for _, appService := range s.appService {
		if appService.Name() == name {
			return appService
		}
	}
	for _, service := range s.appService {
		return service
	}
	return nil
}
func (s Set) GetApps(ctx context.Context, storage string, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error) {
	return s.SafeGetAppService(storage).GetApps(ctx, keywords, current, pageSize)
}

func (s Set) GetAppSource(ctx context.Context) (data map[string]string, total int64, err error) {
	data = map[string]string{}
	for _, userService := range s.userService {
		data[userService.Name()] = userService.Name()
	}
	return
}

func (s Set) PatchApps(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.PatchApps(ctx, patch)
	}
}

func (s Set) DeleteApps(ctx context.Context, storage string, id []string) (total int64, err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.DeleteApps(ctx, id)
	}
}

func (s Set) UpdateApp(ctx context.Context, storage string, app *models.App, updateColumns ...string) (a *models.App, err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.UpdateApp(ctx, app, updateColumns...)
	}
}

func (s Set) GetAppInfo(ctx context.Context, storage string, id string) (app *models.App, err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.GetAppInfo(ctx, id)
	}
}

func (s Set) CreateApp(ctx context.Context, storage string, app *models.App) (a *models.App, err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.CreateApp(ctx, app)
	}
}

func (s Set) PatchApp(ctx context.Context, storage string, fields map[string]interface{}) (app *models.App, err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.PatchApp(ctx, fields)
	}
}

func (s Set) DeleteApp(ctx context.Context, storage string, id string) (err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.DeleteApp(ctx, id)
	}
}
