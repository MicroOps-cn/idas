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
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/sets"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/rand"
	"strconv"
	"strings"
	"testing"
)

func testAppService(t *testing.T, ctx context.Context, storage string, svc Service) {
	var userIds []string
	var getRandomUsers = func(roles models.AppRoles) []*models.User {
		rCount := rand.Intn(20)
		users := make([]*models.User, rCount)

		for i := 0; i < rCount; i++ {
			users[i] = &models.User{Model: models.Model{Id: userIds[rand.Intn(len(userIds))]}}
			if len(roles) > 0 {
				if rand.Intn(3) > 1 {
					users[i].Role = roles[rand.Intn(len(roles))].Name
				}
				if rand.Intn(3) > 0 {
					users[i].RoleId = roles[rand.Intn(len(roles))].Id
				}
			}
		}
		return users
	}

	var getRandomRoles = func() models.AppRoles {
		rCount := rand.Intn(20)
		roles := make(models.AppRoles, rCount)
		for i := 0; i < rCount; i++ {
			var id string
			if i%3 == 0 {
				id = models.NewId()
			}
			roles[i] = &models.AppRole{Model: models.Model{Id: id}, Name: rand.String(10)}
		}
		return roles
	}

	if !t.Run("Test Init User Data", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			_, err := svc.CreateUser(ctx, storage, &models.User{
				Username:    rand.String(5),
				Email:       rand.String(7),
				PhoneNumber: rand.String(9),
				FullName:    rand.String(5),
				Avatar:      rand.String(20),
				Status:      models.UserMeta_UserStatus(rand.Intn(4)),
			})
			require.NoError(t, err)
		}
		_, users, err := svc.GetUsers(ctx, storage, "", models.UserMeta_status_all, "", 1, 200)
		require.NoError(t, err)
		require.Truef(t, len(users) >= 100, "Failed to obtain the information of all current users: the number of users (%s) is less than 100.", len(users))
		for _, user := range users {
			userIds = append(userIds, user.Id)
		}
	}) {
		return
	}

	if !t.Run("Test Create App", func(t *testing.T) {
		{
			count, users, err := svc.GetApps(ctx, storage, "", 1, 1024)
			require.NoError(t, err)
			require.Len(t, users, 0)
			require.Equal(t, count, int64(0))
		}

		t.Run("Test Create Null App", func(t *testing.T) {
			_, err := svc.CreateApp(ctx, storage, &models.App{})
			require.Error(t, err)
		})
		t.Run("Test Create duplicate App", func(t *testing.T) {
			_, err := svc.CreateApp(ctx, storage, &models.App{Name: "Test App - AAA - 01"})
			require.NoError(t, err)
			_, err = svc.CreateApp(ctx, storage, &models.App{Name: "Test App - AAA - 01"})
			require.Error(t, err)
		})
		t.Run("Test Create App with empty role name ", func(t *testing.T) {
			_, err := svc.CreateApp(ctx, storage, &models.App{
				Name: rand.String(rand.IntnRange(1, 20)),
				Role: models.AppRoles{{Name: ""}},
			})
			require.EqualError(t, err, "Parameter Error: role name and id is nil")
		})
		t.Run("Test Create App with duplicate role", func(t *testing.T) {
			_, err := svc.CreateApp(ctx, storage, &models.App{
				Name: rand.String(rand.IntnRange(1, 20)),
				Role: models.AppRoles{{Name: "AAA"}, {Name: "AAA"}},
			})
			require.EqualError(t, err, "Parameter Error: duplicate role: AAA")
		})

		t.Run("Test Random Create App", func(t *testing.T) {
			for i := 0; i < 100; i++ {
				roles := getRandomRoles()
				users := getRandomUsers(roles)
				var roleNames = make([]string, len(roles))
				var userRoles = make(map[string]string, len(users))
				for j, role := range roles {
					roleNames[j] = role.Name
				}
				for _, user := range users {
					userRoles[user.Id] = user.Role
				}
				app := &models.App{
					Name:        rand.String(rand.IntnRange(1, 20)) + "_" + strconv.Itoa(i),
					Description: rand.String(rand.Intn(20)),
					Avatar:      rand.String(rand.Intn(20)),
					GrantType:   models.AppMeta_GrantType(rand.Intn(len(models.AppMeta_GrantType_value))),
					GrantMode:   models.AppMeta_GrantMode(rand.Intn(len(models.AppMeta_GrantMode_value))),
					Status:      models.AppMeta_Status(rand.Intn(len(models.AppMeta_Status_value))),
					User:        users,
					Role:        roles,
				}
				app, err := svc.CreateApp(ctx, storage, app)
				require.NoError(t, err)
				info, err := svc.GetAppInfo(ctx, storage, app.Id)
				require.NoError(t, err)
				var roleNames1 = make([]string, len(info.Role))
				var userRoles1 = make(map[string]string, len(info.User))
				for j, role := range roles {
					roleNames1[j] = role.Name
				}
				for _, user := range users {
					userRoles1[user.Id] = user.Role
				}
				require.Equal(t, roleNames, roleNames1)
				require.Equal(t, userRoles, userRoles1)
			}
		})
	}) {
		return
	}

	var appIds []string

	if !t.Run("Test List App", func(t *testing.T) {
		count, apps, err := svc.GetApps(ctx, storage, "", 1, 1024)
		require.NoError(t, err)
		require.Truef(t, len(apps) == 101, "Failed to obtain the information of all current apps: the number of users (%s) is not equal 101.", len(apps))
		require.Truef(t, count == 101, "Failed to obtain the information of all current apps: the number of users (%s) is not equal 101.", len(apps))
		for _, app := range apps {
			appIds = append(appIds, app.Id)
		}
		t.Run("Filter by keywords", func(t *testing.T) {
			_, apps, err = svc.GetApps(ctx, storage, "AAA", 1, 1024)
			require.NoError(t, err)
			require.True(t, len(apps) >= 1)
			for _, app := range apps {
				if !strings.Contains(app.Name, "AAA") && !strings.Contains(app.Description, "AAA") {
					t.Error("The application name and description do not contain `AAA`, but are queried.")
				}
			}
		})
	}) {
		return
	}

	t.Run("Test Update App", func(t *testing.T) {
		appId := appIds[rand.Intn(len(appIds))]
		appInfo, err := svc.GetAppInfo(ctx, storage, appId)
		require.NoError(t, err)

		var roleNames = make([]string, len(appInfo.Role))
		var userRoles = make(map[string]string, len(appInfo.User))
		for j, role := range appInfo.Role {
			roleNames[j] = role.Name
		}
		for _, user := range appInfo.User {
			userRoles[user.Id] = user.Role
		}
		var newAppInfo = *appInfo
		newAppInfo.Name = rand.String(10)
		newAppInfo.Description = rand.String(10)
		newAppInfo.Avatar = rand.String(10)
		newAppInfo.GrantMode = (newAppInfo.GrantMode + 1) % models.AppMeta_GrantMode(len(models.AppMeta_GrantMode_name))
		newAppInfo.GrantType = (newAppInfo.GrantType + 1) % models.AppMeta_GrantType(len(models.AppMeta_GrantType_name))
		newAppInfo.Status = (newAppInfo.Status + 1) % models.AppMeta_Status(len(models.AppMeta_Status_name))
		aUserIds := sets.New[string](userIds...)
		for _, user := range newAppInfo.User {
			if aUserIds.Has(user.Id) {
				aUserIds.Delete(user.Id)
			}
		}
		app := models.AppRole{
			Name: "X_1asl",
		}
		roleNames = append(roleNames, app.Name)
		userId, ok := aUserIds.PopAny()
		if ok {
			newAppInfo.User = append(newAppInfo.User, &models.User{
				Model: models.Model{Id: userId},
				Role:  app.Name,
			})
			userRoles[userId] = app.Name
		}
		newAppInfo.Role = append(newAppInfo.Role, &app)
		tmpAppInfo := newAppInfo
		_, err = svc.UpdateApp(ctx, storage, &tmpAppInfo)
		require.NoError(t, err)

		newAppInfo1, err := svc.GetAppInfo(ctx, storage, appId)
		require.NoError(t, err)
		require.Equal(t, newAppInfo.Name, newAppInfo1.Name)
		require.Equal(t, newAppInfo.Description, newAppInfo1.Description)
		require.Equal(t, newAppInfo.Avatar, newAppInfo1.Avatar)
		require.Equal(t, newAppInfo.GrantType, newAppInfo1.GrantType)
		require.Equal(t, newAppInfo.GrantMode, newAppInfo1.GrantMode)
		require.Equal(t, newAppInfo.Status, newAppInfo1.Status)

		var roleNames1 = make([]string, len(newAppInfo1.Role))
		var userRoles1 = make(map[string]string, len(newAppInfo1.User))
		for j, role := range newAppInfo1.Role {
			roleNames1[j] = role.Name
		}
		for _, user := range newAppInfo1.User {
			userRoles1[user.Id] = user.Role
		}

		require.Equal(t, roleNames, roleNames1)
		require.Equal(t, userRoles, userRoles1)
	})

	t.Run("Test Delete App", func(t *testing.T) {
		err := svc.DeleteApp(ctx, storage, appIds[rand.Intn(len(appIds))])
		require.NoError(t, err)

		t.Run("Test List App", func(t *testing.T) {
			count, apps, err := svc.GetApps(ctx, storage, "", 1, 1024)
			require.NoError(t, err)
			require.Truef(t, len(apps) == 100, "Failed to obtain the information of all current apps: the number of users (%s) is not equal 100.", len(apps))
			require.Truef(t, count == 100, "Failed to obtain the information of all current apps: the number of users (%s) is not equal 100.", len(apps))
		})
	})
}
