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

type UserServices []UserService

func (s UserServices) Include(name string) bool {
	for _, service := range s {
		if service.Name() == name {
			return true
		}
	}
	return false
}

type UserService interface {
	baseService
	Name() string
	GetUsers(ctx context.Context, keyword string, status models.UserStatus, current int64, pageSize int64) (users []*models.User, total int64, err error)
	PatchUsers(ctx context.Context, patch []map[string]interface{}) (count int64, err error)
	DeleteUsers(ctx context.Context, id []string) (count int64, err error)
	UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (*models.User, error)
	GetUserInfo(ctx context.Context, id string, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	PatchUser(ctx context.Context, user map[string]interface{}) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	VerifyPassword(ctx context.Context, username string, password string) (*models.User, error)
}

func NewUserServices(ctx context.Context) UserServices {
	var userServices UserServices
	if len(config.Get().GetStorage().GetUser()) > 0 {
		for _, userStorage := range config.Get().GetStorage().GetUser() {

			if userServices.Include(userStorage.GetName()) {
				panic(any(fmt.Errorf("Failed to init UserService: duplicate datasource: %T ", userStorage.Name)))
			}
			switch userSource := userStorage.GetStorageSource().(type) {
			case *config.Storage_Mysql:
				if client, err := mysql.NewMySQLClient(ctx, userSource.Mysql); err != nil {
					panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
				} else {
					userServices = append(userServices, mysqlservice.NewUserService(userStorage.GetName(), client))
				}
			default:
				panic(any(fmt.Errorf("Failed to init UserService: Unknown datasource: %T ", userSource)))
			}
		}
	}
	return userServices
}

func (s Set) GetUserService(name string) UserService {
	for _, userService := range s.userService {
		if userService.Name() == name /*|| len(name) == 0 */ {
			return userService
		}
	}
	return nil
}

func (s Set) UserServiceDo(name string, f func(service UserService)) error {
	service := s.GetUserService(name)
	if service == nil {
		return errors.StatusNotFound(fmt.Sprintf("User Source [%s]", name))
	}
	f(service)
	return nil
}

func (s Set) SafeGetUserService(name string) UserService {
	for _, userService := range s.userService {
		if userService.Name() == name {
			return userService
		}
	}
	for _, service := range s.userService {
		return service
	}
	return nil
}

func (s Set) GetUserSource(_ context.Context) (data map[string]string, total int64, err error) {
	data = map[string]string{}
	for _, userService := range s.userService {
		data[userService.Name()] = userService.Name()
	}
	return
}

func (s Set) GetUsers(ctx context.Context, storage string, keyword string, status models.UserStatus, current int64, pageSize int64) (users []*models.User, total int64, err error) {
	return s.SafeGetUserService(storage).GetUsers(ctx, keyword, status, current, pageSize)
}

func (s Set) PatchUsers(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, err error) {
	service := s.GetUserService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.PatchUsers(ctx, patch)
	}
}

func (s Set) DeleteUsers(ctx context.Context, storage string, id []string) (total int64, err error) {
	service := s.GetUserService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.DeleteUsers(ctx, id)
	}
}

func (s Set) UpdateUser(ctx context.Context, storage string, user *models.User, updateColumns ...string) (u *models.User, err error) {
	service := s.GetUserService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.UpdateUser(ctx, user, updateColumns...)
	}
}

func (s Set) GetUserInfo(ctx context.Context, storage string, id string, username string) (user *models.User, err error) {
	service := s.GetUserService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.GetUserInfo(ctx, id, username)
	}
}

func (s Set) CreateUser(ctx context.Context, storage string, user *models.User) (u *models.User, err error) {
	service := s.GetUserService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.CreateUser(ctx, user)
	}
}

func (s Set) PatchUser(ctx context.Context, storage string, user map[string]interface{}) (u *models.User, err error) {
	service := s.GetUserService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.PatchUser(ctx, user)
	}
}

func (s Set) DeleteUser(ctx context.Context, storage string, id string) (err error) {
	service := s.GetUserService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	} else {
		return service.DeleteUser(ctx, id)
	}
}

func (s Set) VerifyPassword(ctx context.Context, username string, password string) (user *models.User, err error) {
	for _, userService := range s.userService {
		user, err = userService.VerifyPassword(ctx, username, password)
		if err == nil {
			user.Storage = userService.Name()
			return user, nil
		}
	}
	return nil, errors.UnauthorizedError
}
