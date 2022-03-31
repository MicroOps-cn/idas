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
	PatchUsers(ctx context.Context, storage string, patch []map[string]interface{}) (count int64, msg string, err error)
	DeleteUsers(ctx context.Context, storage string, id []string) (count int64, msg string, err error)
	UpdateUser(ctx context.Context, storage string, user *models.User, updateColumns ...string) (*models.User, string, error)
	GetUserInfo(ctx context.Context, storage string, id string, username string) (*models.User, string, error)
	CreateUser(ctx context.Context, storage string, user *models.User) (*models.User, string, error)
	PatchUser(ctx context.Context, storage string, user map[string]interface{}) (*models.User, string, error)
	DeleteUser(ctx context.Context, storage string, id string) (string, error)
	CreateLoginSession(ctx context.Context, username string, password string) (string, error)
	GetUserSource(ctx context.Context) (data map[string]string, total int64, err error)
}

type Set struct {
	userService UserServices
	SessionService
}

func (s Set) AutoMigrate(ctx context.Context) error {
	for _, svc := range s.userService {
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
	}
}
