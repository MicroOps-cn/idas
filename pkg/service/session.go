/*
 Copyright © 2022 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service/gormservice"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/redisservice"
	w "github.com/MicroOps-cn/idas/pkg/utils/wrapper"
)

func (s Set) CreateLoginSession(ctx context.Context, username string, password string, rememberMe bool) (session string, err error) {
	users, err := s.VerifyPassword(ctx, username, password)
	if len(users) == 0 {
		return "", errors.UnauthorizedError
	}
	for _, user := range users {
		user.LoginTime = new(time.Time)
		*user.LoginTime = time.Now().UTC()
		if err = s.GetUserAndAppService(user.Storage).UpdateLoginTime(ctx, user.Id); err != nil {
			return "", err
		}
	}
	token, err := s.CreateToken(ctx, models.TokenTypeLoginSession, w.ToInterfaces[*models.User](users)...)
	if err != nil {
		return "", err
	}
	session = fmt.Sprintf("%s=%s; Path=/;Expires=%s", global.LoginSession, token.Id, token.Expiry.Format(global.LoginSessionExpiresFormat))
	return
}

type SessionService interface {
	baseService
	DeleteLoginSession(ctx context.Context, session string) error
	GetSessionByToken(ctx context.Context, sessionIds string, tokenType models.TokenType) ([]*models.User, error)

	GetSessions(ctx context.Context, userId string, current int64, size int64) (int64, []*models.Token, error)
	DeleteSession(ctx context.Context, id string) (err error)
	RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error)
	VerifyToken(ctx context.Context, token string, relationId string, tokenType models.TokenType) bool
	CreateToken(ctx context.Context, token *models.Token) error
}

func NewSessionService(ctx context.Context) SessionService {
	// logger := log.With(logs.GetContextLogger(ctx), "service", "session")
	// ctx = context.WithValue(ctx, global.LoggerName, logger)
	var sessionService SessionService
	sessionStorage := config.Get().GetStorage().GetSession()
	switch sessionSource := sessionStorage.GetStorageSource().(type) {
	case *config.Storage_Mysql:
		sessionService = gormservice.NewSessionService(sessionStorage.Name, sessionSource.Mysql.Client)
	case *config.Storage_Redis:
		sessionService = redisservice.NewSessionService(sessionStorage.Name, sessionSource.Redis)
	default:
		panic(any(fmt.Errorf("初始化SessionService失败: 未知的数据源类型: %T", sessionSource)))
	}
	return sessionService
}

func (s Set) DeleteLoginSession(ctx context.Context, session string) error {
	return s.sessionService.DeleteLoginSession(ctx, session)
}

func (s Set) GetSessionByToken(ctx context.Context, ids string, tokenType models.TokenType) ([]*models.User, error) {
	return s.sessionService.GetSessionByToken(ctx, ids, tokenType)
}

func (s Set) GetOAuthTokenByAuthorizationCode(ctx context.Context, code, clientId string) (accessToken, refreshToken string, expiresIn int, err error) {
	if users, err := s.GetSessionByToken(ctx, code, models.TokenTypeCode); err == nil && len(users) > 0 {
		_ = s.DeleteSession(ctx, code)
		at, err := s.CreateToken(ctx, models.TokenTypeToken, w.ToInterfaces[*models.User](users)...)
		if err != nil {
			return "", "", 0, err
		}
		rt, err := s.CreateToken(ctx, models.TokenTypeRefreshToken, w.ToInterfaces[*models.User](users)...)
		if err != nil {
			return "", "", 0, err
		}
		return at.Id, rt.Id, int(global.TokenExpiration / time.Minute), nil
	}

	return "", "", 0, errors.UnauthorizedError
}

func (s Set) RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error) {
	return s.sessionService.RefreshOAuthTokenByAuthorizationCode(ctx, token, clientId, clientSecret)
}

func (s Set) GetOAuthTokenByPassword(ctx context.Context, username, password string) (accessToken, refreshToken string, expiresIn int, err error) {
	if users, err := s.VerifyPassword(ctx, username, password); err != nil {
		return "", "", 0, err
	} else if len(users) == 0 {
		return "", "", 0, errors.UnauthorizedError
	} else {
		at, err := s.CreateToken(ctx, models.TokenTypeToken, w.ToInterfaces[*models.User](users)...)
		if err != nil {
			return "", "", 0, err
		}
		rt, err := s.CreateToken(ctx, models.TokenTypeRefreshToken, w.ToInterfaces[*models.User](users)...)
		if err != nil {
			return "", "", 0, err
		}
		return at.Id, rt.Id, int(global.TokenExpiration / time.Minute), nil
	}
}

func (s Set) RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error) {
	return s.sessionService.RefreshOAuthTokenByPassword(ctx, token, username, password)
}

func (s Set) GetSessions(ctx context.Context, userId string, current, size int64) (int64, []*models.Token, error) {
	return s.sessionService.GetSessions(ctx, userId, current, size)
}

func (s Set) DeleteSession(ctx context.Context, id string) (err error) {
	return s.sessionService.DeleteSession(ctx, id)
}
