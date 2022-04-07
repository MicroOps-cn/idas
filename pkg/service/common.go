package service

import (
	"context"
	"fmt"
	"github.com/go-kit/log/level"
	"idas/config"
	"idas/pkg/client/mysql"
	"idas/pkg/errors"
	"idas/pkg/logs"
	"idas/pkg/service/mysqlservice"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"
)

type CommonService interface {
	baseService
	RecordUploadFile(ctx context.Context, name string, path string, contentType string, size int64) (id string, err error)
	GetFileInfoFromId(ctx context.Context, id string) (fileName, mimiType, filePath string, err error)
}

func NewCommonService(ctx context.Context) CommonService {
	var commonService CommonService
	commonStorage := config.Get().GetStorage().GetDefault()
	switch commonSource := commonStorage.GetStorageSource().(type) {
	case *config.Storage_Mysql:
		if client, err := mysql.NewMySQLClient(ctx, commonSource.Mysql); err != nil {
			panic(any(fmt.Errorf("初始化UserService失败: MySQL数据库连接失败: %s", err)))
		} else {
			commonService = mysqlservice.NewCommonService(commonStorage.Name, client)
		}
	default:
		panic(any(fmt.Errorf("初始化CommonService失败: 未知的数据源类型: %T", commonSource)))
	}
	return commonService
}

func (s Set) UploadFile(ctx context.Context, name, contentType string, f multipart.File) (fileKey string, err error) {
	logger := logs.GetContextLogger(ctx)
	now := time.Now().UTC()
	if d, err := config.Get().GetUploadDir(); err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to get upload dir")
		return "", errors.InternalServerError
	} else {
		dirName := now.Format("2006-01")
		if _, err = d.Stat(dirName); os.IsNotExist(err) {
			if err = d.MkdirAll(dirName, 0755); err != nil {
				level.Error(logger).Log("msg", "failed to create directory", "err", err)
			}
		} else if err != nil {
			level.Error(logger).Log("msg", "failed to get directory status", "err", err)
		}
		filePath := fmt.Sprintf("%s/%d%s", dirName, now.UnixNano(), path.Ext(name))
		if ff, err := d.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755); err != nil {
			level.Error(logger).Log("err", err, "msg", "failed to open file", "filePath", filePath)
			return "", errors.InternalServerError
		} else {
			defer ff.Close()
			size, err := io.Copy(ff, f)
			if err != nil {
				return "", err
			}
			return s.commonService.RecordUploadFile(ctx, name, filePath, contentType, size)
		}
	}
}

func (s Set) DownloadFile(ctx context.Context, id string) (f io.ReadCloser, mimiType, fileName string, err error) {
	var filePath string
	fileName, mimiType, filePath, err = s.commonService.GetFileInfoFromId(ctx, id)
	if err != nil {
		return nil, "", "", err
	}
	logger := logs.GetContextLogger(ctx)
	if d, err := config.Get().GetUploadDir(); err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to get upload dir")
		return nil, "", "", errors.InternalServerError
	} else {
		if f, err = d.Open(filePath); err != nil {
			level.Error(logger).Log("err", err, "msg", "failed to open file", "filePath", filePath)
			return nil, "", "", errors.InternalServerError
		}
		return f, mimiType, fileName, nil
	}
}
