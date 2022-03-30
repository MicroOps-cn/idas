package mysqlservice

import (
	"context"

	"idas/pkg/client/mysql"
	"idas/pkg/service/models"
)

type AppService struct {
	*mysql.Client
}

func (a AppService) SetupJoinTable() error {
	return a.Session(context.Background()).SetupJoinTable(&models.App{}, "User", models.AppUser{})
}

func NewAppService(client *mysql.Client) *AppService {
	return &AppService{Client: client}
}
