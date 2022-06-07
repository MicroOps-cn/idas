package service

import (
	"context"
	"fmt"
	"github.com/go-kit/log"
	"idas/pkg/client/gorm"
	"idas/pkg/global"
	"idas/pkg/service/gormservice"
	"io"
	"os"
	"path"
	"time"

	"github.com/go-kit/log/level"

	"idas/config"
	"idas/pkg/errors"
	"idas/pkg/logs"
)

type CommonService interface {
	baseService
	RecordUploadFile(ctx context.Context, name string, path string, contentType string, size int64) (id string, err error)
	GetFileInfoFromId(ctx context.Context, id string) (fileName, mimiType, filePath string, err error)
	SendResetPasswordLink(ctx context.Context, token string)
}

func NewCommonService(ctx context.Context) CommonService {
	logger := log.With(logs.GetContextLogger(ctx), "service", "common")
	ctx = context.WithValue(ctx, global.LoggerName, logger)
	var commonService CommonService
	commonStorage := config.Get().GetStorage().GetDefault()
	switch commonSource := commonStorage.GetStorageSource().(type) {
	case *config.Storage_Mysql:
		if client, err := gorm.NewMySQLClient(ctx, commonSource.Mysql); err != nil {
			panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
		} else {
			commonService = gormservice.NewCommonService(commonStorage.Name, client)
		}
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
