package service

import (
	"context"
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

	GetUsers(ctx context.Context, storage string, keyword string, status models.UserStatus, current int64, pageSize int64) (users []*models.User, total int64, err error)
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
	userService    UserServices
	appService     AppServices
	sessionService SessionService
	commonService  CommonService
}

func (s Set) AutoMigrate(ctx context.Context) error {
	svcs := []baseService{
		s.commonService, s.sessionService,
	}
	for _, svc := range s.userService {
		svcs = append(svcs, svc)
	}
	for _, svc := range s.appService {
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
		userService:    NewUserServices(ctx),
		sessionService: NewSessionService(ctx),
		appService:     NewAppService(ctx),
		commonService:  NewCommonService(ctx),
	}
}
