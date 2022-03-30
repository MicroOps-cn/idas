package mysqlservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"idas/pkg/client/mysql"
	"idas/pkg/global"
	"idas/pkg/service/models"
)

type SessionService struct {
	*mysql.Client
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
	if err = s.Session(ctx).Create(&session).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s=%s; Path=/;Expires=%s", global.LoginSession, sessionId, session.Expiry.Format(global.LoginSessionExpiresFormat)), nil
}

func NewSessionService(client *mysql.Client) *SessionService {
	return &SessionService{Client: client}
}

func (s SessionService) GetLoginSession(ctx context.Context, id string) (*models.User, string, error) {
	session := models.Session{Key: id}

	if err := s.Session(ctx).First(&session).Error; err != nil {
		return nil, "用户未登录或身份已过期", err
	}
	if session.Expiry.Before(time.Now().UTC()) {
		return nil, "", fmt.Errorf("用户未登录或身份已过期")
	}
	var user models.User
	if err := json.Unmarshal(session.Data, &user); err != nil {
		return nil, "用户未登录或身份已过期", fmt.Errorf("会话数据异常: %s", err)
	}
	return &user, "", nil
}

func (s SessionService) DeleteLoginSession(ctx context.Context, id string) (string, error) {
	session := models.Session{Key: id}
	if err := s.Session(ctx).Delete(&session).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("%s=%s; Path=/;Expires=%s", global.LoginSession, id, time.Now().UTC().Format(global.LoginSessionExpiresFormat)), nil
}
