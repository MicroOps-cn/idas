package gormservice

import (
	"context"
	"fmt"
	gogorm "gorm.io/gorm"
	"idas/pkg/client/gorm"
	"idas/pkg/utils/sign"

	"idas/pkg/service/models"
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
			UserId: userId,
			Key:    pub1,
			Key2:   pub2,
			Secret: privateKey,
		}, nil
	}
}

const authByKey = `
SELECT 
    t_user_key.id,
    t_user_key.user_id,
    t_user_key.key,
    t_user_key.secret,
    T1.id AS User__id,
    T1.create_time AS User__create_time,
    T1.update_time AS User__update_time,
    T1.is_delete AS User__is_delete,
    T1.username AS User__username,
    T1.email AS User__email,
    T1.full_name AS User__full_name,
    T1.user_type AS User__user_type,
    T1.last_login AS User__last_login,
    T1.is_admin AS User__is_admin
FROM
    t_user_key
        JOIN
    t_user T1 ON T1.id = t_user_key.user_id
WHERE
    t_user_key.key = ?
`

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
