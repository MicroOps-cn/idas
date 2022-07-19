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

func (s Set) GetUserSource(_ context.Context) (total int64, data map[string]string, err error) {
	data = map[string]string{}
	for _, userService := range s.userAndAppService {
		data[userService.Name()] = userService.Name()
	}

	return
}

func (s Set) GetUsers(ctx context.Context, storage string, keywords string, status models.UserStatus, appId string, current, pageSize int64) (total int64, users []*models.User, err error) {
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

func (s Set) VerifyPassword(ctx context.Context, username string, password string) (users []*models.User, err error) {
	for _, userService := range s.userAndAppService {
		for _, user := range userService.VerifyPassword(ctx, username, password) {
			user.Storage = userService.Name()
			users = append(users, user)
		}
	}
	return users, nil
}

func (s Set) Authentication(ctx context.Context, method models.AuthMeta_Method, algorithm models.AuthAlgorithm, key, secret string) ([]*models.User, error) {
	switch method {
	case models.AuthMeta_basic:
		return s.VerifyPassword(ctx, key, secret)
	case models.AuthMeta_signature:
		switch algorithm {
		case "", "HMAC-SHA1":
		}
	default:
		return nil, errors.ParameterError("unknown auth method")
	}
	return nil, errors.ParameterError("unknown auth request")
}
func (s Set) GetAuthCodeByClientId(ctx context.Context, clientId, userId, sessionId, storage string) (code string, err error) {
	svc := s.GetUserAndAppService(storage)
	if svc == nil {
		return "", errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
	}
	scope, err := svc.VerifyUserAuthorizationForApp(ctx, clientId, userId)
	if err != nil {
		return "", err
	}
	return s.sessionService.CreateOAuthAuthCode(ctx, clientId, sessionId, scope, storage)
}
