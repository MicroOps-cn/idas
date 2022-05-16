package redisservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"idas/pkg/client/redis"
	"idas/pkg/global"
	"idas/pkg/service/models"
)

type SessionService struct {
	*redis.Client
	name string
}

func (s SessionService) DeleteSession(ctx context.Context, id string) (err error) {
	panic("implement me")
}

func (s SessionService) GetSessions(ctx context.Context, userId string, current int64, size int64) ([]*models.Session, int64, error) {
	panic("implement me")
}

func (s SessionService) Name() string {
	return s.name
}

func (s SessionService) OAuthAuthorize(ctx context.Context, clientId string) (code string, err error) {
	panic("implement me")
}

func (s SessionService) GetOAuthTokenByAuthorizationCode(ctx context.Context, code, clientId string) (accessToken, refreshToken string, expiresIn int, err error) {
	panic("implement me")
}

func (s SessionService) GetOAuthTokenByPassword(ctx context.Context, username string, password string) (accessToken, refreshToken string, expiresIn int, err error) {
	panic("implement me")
}

func (s SessionService) RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error) {
	panic("implement me")
}

func (s SessionService) RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error) {
	panic("implement me")
}

func (s SessionService) AutoMigrate(ctx context.Context) error {
	return nil
}

func (s SessionService) GetLoginSession(ctx context.Context, ids []string) (users []*models.User, err error) {
	for _, id := range ids {
		user := new(models.User)
		redisClt := s.Redis(ctx)
		sessionValue, err := redisClt.Get(fmt.Sprintf("%s:%s", global.LoginSession, id)).Bytes()
		if err != nil {
			return nil, err
		} else if err = json.Unmarshal(sessionValue, user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s SessionService) SetLoginSession(ctx context.Context, user *models.User) (string, error) {
	sessionId := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	redisClt := s.Redis(ctx)
	user.Password = nil
	user.Salt = nil
	if userb, err := json.Marshal(user); err != nil {
		return "", err
	} else if err := redisClt.Set(fmt.Sprintf("%s:%s", global.LoginSession, sessionId), userb, global.LoginSessionExpiration).Err(); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s=%s; Path=/;Expires=%s", global.LoginSession, sessionId, time.Now().UTC().Add(global.LoginSessionExpiration).Format(global.LoginSessionExpiresFormat)), nil
	}
}

func (s SessionService) DeleteLoginSession(ctx context.Context, sessionId string) error {
	redisClt := s.Redis(ctx)
	_ = redisClt.Del(fmt.Sprintf("%s:%s", global.LoginSession, sessionId)).Val()
	return nil
}

func NewSessionService(name string, client *redis.Client) *SessionService {
	return &SessionService{name: name, Client: client}
}
