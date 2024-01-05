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
	"crypto/md5"
	"time"

	"github.com/MicroOps-cn/fuck/clients/gorm"
	logs "github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/go-kit/log/level"
	uuid "github.com/satori/go.uuid"
	goorm "gorm.io/gorm"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/signals"
)

func NewSessionService(ctx context.Context, name string, client *gorm.Client) *SessionService {
	s := &SessionService{name: name, Client: client}
	go s.startBroom(ctx)
	return s
}

type SessionService struct {
	*gorm.Client
	name string
}

func (s SessionService) startBroom(ctx context.Context) {
	if ctx.Value("command") == "init" {
		return
	}
	ticker := time.NewTicker(time.Hour)
	stopCh := signals.SetupSignalHandler(logs.GetContextLogger(ctx))
	stopCh.AddRequest(1)
loop:
	for {
		logger := logs.NewTraceLogger()
		conn := s.Session(ctx)
		if err := conn.Where("expiry < ?", time.Now().UTC()).Delete(&models.Token{}).Error; err != nil {
			level.Error(logger).Log("msg", "Failed to delete expired token.", "err", err)
		}
		if err := conn.Where("expire_time < ?", time.Now().UTC()).Delete(&models.Counter{}).Error; err != nil {
			level.Error(logger).Log("msg", "Failed to delete expired counter.", "err", err)
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			level.Debug(logs.GetContextLogger(ctx)).Log("msg", "close session broom")
			break loop
		}
	}
	stopCh.DoneRequest()
}

func getCounterKey(seed string) string {
	sum := md5.Sum([]byte(seed))
	u, err := uuid.FromBytes(sum[:])
	if err != nil {
		panic("system error")
	}

	return u.String()
}

func (s SessionService) Counter(ctx context.Context, seed string, expireTime *time.Time, num ...int64) (err error) {
	tx := s.Session(ctx).Begin()
	defer tx.Rollback()
	c := models.Counter{Id: getCounterKey(seed), Seed: seed, ExpireTime: expireTime, Count: 1}
	if len(c.Seed) > 128 {
		c.Seed = c.Seed[:128]
	}
	if len(num) != 0 {
		sum := int64(0)
		for _, n := range num {
			sum += n
		}
		c.Count = sum
	}
	ct := c.Count
	if ret := tx.Where("id = ?", c.Id).FirstOrCreate(&c); ret.Error != nil {
		return ret.Error
	} else if ret.RowsAffected == 0 {
		if ct != 0 {
			if err = tx.Model(&c).UpdateColumn("count", goorm.Expr("count + ?", ct)).Error; err != nil {
				return err
			}
		}
	}
	return tx.Commit().Error
}

func (s SessionService) GetCounter(ctx context.Context, seed string) (count int64, err error) {
	conn := s.Session(ctx)
	var c models.Counter
	if err = conn.Where("id = ?", getCounterKey(seed)).Select("id", "count").First(&c).Error; err != nil && err != goorm.ErrRecordNotFound {
		return 0, err
	}
	return c.Count, nil
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
	if err = s.Session(ctx).Where("`type` = ?", tokenType).Delete(&tk).Error; err != nil {
		return err
	}
	return
}

func (s SessionService) GetSessions(ctx context.Context, userId string, current, pageSize int64) (total int64, sessions []*models.Token, err error) {
	query := s.Session(ctx).Where("relation_id = ?", userId).Where("`type` = ?", models.TokenTypeLoginSession)
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
	return s.Session(ctx).AutoMigrate(&models.Token{}, &models.Counter{})
}

func (s SessionService) UpdateTokenExpires(ctx context.Context, id string, expiry time.Time) error {
	return s.Session(ctx).Select("expiry").Updates(&models.Token{Id: id, Expiry: expiry}).Error
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
