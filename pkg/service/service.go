package service

import (
	"context"
	"fmt"
	"time"

	"idas/config"
	"idas/pkg/client/mysql"
	"idas/pkg/client/redis"
	"idas/pkg/errors"
	"idas/pkg/service/models"
	"idas/pkg/service/mysqlservice"
	"idas/pkg/service/redisservice"
)

type migrator interface {
	AutoMigrate(ctx context.Context) error
}

type Service interface {
	migrator
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

func (s Set) GetUserService(name string) UserService {
	for n, c := range s.userService {
		if n == name /*|| len(name) == 0 */ {
			return c
		}
	}
	return nil
}

func (s Set) UserServiceDo(name string, f func(service UserService)) error {
	service := s.GetUserService(name)
	if service == nil {
		return errors.StatusNotFound("User")
	}
	f(service)
	return nil
}

func (s Set) SafeGetUserService(name string) UserService {
	for n, c := range s.userService {
		if n == name {
			return c
		}
	}
	for _, service := range s.userService {
		return service
	}
	return nil
}

type Set struct {
	userService map[string]UserService
	SessionService
}

func (s Set) GetUserSource(ctx context.Context) (data map[string]string, total int64, err error) {
	data = map[string]string{}
	for name := range s.userService {
		data[name] = name
	}
	return
}

func (s Set) GetUsers(ctx context.Context, storage string, keyword string, status models.UserStatus, current int64, pageSize int64) (users []*models.User, total int64, err error) {
	return s.SafeGetUserService(storage).GetUsers(ctx, keyword, status, current, pageSize)
}

func (s Set) PatchUsers(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, msg string, err error) {
	err = s.UserServiceDo(storage, func(service UserService) {
		total, msg, err = service.PatchUsers(ctx, patch)
	})
	return
}

func (s Set) DeleteUsers(ctx context.Context, storage string, id []string) (total int64, msg string, err error) {
	err = s.UserServiceDo(storage, func(service UserService) {
		total, msg, err = service.DeleteUsers(ctx, id)
	})
	return
}

func (s Set) UpdateUser(ctx context.Context, storage string, user *models.User, updateColumns ...string) (u *models.User, msg string, err error) {
	err = s.UserServiceDo(storage, func(service UserService) {
		u, msg, err = service.UpdateUser(ctx, user, updateColumns...)
	})
	return
}

func (s Set) GetUserInfo(ctx context.Context, storage string, id string, username string) (user *models.User, msg string, err error) {
	err = s.UserServiceDo(storage, func(service UserService) {
		user, msg, err = service.GetUserInfo(ctx, id, username)
	})
	return
}

func (s Set) CreateUser(ctx context.Context, storage string, user *models.User) (u *models.User, msg string, err error) {
	err = s.UserServiceDo(storage, func(service UserService) {
		u, msg, err = service.CreateUser(ctx, user)
	})
	return
}

func (s Set) PatchUser(ctx context.Context, storage string, user map[string]interface{}) (u *models.User, msg string, err error) {
	err = s.UserServiceDo(storage, func(service UserService) {
		u, msg, err = service.PatchUser(ctx, user)
	})
	return
}

func (s Set) DeleteUser(ctx context.Context, storage string, id string) (msg string, err error) {
	err = s.UserServiceDo(storage, func(service UserService) {
		msg, err = service.DeleteUser(ctx, id)
	})
	return
}

func (s Set) VerifyPassword(ctx context.Context, username string, password string) (user *models.User, err error) {
	for storageName, userService := range s.userService {
		user, err = userService.VerifyPassword(ctx, username, password)
		if err == nil {
			user.Storage = storageName
			return user, nil
		}
	}
	return nil, errors.UnauthorizedError
}

func (s Set) AutoMigrate(ctx context.Context) error {
	for _, svc := range s.userService {
		if err := svc.AutoMigrate(ctx); err != nil {
			return err
		}
	}
	return s.SessionService.AutoMigrate(ctx)
}

func (s Set) CreateLoginSession(ctx context.Context, username string, password string) (session string, err error) {
	user, err := s.VerifyPassword(ctx, username, password)
	if user == nil {
		return "", err
	}
	user.LoginTime = time.Now().UTC()
	if user, _, err = s.GetUserService(user.Storage).UpdateUser(ctx, user, "login_time"); err != nil {
		return "", err
	}
	return s.SessionService.SetLoginSession(ctx, user)
}

type SessionService interface {
	migrator
	SetLoginSession(ctx context.Context, user *models.User) (string, error)
	DeleteLoginSession(ctx context.Context, session string) (string, error)
	GetLoginSession(ctx context.Context, id string) (*models.User, string, error)
	OAuthAuthorize(ctx context.Context, responseType, clientId, redirectURI string) (redirect string, err error)
	GetOAuthTokenByAuthorizationCode(ctx context.Context, code, clientId, redirectURI string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error)
	GetOAuthTokenByPassword(ctx context.Context, username string, password string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error)
}

type UserService interface {
	migrator
	GetUsers(ctx context.Context, keyword string, status models.UserStatus, current int64, pageSize int64) (users []*models.User, total int64, err error)
	PatchUsers(ctx context.Context, patch []map[string]interface{}) (count int64, msg string, err error)
	DeleteUsers(ctx context.Context, id []string) (count int64, msg string, err error)
	UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (*models.User, string, error)
	GetUserInfo(ctx context.Context, id string, username string) (*models.User, string, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, string, error)
	PatchUser(ctx context.Context, user map[string]interface{}) (*models.User, string, error)
	DeleteUser(ctx context.Context, id string) (string, error)
	VerifyPassword(ctx context.Context, username string, password string) (*models.User, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(ctx context.Context) Service {
	userService := map[string]UserService{}
	var sessionService SessionService
	if len(config.Get().GetStorage().GetUser()) > 0 {
		for _, userStorage := range config.Get().GetStorage().GetUser() {
			if _, ok := userService[userStorage.GetName()]; ok {
				panic(any(fmt.Errorf("Failed to init UserService: duplicate datasource: %T ", userStorage.Name)))
			}
			switch userSource := userStorage.GetStorageSource().(type) {
			case *config.Storage_Mysql:
				if client, err := mysql.NewMySQLClient(ctx, userSource.Mysql); err != nil {
					panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
				} else {
					userService[userStorage.GetName()] = mysqlservice.NewUserService(userStorage.GetName(), client)
				}
			default:
				panic(any(fmt.Errorf("Failed to init UserService: Unknown datasource: %T ", userSource)))
			}
		}
	}

	switch sessionSource := config.Get().GetStorage().GetSession().GetSource().(type) {
	case *config.Storage_Mysql:
		if client, err := mysql.NewMySQLClient(ctx, sessionSource.Mysql); err != nil {
			panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
		} else {
			sessionService = mysqlservice.NewSessionService(client)
		}
	case *config.Storage_Redis:
		if client, err := redis.NewRedisClient(ctx, sessionSource.Redis); err != nil {
			panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
		} else {
			sessionService = redisservice.NewSessionService(client)
		}
	default:
		panic(any(fmt.Errorf("初始化UserService失败: 未知的数据源类型: %T", sessionSource)))
	}
	return &Set{
		userService:    userService,
		SessionService: sessionService,
	}
}
