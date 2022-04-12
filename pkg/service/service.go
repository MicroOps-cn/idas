package service

import (
	"context"
	"fmt"
	"idas/config"
	"idas/pkg/client/gorm"
	"idas/pkg/client/ldap"
	"idas/pkg/service/gormservice"
	"idas/pkg/service/ldapservice"
	"io"

	"idas/pkg/service/models"
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
	DeleteLoginSession(ctx context.Context, session string) (string, error)
	GetLoginSession(ctx context.Context, id string) (*models.User, error)
	OAuthAuthorize(ctx context.Context, responseType, clientId, redirectURI string) (redirect string, err error)
	GetOAuthTokenByAuthorizationCode(ctx context.Context, code, clientId, redirectURI string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error)
	GetOAuthTokenByPassword(ctx context.Context, username string, password string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error)
	GetSessions(ctx context.Context, userId string, current int64, size int64) ([]*models.Session, int64, error)
	DeleteSession(ctx context.Context, id string) (err error)

	UploadFile(ctx context.Context, name, contentType string, f io.Reader) (fileKey string, err error)

	GetUsers(ctx context.Context, storage string, keywords string, status models.UserStatus, appId string, current int64, pageSize int64) (users []*models.User, total int64, err error)
	PatchUsers(ctx context.Context, storage string, patch []map[string]interface{}) (count int64, err error)
	DeleteUsers(ctx context.Context, storage string, id []string) (count int64, err error)
	UpdateUser(ctx context.Context, storage string, user *models.User, updateColumns ...string) (*models.User, error)
	GetUserInfo(ctx context.Context, storage string, id string, username string) (user *models.User, err error)
	CreateUser(ctx context.Context, storage string, user *models.User) (*models.User, error)
	PatchUser(ctx context.Context, storage string, user map[string]interface{}) (*models.User, error)
	DeleteUser(ctx context.Context, storage string, id string) error
	CreateLoginSession(ctx context.Context, username string, password string) (string, error)
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
}

type Set struct {
	userAndAppService UserAndAppServices
	sessionService    SessionService
	commonService     CommonService
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
	return &Set{
		userAndAppService: NewUserAndAppService(ctx),
		sessionService:    NewSessionService(ctx),
		commonService:     NewCommonService(ctx),
	}
}

type UserAndAppService interface {
	baseService

	Name() string
	GetUsers(ctx context.Context, keywords string, status models.UserStatus, appId string, current int64, pageSize int64) (users []*models.User, total int64, err error)
	PatchUsers(ctx context.Context, patch []map[string]interface{}) (count int64, err error)
	DeleteUsers(ctx context.Context, id []string) (count int64, err error)
	UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (*models.User, error)
	GetUserInfo(ctx context.Context, id string, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	PatchUser(ctx context.Context, user map[string]interface{}) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	VerifyPassword(ctx context.Context, username string, password string) (*models.User, error)

	GetApps(ctx context.Context, keywords string, current int64, pageSize int64) (apps []*models.App, total int64, err error)
	PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error)
	DeleteApps(ctx context.Context, id []string) (total int64, err error)
	UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (*models.App, error)
	GetAppInfo(ctx context.Context, id string, name string) (app *models.App, err error)
	CreateApp(ctx context.Context, app *models.App) (*models.App, error)
	PatchApp(ctx context.Context, fields map[string]interface{}) (app *models.App, err error)
	DeleteApp(ctx context.Context, id string) (err error)
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
	var userServices UserAndAppServices
	if len(config.Get().GetStorage().GetUser()) > 0 {
		for _, userStorage := range config.Get().GetStorage().GetUser() {

			if userServices.Include(userStorage.GetName()) {
				panic(any(fmt.Errorf("Failed to init UserService: duplicate datasource: %T ", userStorage.Name)))
			}
			switch userSource := userStorage.GetStorageSource().(type) {
			case *config.Storage_Mysql:
				if client, err := gorm.NewMySQLClient(ctx, userSource.Mysql); err != nil {
					panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
				} else {
					userServices = append(userServices, gormservice.NewUserAndAppService(userStorage.GetName(), client))
				}
			case *config.Storage_Sqlite:
				if client, err := gorm.NewSQLiteClient(ctx, userSource.Sqlite); err != nil {
					panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
				} else {
					userServices = append(userServices, gormservice.NewUserAndAppService(userStorage.GetName(), client))
				}
			case *config.Storage_Ldap:
				if client, err := ldap.NewLdapClient(ctx, userSource.Ldap); err != nil {
					panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
				} else {
					userServices = append(userServices, ldapservice.NewUserAndAppService(userStorage.GetName(), client))
				}
			default:
				panic(any(fmt.Errorf("Failed to init UserService: Unknown datasource: %T ", userSource)))
			}
		}
	}
	return userServices
}
