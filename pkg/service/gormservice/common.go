package gormservice

import (
	"context"
	"fmt"

	gogorm "gorm.io/gorm"

	"idas/pkg/client/gorm"
	"idas/pkg/service/models"
	"idas/pkg/utils/sign"
)

func NewCommonService(name string, client *gorm.Client) *CommonService {
	return &CommonService{name: name, Client: client}
}

type CommonService struct {
	*gorm.Client
	name string
}

func (c CommonService) Name() string {
	return c.name
}

func (c CommonService) AutoMigrate(ctx context.Context) error {
	err := c.Session(ctx).AutoMigrate(&models.File{}, &models.Permission{}, &models.Role{})
	if err != nil {
		return err
	}
	return nil
}

func (c CommonService) RecordUploadFile(ctx context.Context, name string, path string, contentType string, size int64) (id string, err error) {
	file := &models.File{MimiType: contentType, Name: name, Path: path, Size: size}
	if err = c.Session(ctx).Create(file).Error; err != nil {
		return
	}
	return file.Id, err
}

func (c CommonService) GetFileInfoFromId(ctx context.Context, id string) (fileName, mimiType, filePath string, err error) {
	file := &models.File{Model: models.Model{Id: id}}
	if err = c.Session(ctx).First(file).Error; err != nil {
		return "", "", "", err
	}
	return file.Name, file.MimiType, file.Path, nil
}

func (c CommonService) CreateUserKeyWithId(ctx context.Context, userId string, name string) (userKey *models.UserKey, err error) {
	conn := c.Session(ctx)
	if err != nil {
		return nil, err
		//} else if id := strings.ReplaceAll(uuid.NewV4().String(), "-", ""); len(id) != 32 {
		//	return nil, errors.New("生成ID失败: " + id)
	} else if pub1, pub2, privateKey, err := sign.GenerateECDSAKeyPair(); err != nil {
		return nil, err
	} else {
		userKey = &models.UserKey{
			Name:   name,
			UserId: userId,
			Key:    pub1,
			Secret: pub2,
		}
		if err = conn.Create(&userKey).Error; err != nil {
			return nil, err
		}
		return &models.UserKey{
			Model:   userKey.Model,
			UserId:  userId,
			Key:     pub1,
			Secret:  pub2,
			Private: privateKey,
		}, nil
	}
}

func (c CommonService) GetUserKey(ctx context.Context, key string) (*models.UserKey, error) {
	userKey := &models.UserKey{Key: key}
	conn := c.Session(ctx)
	if err := conn.Where("key = ?", key).First(&userKey).Error; err != nil && err != gogorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to query user key, database exception: %s", err)
	} else if err != nil {
		return nil, nil
	}
	return userKey, nil
}
