package redisservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	"idas/pkg/errors"
	"idas/pkg/logs"
	"time"

	"idas/pkg/client/redis"
	"idas/pkg/global"
	"idas/pkg/service/models"
)

type SessionService struct {
	*redis.Client
	name string
}

func (s SessionService) CreateToken(ctx context.Context, token *models.Token) error {
	sessionId := models.NewId()
	redisClt := s.Redis(ctx)
	if err := redisClt.Set(fmt.Sprintf("%s:%s", global.LoginSession, sessionId), token, -time.Since(token.Expiry)).Err(); err != nil {
		return err
	}
	return nil
}

func (s SessionService) VerifyToken(ctx context.Context, token string, relationId string, tokenType models.TokenType) bool {
	//TODO implement me
	panic("implement me")
}

func (s SessionService) DeleteSession(ctx context.Context, id string) (err error) {
	panic("implement me")
}

func (s SessionService) GetSessions(ctx context.Context, userId string, current int64, size int64) ([]*models.Token, int64, error) {
	panic("implement me")
}

func (s SessionService) Name() string {
	return s.name
}

func (s SessionService) OAuthAuthorize(ctx context.Context, clientId string) (code string, err error) {
	panic("implement me")
}

func (s SessionService) CreateOAuthAuthCode(ctx context.Context, appId, sessionId, scope, storage string) (code string, err error) {
	redisClt := s.Redis(ctx)
	c := models.AppAuthCode{Model: models.Model{Id: models.NewId()}, AppId: appId, SessionId: sessionId, Scope: scope, Storage: storage}

	if cc, err := json.Marshal(c); err != nil {
		return "", err
	} else if err := redisClt.Set(fmt.Sprintf("%s:%s", global.AuthCode, c.Id), cc, global.AuthCodeExpiration).Err(); err != nil {
		return "", err
	} else {
		return c.Id, nil
	}
}

func (s SessionService) GetUserByOAuthAuthorizationCode(ctx context.Context, code, clientId string) (user *models.User, scope string, err error) {
	logger := logs.GetContextLogger(ctx)
	redisClt := s.Redis(ctx)
	c := new(models.AppAuthCode)
	sessionValue, err := redisClt.Get(fmt.Sprintf("%s:%s", global.AuthCode, code)).Bytes()
	if err != nil {
		level.Error(logger).Log("msg", "failed to get auth code info", "err", err)
		return nil, "", errors.BadRequestError
	} else if err = json.Unmarshal(sessionValue, c); err != nil {
		level.Error(logger).Log("msg", "failed to parse auth code info", "err", err)
		return nil, "", errors.BadRequestError
	} else if clientId != c.AppId {
		level.Error(logger).Log("msg", "client id is not match", "err", err)
		return nil, "", errors.BadRequestError
	}
	if users, err := s.GetLoginSession(ctx, c.SessionId); err != nil {
		level.Error(logger).Log("msg", "failed to get session info", "err", err)
		return nil, "", errors.BadRequestError
	} else if len(users) < 0 {
		level.Error(logger).Log("msg", "session expired")
		return nil, "", errors.BadRequestError
	} else {
		return users[0], c.Scope, nil
	}
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

func (s SessionService) GetLoginSession(ctx context.Context, id string) (users []*models.User, err error) {
	redisClt := s.Redis(ctx)
	sessionValue, err := redisClt.Get(fmt.Sprintf("%s:%s", global.LoginSession, id)).Bytes()
	if err != nil {
		return nil, err
	} else if err = json.Unmarshal(sessionValue, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s SessionService) SetLoginSession(ctx context.Context, user *models.User) (string, error) {
	sessionId := models.NewId()
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
