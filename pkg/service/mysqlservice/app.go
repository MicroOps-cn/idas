package mysqlservice

import (
	"context"

	"idas/pkg/client/mysql"
	"idas/pkg/service/models"
)

type AppService struct {
	*mysql.Client
	name string
}

func NewAppService(name string, client *mysql.Client) *AppService {
	if err := client.Session(context.Background()).SetupJoinTable(&models.App{}, "User", models.AppUser{}); err != nil {
		panic(err)
	}
	return &AppService{name: name, Client: client}
}
