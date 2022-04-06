package service

import (
	"context"
	"fmt"
	"time"

	"idas/config"
	"idas/pkg/client/mysql"
	"idas/pkg/client/redis"
	"idas/pkg/service/models"
	"idas/pkg/service/mysqlservice"
	"idas/pkg/service/redisservice"
)

func (s Set) CreateLoginSession(ctx context.Context, username string, password string) (session string, err error) {
	user, err := s.VerifyPassword(ctx, username, password)
	if user == nil {
		return "", err
	}
	user.LoginTime = time.Now().UTC()
	if user, err = s.GetUserService(user.Storage).UpdateUser(ctx, user, "login_time"); err != nil {
		return "", err
	}
	return s.SessionService.SetLoginSession(ctx, user)
}

type SessionService interface {
	baseService
	SetLoginSession(ctx context.Context, user *models.User) (string, error)
	DeleteLoginSession(ctx context.Context, session string) (string, error)
	GetLoginSession(ctx context.Context, id string) (*models.User, error)
	OAuthAuthorize(ctx context.Context, responseType, clientId, redirectURI string) (redirect string, err error)
	GetOAuthTokenByAuthorizationCode(ctx context.Context, code, clientId, redirectURI string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error)
	GetOAuthTokenByPassword(ctx context.Context, username string, password string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error)
	GetSessions(ctx context.Context, userId string, current int64, size int64) ([]*models.Session, int64, error)
	DeleteSession(ctx context.Context, id string) (err error)
}

func NewSessionService(ctx context.Context) SessionService {
	var sessionService SessionService
	sessionStorage := config.Get().GetStorage().GetSession()
	switch sessionSource := sessionStorage.GetStorageSource().(type) {
	case *config.Storage_Mysql:
		if client, err := mysql.NewMySQLClient(ctx, sessionSource.Mysql); err != nil {
			panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
		} else {
			sessionService = mysqlservice.NewSessionService(sessionStorage.Name, client)
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
