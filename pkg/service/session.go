package service

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"idas/pkg/client/gorm"
	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/logs"
	"idas/pkg/service/gormservice"
	"time"

	"idas/config"
	"idas/pkg/client/redis"
	"idas/pkg/service/models"
	"idas/pkg/service/redisservice"
)

func (s Set) CreateLoginSession(ctx context.Context, username string, password string) (session []string, err error) {
	users, err := s.VerifyPassword(ctx, username, password)
	if len(users) == 0 {
		return nil, errors.UnauthorizedError
	}
	for _, user := range users {
		user.LoginTime = new(time.Time)
		*user.LoginTime = time.Now().UTC()
		if user, err = s.GetUserAndAppService(user.Storage).UpdateUser(ctx, user, "login_time"); err != nil {
			return nil, err
		}

		loginSession, err := s.sessionService.SetLoginSession(ctx, user)
		if err != nil {
			return nil, err
		}
		session = append(session, loginSession)
	}
	return
}

type SessionService interface {
	baseService
	SetLoginSession(ctx context.Context, user *models.User) (string, error)
	DeleteLoginSession(ctx context.Context, session string) error
	GetLoginSession(ctx context.Context, sessionIds []string) ([]*models.User, error)

	GetSessions(ctx context.Context, userId string, current int64, size int64) ([]*models.Session, int64, error)
	DeleteSession(ctx context.Context, id string) (err error)

	CreateOAuthAuthCode(ctx context.Context, appId, sessionId, scope, storage string) (code string, err error)
	CreateOAuthAccessToken(ctx context.Context, appId, sessionId, scope, storage string) (code string, err error)
	GetUserByOAuthAuthorizationCode(ctx context.Context, code, clientId string) (user *models.User, scope string, err error)
	RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error)
	GetOAuthTokenByPassword(ctx context.Context, username string, password string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error)
}

func NewSessionService(ctx context.Context) SessionService {
	logger := log.With(logs.GetContextLogger(ctx), "service", "session")
	ctx = context.WithValue(ctx, global.LoggerName, logger)
	var sessionService SessionService
	sessionStorage := config.Get().GetStorage().GetSession()
	switch sessionSource := sessionStorage.GetStorageSource().(type) {
	case *config.Storage_Mysql:
		if client, err := gorm.NewMySQLClient(ctx, sessionSource.Mysql); err != nil {
			panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
		} else {
			sessionService = gormservice.NewSessionService(sessionStorage.Name, client)
		}
	case *config.Storage_Redis:
		if client, err := redis.NewRedisClient(ctx, sessionSource.Redis); err != nil {
			panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
		} else {
			sessionService = redisservice.NewSessionService(sessionStorage.Name, client)
		}
	default:
		panic(any(fmt.Errorf("初始化SessionService失败: 未知的数据源类型: %T", sessionSource)))
	}
	return sessionService
}

func (s Set) SetLoginSession(ctx context.Context, user *models.User) (string, error) {
	return s.sessionService.SetLoginSession(ctx, user)
}

func (s Set) DeleteLoginSession(ctx context.Context, session string) error {
	return s.sessionService.DeleteLoginSession(ctx, session)
}

func (s Set) GetLoginSession(ctx context.Context, ids []string) ([]*models.User, error) {
	return s.sessionService.GetLoginSession(ctx, ids)
}

func (s Set) GetOAuthTokenByAuthorizationCode(ctx context.Context, code, clientId string) (accessToken, refreshToken string, expiresIn int, err error) {
	user, _, err := s.sessionService.GetUserByOAuthAuthorizationCode(ctx, code, clientId)
	if err != nil {
		return "", "", 0, err
	}
	if len(user.Id) == 0 {
		return "", "", 0, errors.BadRequestError
	}
	//info, err := s.GetUserInfo(ctx, storage, userId, "")
	s.sessionService.CreateOAuthAccessToken()
	return "", "", 0, err
}

func (s Set) RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error) {
	return s.sessionService.RefreshOAuthTokenByAuthorizationCode(ctx, token, clientId, clientSecret)
}

func (s Set) GetOAuthTokenByPassword(ctx context.Context, username, password string) (accessToken, refreshToken string, expiresIn int, err error) {
	return s.sessionService.GetOAuthTokenByPassword(ctx, username, password)
}

func (s Set) RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error) {
	return s.sessionService.RefreshOAuthTokenByPassword(ctx, token, username, password)
}

func (s Set) GetSessions(ctx context.Context, userId string, current, size int64) ([]*models.Session, int64, error) {
	return s.sessionService.GetSessions(ctx, userId, current, size)
}

func (s Set) DeleteSession(ctx context.Context, id string) (err error) {
	return s.sessionService.DeleteSession(ctx, id)
}
