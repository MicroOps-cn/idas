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

package redisservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MicroOps-cn/idas/pkg/client/redis"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

type SessionService struct {
	*redis.Client
	name string
}

func (s SessionService) CreateToken(ctx context.Context, token *models.Token) error {
	sessionId := models.NewId()
	redisClt := s.Redis(ctx)
	return redisClt.Set(fmt.Sprintf("%s:%s", global.LoginSession, sessionId), token, -time.Since(token.Expiry)).Err()
}

func (s SessionService) VerifyToken(ctx context.Context, token string, relationId string, tokenType models.TokenType) bool {
	panic("implement me")
}

func (s SessionService) DeleteSession(ctx context.Context, id string) (err error) {
	panic("implement me")
}

func (s SessionService) GetSessions(ctx context.Context, userId string, current int64, size int64) (int64, []*models.Token, error) {
	panic("implement me")
}

func (s SessionService) Name() string {
	return s.name
}

func (s SessionService) OAuthAuthorize(ctx context.Context, clientId string) (code string, err error) {
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

func (s SessionService) GetSessionByToken(ctx context.Context, id string, tokenType models.TokenType) (users []*models.User, err error) {
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
