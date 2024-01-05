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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MicroOps-cn/fuck/clients/redis"
	g "github.com/MicroOps-cn/fuck/generator"
	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/fuck/sets"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/go-kit/log/level"
	goredis "github.com/go-redis/redis"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

type SessionService struct {
	*redis.Client
	name string
}

func getCounterKey(seed string) string {
	return strings.Join([]string{global.RedisKeyPrefix, "COUNTER", fmt.Sprintf("%x", sha256.Sum256([]byte(seed)))}, ":")
}

func (s SessionService) Counter(ctx context.Context, seed string, expireTime *time.Time, num ...int64) (err error) {
	redisClt := s.Redis(ctx)
	var ret *goredis.IntCmd
	key := getCounterKey(seed)
	if len(num) == 0 {
		if ret = redisClt.Incr(key); ret.Err() != nil {
			return ret.Err()
		}
	} else {
		sum := int64(0)
		for _, n := range num {
			sum += n
		}
		if sum != 0 {
			if ret = redisClt.IncrBy(key, sum); ret.Err() != nil {
				return ret.Err()
			}
		}
	}
	if expireTime != nil {
		redisClt.ExpireAt(key, *expireTime)
	}
	return nil
}

func (s SessionService) GetCounter(ctx context.Context, seed string) (count int64, err error) {
	count, err = s.Redis(ctx).Get(getCounterKey(seed)).Int64()
	if err != nil && err != goredis.Nil {
		return 0, err
	}
	return count, nil
}

func (s SessionService) UpdateToken(ctx context.Context, token *models.Token) error {
	return s.Redis(ctx).Set(getTokenKey(token.Id), NewToken(token), -time.Since(token.Expiry)).Err()
}

func (s SessionService) UpdateTokenExpires(ctx context.Context, id string, expiry time.Time) error {
	return s.Redis(ctx).ExpireAt(getTokenKey(id), expiry).Err()
}

type Token struct {
	Id          string           `json:"id"`
	CreateTime  time.Time        `json:"createTime"`
	Data        string           `json:"data,omitempty"`
	RelationId  string           `json:"relationId"`
	Expiry      time.Time        `json:"expiry"`
	Type        models.TokenType `json:"type"`
	LastSeen    time.Time        `json:"lastSeen"`
	ParentId    string           `json:"parentId,omitempty"`
	Childrens   []*Token         `json:"-"`
	ChildrensId []string         `json:"childrensId,omitempty"`
}

func NewToken(token *models.Token) *Token {
	tk := &Token{
		Id:         token.Id,
		CreateTime: token.CreateTime,
		Data:       string(token.Data),
		ParentId:   token.ParentId,
		LastSeen:   token.LastSeen,
		Expiry:     token.Expiry,
		Type:       token.Type,
		RelationId: token.RelationId,
	}
	return tk
}

func (s *Token) ToToken() *models.Token {
	tk := &models.Token{
		Id:         s.Id,
		CreateTime: s.CreateTime,
		Data:       []byte(s.Data),
		ParentId:   s.ParentId,
		LastSeen:   s.LastSeen,
		Expiry:     s.Expiry,
		Type:       s.Type,
		RelationId: s.RelationId,
	}
	return tk
}

func (s *Token) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *Token) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}

func getTokensKey(tokenId string) string {
	return strings.Join([]string{global.RedisKeyPrefix, "TOKENS", fmt.Sprintf("%x", sha256.Sum256([]byte(tokenId)))}, ":")
}

func getTokenKey(tokenId string) string {
	return strings.Join([]string{global.RedisKeyPrefix, "TOKEN", fmt.Sprintf("%x", sha256.Sum256([]byte(tokenId)))}, ":")
}

func (s SessionService) CreateToken(ctx context.Context, token *models.Token) error {
	redisClt := s.Redis(ctx)
	if len(token.Id) == 0 {
		token.Id = g.NewUUID(token.RelationId).String()
	}
	if token.CreateTime.IsZero() {
		token.CreateTime = time.Now().UTC()
	}
	tk := NewToken(token)
	// 创建Token
	if err := redisClt.Set(getTokenKey(token.Id), tk, -time.Since(token.Expiry)).Err(); err != nil {
		return err
	}

	if len(token.RelationId) > 0 {
		objTokenListKey := getTokensKey(token.RelationId)
		err := redisClt.SAdd(objTokenListKey, token.Id).Err()
		if err != nil {
			return err
		}
		ttl, err := redisClt.TTL(objTokenListKey).Result()
		if err != nil || ttl < 0 || time.Since(token.Expiry) > ttl {
			if err = redisClt.ExpireAt(objTokenListKey, token.Expiry).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s SessionService) GetToken(ctx context.Context, tokenId string, tokenType models.TokenType, relationId ...string) (*models.Token, error) {
	var token Token
	redisClt := s.Redis(ctx)
	if err := redisClt.Get(getTokenKey(tokenId)).Scan(&token); err != nil {
		return nil, nil
	}
	retToken := models.Token{
		Id:         token.Id,
		CreateTime: token.CreateTime,
		Data:       []byte(token.Data),
		ParentId:   token.ParentId,
		LastSeen:   token.LastSeen,
		Expiry:     token.Expiry,
		Type:       token.Type,
		RelationId: token.RelationId,
	}
	if token.Type != tokenType || (len(relationId) > 0 && !w.Include[string](relationId, token.RelationId)) {
		return nil, errors.StatusNotFound("token")
	}
	return &retToken, nil
}

func (s SessionService) DeleteToken(ctx context.Context, tokenType models.TokenType, id string) (err error) {
	conn := s.Redis(ctx)
	token, err := s.GetToken(ctx, id, tokenType)
	keys := []string{getTokenKey(id)}
	if err != nil {
		level.Warn(logs.GetContextLogger(ctx)).Log("msg", "failed to get token", "err", err)
		return conn.Del(getTokenKey(id)).Err()
	}
	conn.SRem(getTokensKey(token.RelationId), token.Id)

	return conn.Del(keys...).Err()
}

func (s SessionService) GetSessions(ctx context.Context, userId string, current int64, pageSize int64) (int64, []*models.Token, error) {
	redisClt := s.Redis(ctx)
	pChildKey := getTokensKey(userId)
	expireKeys := sets.New[string]()
	var ret []*models.Token
	err := redis.ForeachSet(ctx, redisClt, pChildKey, uint64((current-1)*pageSize), pageSize, func(key, tokenId string) error {
		var token Token
		tokenKey := getTokenKey(tokenId)
		if err := redisClt.Get(tokenKey).Scan(&token); err == goredis.Nil {
			expireKeys.Insert(tokenKey)
		} else if err != nil {
			return err
		} else {
			if token.Type == models.TokenTypeLoginSession {
				ret = append(ret, token.ToToken())
			}
			if int64(len(ret)) >= pageSize {
				return redis.ErrStopLoop
			}
		}
		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	if expireKeys.Len() > 0 {
		err = redisClt.SRem(pChildKey, w.Interfaces(expireKeys.List())...).Err()
		if err != nil {
			return 0, nil, err
		}
	}
	count, err := redisClt.SCard(pChildKey).Result()
	if err != nil {
		return 0, nil, err
	}
	return count, ret, nil
}

func (s SessionService) Name() string {
	return s.name
}

//func (s SessionService) OAuthAuthorize(ctx context.Context, clientId string) (code string, err error) {
//	panic("implement me")
//}

func (s SessionService) RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error) {
	panic("implement me")
}

func (s SessionService) RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error) {
	panic("implement me")
}

func (s SessionService) AutoMigrate(ctx context.Context) error {
	return nil
}

//func (s SessionService) GetSessionByToken(ctx context.Context, id string, tokenType models.TokenType, receiver interface{}) (err error) {
//	redisClt := s.Redis(ctx)
//	sessionValue, err := redisClt.Get(fmt.Sprintf("%s:%s", global.LoginSession, id)).Bytes()
//	if err != nil {
//		return err
//	} else if err = json.Unmarshal(sessionValue, receiver); err != nil {
//		return err
//	}
//	return nil
//}

//func (s SessionService) SetLoginSession(ctx context.Context, user *models.User) (string, error) {
//	sessionId := models.NewId()
//	redisClt := s.Redis(ctx)
//	user.Password = nil
//	user.Salt = nil
//	if userb, err := json.Marshal(user); err != nil {
//		return "", err
//	} else if err := redisClt.Set(fmt.Sprintf("%s:%s", global.LoginSession, sessionId), userb, global.LoginSessionExpiration).Err(); err != nil {
//		return "", err
//	} else {
//		return fmt.Sprintf("%s=%s; Path=/;Expires=%s", global.LoginSession, sessionId, time.Now().UTC().Add(global.LoginSessionExpiration).Format(global.LoginSessionExpiresFormat)), nil
//	}
//}

func NewSessionService(name string, client *redis.Client) *SessionService {
	return &SessionService{name: name, Client: client}
}
