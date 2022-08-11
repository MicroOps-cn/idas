package service

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"

	"idas/pkg/errors"
	"idas/pkg/service/models"
	"idas/pkg/utils/sign"
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

func (s Set) GetUsers(ctx context.Context, storage string, keywords string, status models.UserMeta_UserStatus, appId string, current, pageSize int64) (total int64, users []*models.User, err error) {
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

func (s Set) CreateUserKey(ctx context.Context, userId, name string) (keyPair *models.UserKey, err error) {
	return s.commonService.CreateUserKeyWithId(ctx, userId, name)
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

func (s Set) VerifyPasswordById(ctx context.Context, storage, username, password string) (users []*models.User) {
	service := s.GetUserAndAppService(storage)
	if service == nil {
		return service.VerifyPassword(ctx, username, password)
	}
	return nil
}

func (s Set) VerifyPassword(ctx context.Context, username string, password string) (users []*models.User, err error) {
	for _, userService := range s.userAndAppService {
		for _, user := range userService.VerifyPassword(ctx, username, password) {
			if user.Status == models.UserMeta_inactive || user.Status == models.UserMeta_disable {
				continue
			}
			user.Storage = userService.Name()
			users = append(users, user)
		}
	}
	return users, nil
}

func (s Set) Authentication(ctx context.Context, method models.AuthMeta_Method, algorithm sign.AuthAlgorithm, key, secret, payload, signStr string) ([]*models.User, error) {
	if method == models.AuthMeta_basic {
		if _, err := uuid.FromString(key); err != nil {
			return s.VerifyPassword(ctx, key, secret)
		}
	}
	userKey, err := s.commonService.GetUserKey(ctx, key)
	if err != nil {
		return nil, err
	} else if userKey == nil {
		return nil, nil
	} else {
		for _, service := range s.userAndAppService {
			if info, err := service.GetUserInfo(ctx, userKey.UserId, ""); err != nil {
				continue
			} else {
				userKey.User = info
				break
			}
		}
	}
	if userKey.User == nil {
		return nil, errors.StatusNotFound("user")
	}
	switch method {
	case models.AuthMeta_basic:
		if userKey.Secret == secret {
			return []*models.User{userKey.User}, nil
		}
	case models.AuthMeta_signature:
		if sign.Verify(userKey.Key, userKey.Secret, userKey.Private, algorithm, signStr, payload) {
			return []*models.User{userKey.User}, nil
		}
		return nil, errors.ParameterError("Failed to verify the signature")
	//case models.AuthMeta_token:
	//	if s.VerifyToken(ctx, secret, "", models.TokenTypeToken) {
	//		return nil, nil
	//	}
	//	return nil, errors.ParameterError("Failed to verify the signature")
	default:
		return nil, errors.ParameterError("unknown auth method")
	}
	return nil, errors.ParameterError("unknown auth request")
}

func (s Set) GetAuthCodeByClientId(ctx context.Context, clientId string, user *models.User, sessionId, storage string) (code string, err error) {
	svc := s.GetUserAndAppService(storage)
	if svc == nil {
		return "", errors.StatusNotFound(fmt.Sprintf("App Source [%s]", storage))
	}
	user.Role, err = svc.VerifyUserAuthorizationForApp(ctx, clientId, user.Id)
	if err != nil {
		return "", err
	}

	token, err := s.CreateToken(ctx, models.TokenTypeCode, user)
	if err != nil {
		return "", err
	}
	return token.Id, nil
}
