package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/go-kit/log/level"

	"idas/config"
	"idas/pkg/errors"
	"idas/pkg/logs"
	"idas/pkg/service/gormservice"
	"idas/pkg/service/models"
	"idas/pkg/utils/sets"
)

type CommonService interface {
	baseService
	RecordUploadFile(ctx context.Context, name string, path string, contentType string, size int64) (id string, err error)
	GetFileInfoFromId(ctx context.Context, id string) (fileName, mimiType, filePath string, err error)
	CreateRole(ctx context.Context, role *models.Role) (newRole *models.Role, err error)
	UpdateRole(ctx context.Context, role *models.Role) (newRole *models.Role, err error)
	GetRoles(ctx context.Context, keywords string, current, pageSize int64) (count int64, roles []*models.Role, err error)
	GetPermissions(ctx context.Context, keywords string, current int64, pageSize int64) (count int64, permissions []*models.Permission, err error)
	DeleteRoles(ctx context.Context, ids []string) error
	RegisterPermission(ctx context.Context, permissions models.Permissions) error
	CreateOrUpdateRoleByName(ctx context.Context, role *models.Role) error
	Authorization(ctx context.Context, roles []string, method string) bool

	GetUserKey(ctx context.Context, key string) (*models.UserKey, error)
	CreateUserKeyWithId(ctx context.Context, userId string, name string) (userKey *models.UserKey, err error)
}

func NewCommonService(ctx context.Context) CommonService {
	// logger := log.With(logs.GetContextLogger(ctx), "service", "common")
	// ctx = context.WithValue(ctx, global.LoggerName, logger)
	var commonService CommonService
	commonStorage := config.Get().GetStorage().GetDefault()
	switch commonSource := commonStorage.GetStorageSource().(type) {
	case *config.Storage_Mysql:
		commonService = gormservice.NewCommonService(commonStorage.Name, commonSource.Mysql.Client)
	case *config.Storage_Sqlite:
		commonService = gormservice.NewCommonService(commonStorage.Name, commonSource.Sqlite.Client)
	default:
		panic(any(fmt.Errorf("初始化CommonService失败: 未知的数据源类型: %T", commonSource)))
	}
	return commonService
}

func (s Set) UploadFile(ctx context.Context, name, contentType string, f io.Reader) (fileKey string, err error) {
	logger := logs.GetContextLogger(ctx)
	now := time.Now().UTC()
	d, err := config.Get().GetUploadDir()
	if err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to get upload dir")
		return "", errors.InternalServerError
	}
	dirName := now.Format("2006-01")
	if _, err = d.Stat(dirName); os.IsNotExist(err) {
		//nolint:gofumpt
		if err = d.MkdirAll(dirName, 0755); err != nil {
			level.Error(logger).Log("msg", "failed to create directory", "err", err)
		}
	} else if err != nil {
		level.Error(logger).Log("msg", "failed to get directory status", "err", err)
	}
	filePath := fmt.Sprintf("%s/%d%s", dirName, now.UnixNano(), path.Ext(name))

	var ff io.ReadWriteCloser
	//nolint:gofumpt
	if ff, err = d.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to open file", "filePath", filePath)
		return "", errors.InternalServerError
	}
	defer ff.Close()
	size, err := io.Copy(ff, f)
	if err != nil {
		return "", err
	}
	return s.commonService.RecordUploadFile(ctx, name, filePath, contentType, size)
}

func (s Set) DownloadFile(ctx context.Context, id string) (f io.ReadCloser, mimiType, fileName string, err error) {
	var filePath string
	fileName, mimiType, filePath, err = s.commonService.GetFileInfoFromId(ctx, id)
	if err != nil {
		return nil, "", "", err
	}
	logger := logs.GetContextLogger(ctx)
	d, err := config.Get().GetUploadDir()
	if err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to get upload dir")
		return nil, "", "", errors.InternalServerError
	}
	if f, err = d.Open(filePath); err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to open file", "filePath", filePath)
		return nil, "", "", errors.InternalServerError
	}
	return f, mimiType, fileName, nil
}

func (s Set) CreateRole(ctx context.Context, role *models.Role) (newRole *models.Role, err error) {
	return s.commonService.CreateRole(ctx, role)
}

func (s Set) UpdateRole(ctx context.Context, role *models.Role) (newRole *models.Role, err error) {
	return s.commonService.UpdateRole(ctx, role)
}

func (s Set) GetRoles(ctx context.Context, keywords string, current, pageSize int64) (count int64, roles []*models.Role, err error) {
	return s.commonService.GetRoles(ctx, keywords, current, pageSize)
}

func (s Set) GetPermissions(ctx context.Context, keywords string, current int64, pageSize int64) (count int64, permissions []*models.Permission, err error) {
	return s.commonService.GetPermissions(ctx, keywords, current, pageSize)
}

func (s Set) DeleteRoles(ctx context.Context, ids []string) error {
	return s.commonService.DeleteRoles(ctx, ids)
}

func (s Set) Authorization(ctx context.Context, users []*models.User, method string) bool {
	roles := sets.New[string]()
	for _, user := range users {
		roles.Insert(user.Role)
	}
	return s.commonService.Authorization(ctx, roles.List(), method)
}

func (s Set) RegisterPermission(ctx context.Context, permissions models.Permissions) error {
	err := s.commonService.RegisterPermission(ctx, permissions)
	if err != nil {
		return err
	}

	for _, role := range permissions.GetRoles() {
		if err = s.commonService.CreateOrUpdateRoleByName(ctx, role); err != nil {
			return err
		}
	}
	return nil
}
