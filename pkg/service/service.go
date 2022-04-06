package service

import (
	"context"

	"idas/pkg/service/models"
)

type migrator interface {
	AutoMigrate(ctx context.Context) error
}

type baseService interface {
	migrator
	Name() string
}

type Service interface {
	baseService
	SessionService
	GetUsers(ctx context.Context, storage string, keyword string, status models.UserStatus, current int64, pageSize int64) (users []*models.User, total int64, err error)
	PatchUsers(ctx context.Context, storage string, patch []map[string]interface{}) (count int64, err error)
	DeleteUsers(ctx context.Context, storage string, id []string) (count int64, err error)
	UpdateUser(ctx context.Context, storage string, user *models.User, updateColumns ...string) (*models.User, error)
	GetUserInfo(ctx context.Context, storage string, id string, username string) (user *models.User, err error)
	CreateUser(ctx context.Context, storage string, user *models.User) (*models.User, error)
	PatchUser(ctx context.Context, storage string, user map[string]interface{}) (*models.User, error)
	DeleteUser(ctx context.Context, storage string, id string) error
	CreateLoginSession(ctx context.Context, username string, password string) (string, error)
	GetUserSource(ctx context.Context) (data map[string]string, total int64, err error)

	GetApps(ctx context.Context, storage string, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error)
	GetAppSource(ctx context.Context) (data map[string]string, total int64, err error)
	PatchApps(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, err error)
	DeleteApps(ctx context.Context, storage string, id []string) (total int64, err error)
	UpdateApp(ctx context.Context, storage string, app *models.App, updateColumns ...string) (*models.App, error)
	GetAppInfo(ctx context.Context, storage string, id string) (app *models.App, err error)
	CreateApp(ctx context.Context, storage string, app *models.App) (*models.App, error)
	PatchApp(ctx context.Context, storage string, fields map[string]interface{}) (app *models.App, err error)
	DeleteApp(ctx context.Context, storage string, id string) (err error)
}

type Set struct {
	userService UserServices
	appService  AppServices
	SessionService
}

func (s Set) AutoMigrate(ctx context.Context) error {
	for _, svc := range s.userService {
		if err := svc.AutoMigrate(ctx); err != nil {
			return err
		}
	}
	for _, svc := range s.appService {
		if err := svc.AutoMigrate(ctx); err != nil {
			return err
		}
	}
	return s.SessionService.AutoMigrate(ctx)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(ctx context.Context) Service {
	return &Set{
		userService:    NewUserServices(ctx),
		SessionService: NewSessionService(ctx),
		appService:     NewAppService(ctx),
	}
}
