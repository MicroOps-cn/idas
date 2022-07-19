package gormservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	"idas/pkg/client/gorm"
	"idas/pkg/logs"
	"time"

	gogorm "gorm.io/gorm"

	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/service/models"
)

func NewSessionService(name string, client *gorm.Client) *SessionService {
	return &SessionService{name: name, Client: client}
}

type SessionService struct {
	*gorm.Client
	name string
}

func (s SessionService) CreateToken(ctx context.Context, token *models.Token) error {
	return s.Session(ctx).Create(token).Error
}

func (s SessionService) VerifyToken(ctx context.Context, token string, relationId string, tokenType models.TokenType) bool {
	conn := s.Session(ctx)
	tk := &models.Token{}
	if err := conn.Model(&models.Token{}).Where("id = ? and relation_id = ? and `type` = ?", token, relationId, tokenType).First(tk).Error; err != nil {
		return false
	}
	if tokenType == models.TokenTypeResetPassword {
		if err := conn.Delete(tk).Error; err != nil {
			return false
		}
	}
	if tk.Expiry.After(time.Now().UTC()) {
		return true
	}
	return false
}

func (s SessionService) DeleteSession(ctx context.Context, id string) (err error) {
	session := models.Token{Id: id}
	if err = s.Session(ctx).Delete(&session).Error; err != nil {
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

func (s SessionService) CreateOAuthAuthCode(ctx context.Context, appId, sessionId string, scope, storage string) (code string, err error) {
	c := models.AppAuthCode{AppId: appId, SessionId: sessionId, Scope: scope, Storage: storage}
	if err = s.Session(ctx).Create(&c).Error; err != nil {
		return "", err
	}
	return c.Id, nil
}

func (s SessionService) GetUserByOAuthAuthorizationCode(ctx context.Context, code, clientId string) (user *models.User, scope string, err error) {
	logger := logs.GetContextLogger(ctx)
	c := models.AppAuthCode{}
	if err = s.Session(ctx).Where("id = ? and app_id = ? and create_time > ?", code, clientId, time.Now().Add(-global.AuthCodeExpiration)).First(&c).Error; err != nil {
		if err != gogorm.ErrRecordNotFound {
			level.Error(logger).Log("msg", "failed to get auth code info", "err", err)
		}
		return nil, "", errors.BadRequestError
	} else if err = s.Session(ctx).Delete(&c).Error; err != nil {
		level.Error(logger).Log("msg", "failed to remove auth code", "err", err)
		return
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
	return s.Session(ctx).AutoMigrate(&models.Token{})
}

func (s SessionService) GetLoginSession(ctx context.Context, id string) (users []*models.User, err error) {
	session := models.Token{Id: id}
	if err = s.Session(ctx).Where("`type` = ?", models.TokenTypeLoginSession).Omit("last_seen", "create_time", "user_id").First(&session).Error; err == gogorm.ErrRecordNotFound {
		return nil, errors.NotLoginError
	} else if err != nil {
		return nil, err
	}
	if session.Expiry.Before(time.Now().UTC()) {
		return nil, errors.NotLoginError
	}
	if err = json.Unmarshal(session.Data, &users); err != nil {
		return nil, fmt.Errorf("session data exception: %s,data=%s,string(data)=%s", err, session.Data, string(session.Data))
	}
	session.LastSeen = time.Now()
	_ = s.Session(ctx).Select("last_seen").Updates(&session).Error
	return users, nil
}

func (s SessionService) DeleteLoginSession(ctx context.Context, id string) error {
	session := models.Token{Id: id}
	return s.Session(ctx).Select("`type` = ?", models.TokenTypeLoginSession).Delete(&session).Error
}
