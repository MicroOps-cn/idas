package service

import (
	"context"
	"fmt"
	"idas/pkg/errors"
	"idas/pkg/service/models"
)

func (s Set) UserServiceDo(name string, f func(service UserAndAppService)) error {
	service := s.GetUserAndAppService(name)
	if service == nil {
		return errors.StatusNotFound(fmt.Sprintf("User Source [%s]", name))
	}
	f(service)
	return nil
}

func (s Set) GetUserSource(_ context.Context) (data map[string]string, total int64, err error) {
	data = map[string]string{}
	for _, userService := range s.userAndAppService {
		data[userService.Name()] = userService.Name()
	}
	return
}

func (s Set) GetUsers(ctx context.Context, storage string, keywords string, status models.UserStatus, appId string, current int64, pageSize int64) (users []*models.User, total int64, err error) {
	return s.SafeGetUserAndAppService(storage).GetUsers(ctx, keywords, status, appId, current, pageSize)
}

func (s Set) PatchUsers(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.PatchUsers(ctx, patch)
}

func (s Set) DeleteUsers(ctx context.Context, storage string, id []string) (total int64, err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.DeleteUsers(ctx, id)
}

func (s Set) UpdateUser(ctx context.Context, storage string, user *models.User, updateColumns ...string) (u *models.User, err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.UpdateUser(ctx, user, updateColumns...)
}

func (s Set) GetUserInfo(ctx context.Context, storage string, id string, username string) (user *models.User, err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.GetUserInfo(ctx, id, username)
}

func (s Set) CreateUser(ctx context.Context, storage string, user *models.User) (u *models.User, err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.CreateUser(ctx, user)
}

func (s Set) PatchUser(ctx context.Context, storage string, user map[string]interface{}) (u *models.User, err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.PatchUser(ctx, user)
}

func (s Set) DeleteUser(ctx context.Context, storage string, id string) (err error) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		err = errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
		return
	}
	return service.DeleteUser(ctx, id)
}

func (s Set) VerifyPassword(ctx context.Context, username string, password string) (user *models.User, err error) {
	for _, userService := range s.userAndAppService {
		user, err = userService.VerifyPassword(ctx, username, password)
		if err == nil {
			user.Storage = userService.Name()
			return user, nil
		}
	}
	return nil, errors.UnauthorizedError
}
