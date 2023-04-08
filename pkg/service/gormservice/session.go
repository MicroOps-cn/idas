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

package gormservice

import (
	"context"

	w "github.com/MicroOps-cn/fuck/wrapper"

	"github.com/MicroOps-cn/idas/pkg/client/gorm"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

func NewSessionService(name string, client *gorm.Client) *SessionService {
	return &SessionService{name: name, Client: client}
}

type SessionService struct {
	*gorm.Client
	name string
}

func (s SessionService) UpdateToken(ctx context.Context, token *models.Token) error {
	return s.Session(ctx).Updates(token).Error
}

func (s SessionService) CreateToken(ctx context.Context, token *models.Token) error {
	conn := s.Session(ctx)
	return conn.Create(token).Error
}

func (s SessionService) GetToken(ctx context.Context, tokenId string, tokenType models.TokenType, relationId ...string) (*models.Token, error) {
	conn := s.Session(ctx)
	tk := &models.Token{}
	if err := conn.Where("id = ?", tokenId).First(tk).Error; err != nil {
		return nil, err
	}
	if tokenType != tk.Type || (len(relationId) > 0 && !w.Include[string](relationId, tk.RelationId)) {
		return nil, errors.StatusNotFound("token")
	}
	return tk, nil
}

func (s SessionService) DeleteToken(ctx context.Context, tokenType models.TokenType, id string) (err error) {
	tk := models.Token{Id: id}
	if err = s.Session(ctx).Where("token_type = ?", tokenType).Delete(&tk).Error; err != nil {
		return err
	}
	return
}

func (s SessionService) GetSessions(ctx context.Context, userId string, current, pageSize int64) (total int64, sessions []*models.Token, err error) {
	query := s.Session(ctx).Where("relation_id = ?", userId)
	if err = query.Order("last_seen").Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&sessions).Error; err != nil {
		return 0, nil, err
	} else if err = query.Count(&total).Error; err != nil {
		return 0, nil, err
	} else {
		return total, sessions, nil
	}
}

func (s SessionService) Name() string {
	return s.name
}

func (s SessionService) RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error) {
	panic("implement me")
}

func (s SessionService) RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error) {
	panic("implement me")
}

func (s SessionService) AutoMigrate(ctx context.Context) error {
	return s.Session(ctx).AutoMigrate(&models.Token{})
}

//func (s SessionService) GetSessionByToken(ctx context.Context, id string, tokenType models.TokenType, receiver interface{}) (err error) {
//	session := models.Token{Id: id}
//	if err = s.Session(ctx).Where("`type` = ?", tokenType).Omit("create_time", "user_id").First(&session).Error; err == gogorm.ErrRecordNotFound {
//		return errors.NotLoginError()
//	} else if err != nil {
//		return err
//	}
//	if session.Expiry.Before(time.Now().UTC()) {
//		return errors.NotLoginError()
//	}
//	if err = session.To(receiver); err != nil {
//		return errors.WithServerError(500, err, fmt.Sprintf("session data exception: string(data)=%s", string(session.Data)))
//	}
//	if time.Since(session.LastSeen) > time.Minute {
//		session.LastSeen = time.Now()
//		_ = s.Session(ctx).Select("last_seen").Updates(&session).Error
//	}
//	return nil
//}
