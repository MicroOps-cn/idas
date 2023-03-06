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

package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	gohttp "net/http"
	"net/url"
	"time"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/log/level"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/client/email"
	"github.com/MicroOps-cn/idas/pkg/client/http"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service/gormservice"
	"github.com/MicroOps-cn/idas/pkg/service/ldapservice"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
	"github.com/MicroOps-cn/idas/pkg/utils/sign"
)

type migrator interface {
	AutoMigrate(ctx context.Context) error
}

type baseService interface {
	migrator
}

type Service interface {
	baseService
	InitData(ctx context.Context) error
	DeleteLoginSession(ctx context.Context, session string) error
	GetSessionByToken(ctx context.Context, id string, tokenType models.TokenType, receiver interface{}) error
	VerifyPassword(ctx context.Context, username string, password string) (users models.Users, err error)
	GetAuthCodeByAppId(ctx context.Context, clientId string, userId *models.User, sessionId, storage string) (code string, err error)
	GetOAuthTokenByAuthorizationCode(ctx context.Context, code, clientId string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByAuthorizationCode(ctx context.Context, token, clientId, clientSecret string) (accessToken, refreshToken string, expiresIn int, err error)
	RefreshOAuthTokenByPassword(ctx context.Context, token, username, password string) (accessToken, refreshToken string, expiresIn int, err error)
	GetSessions(ctx context.Context, userId string, current int64, size int64) (int64, []*models.Token, error)
	DeleteToken(ctx context.Context, tokenType models.TokenType, id string) (err error)

	UploadFile(ctx context.Context, name, contentType string, f io.Reader) (fileKey string, err error)
	DownloadFile(ctx context.Context, id string) (f io.ReadCloser, mimiType, fileName string, err error)

	CreateRole(ctx context.Context, role *models.Role) (err error)
	UpdateRole(ctx context.Context, role *models.Role) (err error)
	GetRoles(ctx context.Context, keywords string, current, pageSize int64) (count int64, roles []*models.Role, err error)
	GetPermissions(ctx context.Context, keywords string, current int64, pageSize int64) (count int64, permissions []*models.Permission, err error)
	DeleteRoles(ctx context.Context, ids []string) error

	GetUsers(ctx context.Context, storage string, keywords string, status models.UserMeta_UserStatus, appId string, current, pageSize int64) (total int64, users models.Users, err error)
	PatchUsers(ctx context.Context, storage string, patch []map[string]interface{}) (count int64, err error)
	DeleteUsers(ctx context.Context, storage string, id []string) (count int64, err error)
	UpdateUser(ctx context.Context, storage string, user *models.User, updateColumns ...string) (err error)
	GetUserInfo(ctx context.Context, storage, id, username string) (user *models.User, err error)
	GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (user models.Users)
	CreateUser(ctx context.Context, storage string, user *models.User) (err error)
	PatchUser(ctx context.Context, storage string, user map[string]interface{}) (err error)
	DeleteUser(ctx context.Context, storage, id string) error
	GetUserSource(ctx context.Context) (total int64, data map[string]string, err error)
	Authentication(ctx context.Context, method models.AuthMeta_Method, algorithm sign.AuthAlgorithm, key, secret, payload, signStr string) (models.Users, error)
	CreateUserKey(ctx context.Context, userId, name string) (keyPair *models.UserKey, err error)
	GetUserKeys(ctx context.Context, userId string, current, pageSize int64) (count int64, keyPairs []*models.UserKey, err error)
	DeleteUserKey(ctx context.Context, userId string, id string) error

	GetApps(ctx context.Context, storage string, keywords string, filter map[string]interface{}, current, pageSize int64) (total int64, apps []*models.App, err error)
	GetAppSource(ctx context.Context) (total int64, data map[string]string, err error)
	PatchApps(ctx context.Context, storage string, patch []map[string]interface{}) (total int64, err error)
	DeleteApps(ctx context.Context, storage string, id []string) (total int64, err error)
	UpdateApp(ctx context.Context, storage string, app *models.App, updateColumns ...string) (err error)
	GetAppInfo(ctx context.Context, storage string, options ...opts.WithGetAppOptions) (app *models.App, err error)
	GetAppRoleByUserId(ctx context.Context, appId string, userId string) (role *models.AppRole, err error)

	CreateApp(ctx context.Context, storage string, app *models.App) (err error)
	PatchApp(ctx context.Context, storage string, fields map[string]interface{}) (err error)
	DeleteApp(ctx context.Context, storage, id string) (err error)
	ResetPassword(ctx context.Context, id, storage, password string) error
	VerifyPasswordById(ctx context.Context, storage, userId, password string) (users models.Users)

	CreateToken(ctx context.Context, tokenType models.TokenType, data ...interface{}) (token *models.Token, err error)
	VerifyToken(ctx context.Context, token string, tokenType models.TokenType, receiver interface{}, relationId ...string) bool
	SendEmail(ctx context.Context, data map[string]interface{}, topic string, to ...string) error
	Authorization(ctx context.Context, users models.Users, method string) bool
	RegisterPermission(ctx context.Context, permissions models.Permissions) error
	GetProxyConfig(ctx context.Context, host string) (*models.AppProxyConfig, error)
	SendProxyRequest(ctx context.Context, r *gohttp.Request, proxyConfig *models.AppProxyConfig) (*gohttp.Response, error)
	AppAuthentication(ctx context.Context, username string, password string) (*models.App, error)
	CreateAppKey(ctx context.Context, appId string, name string) (appKey *models.AppKey, err error)
	GetAppKeys(ctx context.Context, appId string, current int64, pageSize int64) (count int64, keys []*models.AppKey, err error)
	DeleteAppKey(ctx context.Context, appId string, id []string) (affected int64, err error)
	GetAppKeyFromKey(ctx context.Context, key string) (appKey *models.AppKey, err error)

	DeletePages(ctx context.Context, strings []string) error
	UpdatePage(ctx context.Context, role *models.PageConfig) error
	CreatePage(ctx context.Context, page *models.PageConfig) error
	GetPages(ctx context.Context, filter map[string]interface{}, keywords string, current int64, size int64) (int64, []*models.PageConfig, error)
	GetPage(ctx context.Context, id string) (*models.PageConfig, error)
	PatchPages(ctx context.Context, pages []map[string]interface{}) error

	PatchPageDatas(ctx context.Context, patch []models.PageData) error
	UpdatePageData(ctx context.Context, pageId string, id string, data *json.RawMessage) error
	CreatePageData(ctx context.Context, pageId string, data *json.RawMessage) error
	GetPageData(ctx context.Context, pageId string, id string) (*models.PageData, error)
	GetPageDatas(ctx context.Context, filters map[string]string, keywords string, current int64, size int64) (int64, []*models.PageData, error)
	CreateTOTP(ctx context.Context, ids []string, secret string) error
	GetTOTPSecrets(ctx context.Context, ids []string) ([]string, error)
}

type Set struct {
	userAndAppService UserAndAppServices
	sessionService    SessionService
	commonService     CommonService
}

func (s Set) SendProxyRequest(ctx context.Context, r *gohttp.Request, proxyConfig *models.AppProxyConfig) (*gohttp.Response, error) {
	u, err := url.Parse(proxyConfig.Upstream)
	if err != nil {
		return nil, errors.NewServerError(500, fmt.Sprintf("system error: failed to parse upstream domain: %s", proxyConfig.Upstream))
	}
	r.URL.Scheme = u.Scheme
	r.URL.Host = u.Host

	req, err := gohttp.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		return nil, errors.NewServerError(500, fmt.Sprintf("system error: failed to make request: %s", err))
	}
	req = req.WithContext(ctx)
	if proxyConfig.TransparentServerName {
		req = http.WithTransparentServerName(req, proxyConfig.Domain)
	}
	if proxyConfig.InsecureSkipVerify {
		req = http.WithInsecureSkipVerify(req)
	}
	req.Header = r.Header.Clone()
	req.Header.Del(global.LoginSession)
	return http.SendProxyRequest(req)
}

func (s Set) GetProxyConfig(ctx context.Context, host string) (*models.AppProxyConfig, error) {
	return s.commonService.GetProxyConfig(ctx, host)
}

func (s Set) GetUserAndAppService(name string) UserAndAppService {
	for _, svc := range s.userAndAppService {
		if svc.Name() == name {
			return svc
		}
	}
	return newNullService("", name)
}

func (s Set) SendEmail(ctx context.Context, data map[string]interface{}, topic string, to ...string) error {
	if len(to) == 0 {
		level.Error(logs.GetContextLogger(ctx)).Log("err")
		return errors.ParameterError("recipient is empty")
	}
	smtpConfig := config.Get().Smtp
	if smtpConfig == nil {
		return errors.NewServerError(500, "failed to get smtp options")
	}
	subject, body, err := smtpConfig.GetSubjectAndBody(data, topic)
	if err != nil {
		return errors.WithServerError(500, err, fmt.Sprintf("failed to get email body: topic=%s", topic))
	}
	client, err := email.NewSMTPClient(ctx, smtpConfig)
	if err != nil {
		level.Error(logs.GetContextLogger(ctx)).Log("err", fmt.Sprintf("failed to create SMTP client: %s", err))
		return errors.WithServerError(500, err, "failed to create SMTP client")
	}
	client.SetSubject(subject)
	client.SetBody("text/html", body)
	client.SetTo(to)
	return client.Send()
}

func (s Set) CreateToken(ctx context.Context, tokenType models.TokenType, data ...interface{}) (token *models.Token, err error) {
	tk, err := models.NewToken(tokenType, data...)
	if err != nil {
		return nil, err
	}
	err = s.sessionService.CreateToken(ctx, tk)
	if err != nil {
		return nil, err
	}
	return tk, nil
}

func (s Set) ResetPassword(ctx context.Context, id string, storage string, password string) error {
	if len(storage) == 0 {
		for _, service := range s.userAndAppService {
			if err := service.ResetPassword(ctx, id, password); err != nil {
				level.Error(logs.GetContextLogger(ctx)).Log("err", err, "msg", "Failed to reset password", "userId", id)
			}
		}
		return nil
	}
	return s.GetUserAndAppService(storage).ResetPassword(ctx, id, password)
}

func (s Set) VerifyToken(ctx context.Context, token string, tokenType models.TokenType, receiver interface{}, relationId ...string) bool {
	logger := logs.GetContextLogger(ctx)
	tk, err := s.sessionService.GetToken(ctx, token, tokenType, relationId...)
	if err != nil {
		return false
	} else if tk == nil {
		return false
	}
	switch tokenType {
	case models.TokenTypeResetPassword,
		models.TokenTypeCode,
		models.TokenTypeAppProxyLogin,
		models.TokenTypeActive:
		if err := s.DeleteToken(ctx, tokenType, tk.Id); err != nil {
			level.Warn(logger).Log("msg", "failed to delete token.", "err", err)
		}
	}

	if !tk.Expiry.After(time.Now().UTC()) {
		return false
	}
	if receiver != nil {
		if err = tk.To(receiver); err != nil {
			level.Warn(logger).Log("msg", "failed to parse token data.", "err", err)
			return false
		}
	}
	return true
}

func (s Set) InitData(ctx context.Context) error {
	for _, svc := range s.userAndAppService {
		adminUser, err := svc.GetUserInfo(ctx, "", "admin")
		if errors.IsNotFount(err) {
			adminUser = &models.User{
				Username: "admin",
				Password: sql.RawBytes("idas"),
			}
			err = svc.CreateUser(ctx, adminUser)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		idasApp, err := svc.GetAppInfo(ctx, opts.WithAppName(global.IdasAppName))
		if errors.IsNotFount(err) {
			idasApp = &models.App{
				Name:        global.IdasAppName,
				Description: "Identity authentication service. It is bound to the current service. Please do not delete it at will.",
				GrantMode:   models.AppMeta_manual,
				Roles: models.AppRoles{{
					Name: "admin",
				}, {
					Name:      "viewer",
					IsDefault: true,
				}},
				Users: []*models.User{{
					Model: models.Model{Id: adminUser.Id},
					Role:  "admin",
				}},
			}
			if err = svc.CreateApp(ctx, idasApp); err != nil {
				return errors.WithMessage(err, "failed to initialize application data")
			}
		} else if err != nil {
			return errors.WithMessage(err, "failed to initialize application data")
		}
		if len(idasApp.Roles) == 0 {
			if len(idasApp.Roles) == 0 {
				idasApp.Roles = models.AppRoles{{
					Name: "admin",
				}, {
					Name:      "viewer",
					IsDefault: true,
				}}
			}
			if err = svc.UpdateApp(ctx, idasApp); err != nil {
				return err
			}
		}
		if len(idasApp.Users) == 0 {
			idasApp.Users = []*models.User{{
				Model:  models.Model{Id: adminUser.Id},
				RoleId: idasApp.Roles.GetRole("admin").GetId(),
			}}
			if err = svc.UpdateApp(ctx, idasApp); err != nil {
				return err
			}
		}
	}
	return nil
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
	GetUsers(ctx context.Context, keywords string, status models.UserMeta_UserStatus, appId string, current, pageSize int64) (total int64, users []*models.User, err error)
	PatchUsers(ctx context.Context, patch []map[string]interface{}) (count int64, err error)
	DeleteUsers(ctx context.Context, id []string) (count int64, err error)
	UpdateUser(ctx context.Context, user *models.User, updateColumns ...string) (err error)
	UpdateLoginTime(ctx context.Context, id string) error
	GetUserInfo(ctx context.Context, id string, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (err error)
	PatchUser(ctx context.Context, user map[string]interface{}) (err error)
	DeleteUser(ctx context.Context, id string) error
	VerifyPassword(ctx context.Context, username string, password string) []*models.User

	GetApps(ctx context.Context, keywords string, filter map[string]interface{}, current, pageSize int64) (total int64, apps []*models.App, err error)
	PatchApps(ctx context.Context, patch []map[string]interface{}) (total int64, err error)
	DeleteApps(ctx context.Context, id []string) (total int64, err error)
	UpdateApp(ctx context.Context, app *models.App, updateColumns ...string) (err error)
	GetAppInfo(ctx context.Context, o ...opts.WithGetAppOptions) (app *models.App, err error)
	CreateApp(ctx context.Context, app *models.App) (err error)
	PatchApp(ctx context.Context, fields map[string]interface{}) (err error)
	DeleteApp(ctx context.Context, id string) (err error)

	ResetPassword(ctx context.Context, id string, password string) error
	GetUserInfoByUsernameAndEmail(ctx context.Context, username, email string) (*models.User, error)
	VerifyPasswordById(ctx context.Context, id, password string) (users []*models.User)
	GetUsersById(ctx context.Context, id []string) (models.Users, error)
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
	ctx, _ = logs.NewContextLogger(ctx, logs.WithKeyValues("service", "userAndApp"))
	var userServices UserAndAppServices
	if len(config.Get().GetStorage().GetUser()) > 0 {
		for _, userStorage := range config.Get().GetStorage().GetUser() {
			if userServices.Include(userStorage.GetName()) {
				panic(fmt.Sprintf("Failed to init UserAndAppService: duplicate datasource: %T ", userStorage.Name))
			}
			switch userSource := userStorage.GetStorageSource().(type) {
			case *config.Storage_Mysql:
				userServices = append(userServices, gormservice.NewUserAndAppService(ctx, userStorage.GetName(), userSource.Mysql.Client))
			case *config.Storage_Sqlite:
				userServices = append(userServices, gormservice.NewUserAndAppService(ctx, userStorage.GetName(), userSource.Sqlite.Client))
			case *config.Storage_Ldap:
				userServices = append(userServices, ldapservice.NewUserAndAppService(ctx, userStorage.GetName(), userSource.Ldap))
			default:
				panic(fmt.Sprintf("Failed to init UserAndAppService: Unknown datasource: %T ", userSource))
			}
		}
	}
	return userServices
}

func (s Set) CreateTOTP(ctx context.Context, ids []string, secret string) error {
	return s.commonService.CreateTOTP(ctx, ids, secret)
}

func (s Set) GetTOTPSecrets(ctx context.Context, ids []string) ([]string, error) {
	return s.commonService.GetTOTPSecrets(ctx, ids)
}
