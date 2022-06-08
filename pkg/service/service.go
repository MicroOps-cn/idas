package service

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"idas/config"
	"idas/pkg/client/email"
	"idas/pkg/global"
	"idas/pkg/logs"
	"idas/pkg/service/gormservice"
	"idas/pkg/service/ldapservice"
	"idas/pkg/service/models"
	"io"
)

type migrator interface {
	AutoMigrate(ctx context.Context) error
}

type baseService interface {
	migrator
}

type Service interface {
	baseService

	SetLoginSession(ctx context.Context, user *models.User) (string, error)
	DeleteLoginSession(ctx context.Context, session string) error
	GetLoginSession(ctx context.Context, ids []string) ([]*models.User, error)
	GetAuthCodeByClientId(ctx context.Context, clientId, userId, sessionId, storage string) (code string, err error)
	GetOAuthTokenByAuthorizationCode(ctx context.Context, code, clientId string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error)
	GetOAuthTokenByPassword(ctx context.Context, username string, password string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error)
	GetSessions(ctx context.Context, userId string, current int64, size int64) ([]*models.Token, int64, error)
	DeleteSession(ctx context.Context, id string) (err error)

	UploadFile(ctx context.Context, name, contentType string, f io.Reader) (fileKey string, err error)

	GetUsers(ctx context.Context, storage string, keywords string, status models.UserStatus, appId string, current int64, pageSize int64) (users []*models.User, total int64, err error)
	PatchUsers(ctx context.Context, storage string, patch []map[string]interface{}) (count int64, err error)
	DeleteUsers(ctx context.Context, storage string, id []string) (count int64, err error)
	UpdateUser(ctx context.Context, storage string, user *models.User, updateColumns ...string) (*models.User, error)
	GetUserInfo(ctx context.Context, storage string, id string, username string) (user *models.User, err error)
	GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (user []*models.User)
	CreateUser(ctx context.Context, storage string, user *models.User) (*models.User, error)
	PatchUser(ctx context.Context, storage string, user map[string]interface{}) (*models.User, error)
	DeleteUser(ctx context.Context, storage string, id string) error
	CreateLoginSession(ctx context.Context, username string, password string) ([]string, error)
	GetUserSource(ctx context.Context) (data map[string]string, total int64, err error)

	GetApps(ctx context.Context, storage string, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error)
	GetAppSource(ctx context.Context) (data map[string]string, total int64, err error)
	PatchApps(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, err error)
	DeleteApps(ctx context.Context, storage string, id []string) (total int64, err error)
	UpdateApp(ctx context.Context, storage string, app *models.App, updateColumns ...string) (*models.App, error)
	GetAppInfo(ctx context.Context, storage string, id string) (app *models.App, err error)
	CreateApp(ctx context.Context, storage string, app *models.App) (*models.App, error)
	PatchApp(ctx context.Context, storage string, fields map[string]interface{}) (app *models.App, err error)
	DeleteApp(ctx context.Context, storage string, id string) (err error)
	DownloadFile(ctx context.Context, id string) (f io.ReadCloser, mimiType, fileName string, err error)
	ResetPassword(ctx context.Context, id string, storage string, password string) error

	CreateToken(ctx context.Context, data interface{}, tokenType models.TokenType) (token string, err error)
	VerifyToken(ctx context.Context, token string, relationId string, tokenType models.TokenType) bool
	SendResetPasswordLink(ctx context.Context, token string)
}

type Set struct {
	userAndAppService UserAndAppServices
	sessionService    SessionService
	commonService     CommonService
	smtpClient        *email.SmtpClient
}

func (s Set) GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (users []*models.User) {
	for _, service := range s.userAndAppService {
		if info, err := service.GetUserInfoByUsernameAndEmail(ctx, username, email); err == nil {
			users = append(users, info)
		}
	}
	return users
}

func (s Set) SendResetPasswordLink(ctx context.Context, token string) {
	s.commonService.SendResetPasswordLink(ctx, token)
}

func (s Set) CreateToken(ctx context.Context, data interface{}, tokenType models.TokenType) (token string, err error) {
	tk, err := models.NewToken(tokenType, data)
	if err != nil {
		return "", err
	}
	err = s.sessionService.CreateToken(ctx, tk)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s Set) ResetPassword(ctx context.Context, id string, storage string, password string) error {
	return s.GetUserAndAppService(storage).ResetPassword(ctx, id, password)
}

func (s Set) VerifyToken(ctx context.Context, token, relationId string, tokenType models.TokenType) bool {
	return s.sessionService.VerifyToken(ctx, token, relationId, tokenType)
}

func (s Set) AutoMigrate(ctx context.Context) error {
	svcs := []baseService{
		s.commonService, s.sessionService,
	}
	for _, svc := range s.userAndAppService {
		svcs = append(svcs, svc)
	}
	for _, svc := range svcs {
		if err := svc.AutoMigrate(ctx); err != nil {
			return err
		}
	}
	return nil
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(ctx context.Context) Service {
	client, err := email.NewSmtpClient(ctx, config.Get().GetSmtp())
	if err != nil {
		panic(fmt.Errorf("Failed to init SMTP Client: %s. ", err))
	}
	return &Set{
		userAndAppService: NewUserAndAppService(ctx),
		sessionService:    NewSessionService(ctx),
		commonService:     NewCommonService(ctx),
		smtpClient:        client,
	}
}

type UserAndAppService interface {
	baseService

	Name() string
	GetUsers(ctx context.Context, keywords string, status models.UserStatus, appId string, current int64, pageSize int64) (users []*models.User, total int64, err error)
	PatchUsers(ctx context.Context, patch []map[string]interface{}) (count int64, err error)
	DeleteUsers(ctx context.Context, id []string) (count int64, err error)
	UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (*models.User, error)
	UpdateLoginTime(ctx context.Context, id string) error
	GetUserInfo(ctx context.Context, id string, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	PatchUser(ctx context.Context, user map[string]interface{}) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	VerifyPassword(ctx context.Context, username string, password string) []*models.User

	GetApps(ctx context.Context, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error)
	PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error)
	DeleteApps(ctx context.Context, id []string) (total int64, err error)
	UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (*models.App, error)
	GetAppInfo(ctx context.Context, id string, name string) (app *models.App, err error)
	CreateApp(ctx context.Context, app *models.App) (*models.App, error)
	PatchApp(ctx context.Context, fields map[string]interface{}) (app *models.App, err error)
	DeleteApp(ctx context.Context, id string) (err error)

	VerifyUserAuthorizationForApp(ctx context.Context, appId string, userId string) (scope string, err error)
	ResetPassword(ctx context.Context, id string, password string) error
	GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (*models.User, error)
}

type UserAndAppServices []UserAndAppService

func (s UserAndAppServices) Include(name string) bool {
	for _, service := range s {
		if service.Name() == name {
			return true
		}
	}
	return false
}

func NewUserAndAppService(ctx context.Context) UserAndAppServices {
	logger := log.With(logs.GetContextLogger(ctx), "service", "userAndApp")
	ctx = context.WithValue(ctx, global.LoggerName, logger)
	var userServices UserAndAppServices
	if len(config.Get().GetStorage().GetUser()) > 0 {
		for _, userStorage := range config.Get().GetStorage().GetUser() {
			if userServices.Include(userStorage.GetName()) {
				panic(any(fmt.Errorf("Failed to init UserAndAppService: duplicate datasource: %T ", userStorage.Name)))
			}
			switch userSource := userStorage.GetStorageSource().(type) {
			case *config.Storage_Mysql:
				userServices = append(userServices, gormservice.NewUserAndAppService(userStorage.GetName(), userSource.Mysql.Client))
			case *config.Storage_Sqlite:
				userServices = append(userServices, gormservice.NewUserAndAppService(userStorage.GetName(), userSource.Sqlite.Client))
			case *config.Storage_Ldap:
				userServices = append(userServices, ldapservice.NewUserAndAppService(userStorage.GetName(), userSource.Ldap))
			default:
				panic(any(fmt.Errorf("Failed to init UserAndAppService: Unknown datasource: %T ", userSource)))
			}
		}
	}
	return userServices
}
