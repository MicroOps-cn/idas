package gormservice

import (
	"context"
	"encoding/json"
	"fmt"
	"idas/pkg/client/gorm"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
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

func (s SessionService) DeleteSession(ctx context.Context, id string) (err error) {
	session := models.Session{Id: id}
	if err = s.Session(ctx).Delete(&session).Error; err != nil {
		return err
	}
	return
}

func (s SessionService) GetSessions(ctx context.Context, userId string, current int64, pageSize int64) (sessions []*models.Session, total int64, err error) {
	query := s.Session(ctx).Where("user_id = ?", userId)
	if err = query.Order("last_seen").Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&sessions).Error; err != nil {
		return nil, 0, err
	} else if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	} else {
		return sessions, total, nil
	}
}

func (s SessionService) Name() string {
	return s.name
}

func (s SessionService) OAuthAuthorize(ctx context.Context, responseType, clientId, redirectURI string) (redirect string, err error) {
	panic("implement me")
}

func (s SessionService) GetOAuthTokenByAuthorizationCode(ctx context.Context, code, clientId, redirectURI string) (accessToken, refreshToken string, expiresIn int, err error) {
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
	return s.Session(ctx).AutoMigrate(&models.Session{})
}

func (s SessionService) SetLoginSession(ctx context.Context, user *models.User) (cookie string, err error) {
	sessionId := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	var session models.Session
	session.Data, err = user.MarshalJSON()
	if err != nil {
		return "", err
	}
	session.Expiry = time.Now().UTC().Add(global.LoginSessionExpiration)
	session.Key = sessionId
	session.UserId = user.Id
	if err = s.Session(ctx).Create(&session).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s=%s; Path=/;Expires=%s", global.LoginSession, sessionId, session.Expiry.Format(global.LoginSessionExpiresFormat)), nil
}

func (s SessionService) GetLoginSession(ctx context.Context, id string) (*models.User, error) {
	session := models.Session{Key: id}

	if err := s.Session(ctx).Where("`key` = ?", id).Omit("last_seen", "create_time", "user_id").First(&session).Error; err == gogorm.ErrRecordNotFound {
		return nil, errors.NotLoginError
	} else if err != nil {
		return nil, err
	}
	if session.Expiry.Before(time.Now().UTC()) {
		return nil, errors.NotLoginError
	}
	session.LastSeen = time.Now()
	_ = s.Session(ctx).Select("last_seen").Updates(&session).Error
	var user models.User
	if err := json.Unmarshal(session.Data, &user); err != nil {
		return nil, fmt.Errorf("session data exception: %s", err)
	}
	return &user, nil
}

func (s SessionService) DeleteLoginSession(ctx context.Context, id string) (string, error) {
	session := models.Session{Key: id}
	if err := s.Session(ctx).Where("`key` = ?", id).Delete(&session).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s=%s; Path=/;Expires=%s", global.LoginSession, id, time.Now().UTC().Format(global.LoginSessionExpiresFormat)), nil
}
