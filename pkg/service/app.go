package service

import (
	"context"
	"fmt"
	"github.com/go-kit/log/level"
	"idas/pkg/client/gorm"
	"idas/pkg/client/ldap"
	"idas/pkg/service/gormservice"
	"idas/pkg/service/ldapservice"

	"idas/config"
	"idas/pkg/errors"
	"idas/pkg/logs"
	"idas/pkg/service/models"
	"idas/pkg/utils/image"
)

type AppService interface {
	baseService
	Name() string
	GetApps(ctx context.Context, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error)
	PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error)
	DeleteApps(ctx context.Context, id []string) (total int64, err error)
	UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (*models.App, error)
	GetAppInfo(ctx context.Context, id string, name string) (app *models.App, err error)
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
	if len(config.Get().GetStorage().GetUser()) > 0 {
		for _, appStorage := range config.Get().GetStorage().GetUser() {
			if appServices.Include(appStorage.GetName()) {
				panic(any(fmt.Errorf("Failed to init AppService: duplicate datasource: %T ", appStorage.Name)))
			}
			switch appSource := appStorage.GetStorageSource().(type) {
			case *config.Storage_Mysql:
				if client, err := gorm.NewMySQLClient(ctx, appSource.Mysql); err != nil {
					panic(any(fmt.Errorf("初始化AppService失败: MySQL数据库连接失败: %s", err)))
				} else {
					appServices = append(appServices, gormservice.NewAppService(appStorage.GetName(), client))
				}
			case *config.Storage_Sqlite:
				if client, err := gorm.NewSQLiteClient(ctx, appSource.Sqlite); err != nil {
					panic(any(fmt.Errorf("初始化AppService失败: MySQL数据库连接失败: %s", err)))
				} else {
					appServices = append(appServices, gormservice.NewAppService(appStorage.GetName(), client))
				}
			case *config.Storage_Ldap:
				if client, err := ldap.NewLdapClient(ctx, appSource.Ldap); err != nil {
					panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
				} else {
					appServices = append(appServices, ldapservice.NewAppService(appStorage.GetName(), client))
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
	for _, appService := range s.appService {
		data[appService.Name()] = appService.Name()
	}
	return
}

func (s Set) PatchApps(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.PatchApps(ctx, patch)
}

func (s Set) DeleteApps(ctx context.Context, storage string, id []string) (total int64, err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.DeleteApps(ctx, id)
}

func (s Set) UpdateApp(ctx context.Context, storage string, app *models.App, updateColumns ...string) (a *models.App, err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.UpdateApp(ctx, app, updateColumns...)
}

func (s Set) GetAppInfo(ctx context.Context, storage string, id string) (app *models.App, err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.GetAppInfo(ctx, id, "")
}

func (s Set) CreateApp(ctx context.Context, storage string, app *models.App) (a *models.App, err error) {
	logger := logs.GetContextLogger(ctx)
	service := s.GetAppService(storage)
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
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.PatchApp(ctx, fields)
}

func (s Set) DeleteApp(ctx context.Context, storage string, id string) (err error) {
	service := s.GetAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.DeleteApp(ctx, id)
}
