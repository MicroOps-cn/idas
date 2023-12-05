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
	g "github.com/MicroOps-cn/fuck/generator"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/MicroOps-cn/fuck/sets"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/service/opts"
)

func testAppService(ctx context.Context, t *testing.T, svc Service) {
	var userIds []string
	getRandomUsers := func(roles models.AppRoles) []*models.User {
		aUserIds := sets.New[string](userIds...)
		rCount := rand.IntnRange(1, 21)
		users := make([]*models.User, rCount)

		for i := 0; i < rCount; i++ {
			userId, ok := aUserIds.PopAny()
			if !ok {
				continue
			}
			users[i] = &models.User{Model: models.Model{Id: userId}}
			if len(roles) > 0 {
				r := rand.Intn(len(roles))
				if rand.Intn(3) > 1 {
					users[i].Role = roles[r].Name
				}
				if rand.Intn(3) > 0 {
					users[i].RoleId = roles[r].Id
				}
			}
		}
		return users
	}
	getRandomUrls := func() models.AppProxyUrls {
		rCount := rand.Intn(20)
		urls := make(models.AppProxyUrls, rCount)
		for i := 0; i < rCount; i++ {
			var id string
			if i%3 == 0 {
				id = g.NewUUID().String()
			}
			urls[i] = &models.AppProxyUrl{Model: models.Model{Id: id}, Name: rand.String(10)}
		}
		return urls
	}
	getRandomRoles := func(urls ...string) models.AppRoles {
		rCount := rand.Intn(20)
		roles := make(models.AppRoles, rCount)
		for i := 0; i < rCount; i++ {
			var id string
			if i%3 == 0 {
				id = g.NewUUID().String()
			}
			roles[i] = &models.AppRole{Model: models.Model{Id: id}, Name: rand.String(10)}
		}
		return roles
	}

	if !t.Run("Test Init User Data", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			err := svc.CreateUser(ctx, &models.User{
				Username:    rand.String(5),
				Email:       rand.String(7),
				PhoneNumber: rand.String(9),
				FullName:    rand.String(5),
				Avatar:      rand.String(20),
				Status:      models.UserMeta_UserStatus(rand.Intn(4)),
			})
			require.NoError(t, err)
		}
		_, users, err := svc.GetUsers(ctx, "", models.UserMetaStatusAll, "", 1, 200)
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
			count, users, err := svc.GetApps(ctx, "", nil, 1, 1024)
			require.NoError(t, err)
			require.Len(t, users, 0)
			require.Equal(t, count, int64(0))
		}

		t.Run("Test Create Null App", func(t *testing.T) {
			err := svc.CreateApp(ctx, &models.App{})
			require.Error(t, err)
		})
		t.Run("Test Create duplicate Name", func(t *testing.T) {
			users := models.Users{{Model: models.Model{Id: userIds[rand.Intn(len(userIds))]}}}
			err := svc.CreateApp(ctx, &models.App{Name: "Test App - AAA - 01", Users: users})
			require.NoError(t, err)
			err = svc.CreateApp(ctx, &models.App{Name: "Test App - AAA - 01", Users: users})
			require.Error(t, err)
		})
		t.Run("Test Create App with empty role name ", func(t *testing.T) {
			err := svc.CreateApp(ctx, &models.App{
				Name:  rand.String(rand.IntnRange(1, 20)),
				Roles: models.AppRoles{{Name: ""}},
			})
			require.EqualError(t, err, "Parameter Error: role name and id is nil")
		})
		t.Run("Test Create App with duplicate role", func(t *testing.T) {
			err := svc.CreateApp(ctx, &models.App{
				Name:  rand.String(rand.IntnRange(1, 20)),
				Roles: models.AppRoles{{Name: "AAA"}, {Name: "AAA"}},
			})
			require.EqualError(t, err, "Parameter Error: duplicate role: AAA")
		})

		t.Run("Test Random Create App", func(t *testing.T) {
			for i := 0; i < 100; i++ {
				urls := getRandomUrls()
				roles := getRandomRoles(urls.Id()...)
				users := getRandomUsers(roles)
				roleNames := make([]string, len(roles))
				userRoles := make(map[string]string, len(users))
				roleURL := make(map[string][]string)
				for j, role := range roles {
					roleNames[j] = role.Name
				}
				defaultRole := -1
				if len(roles) > 0 {
					defaultRole = rand.Intn(len(roles))
					roles[defaultRole].IsDefault = true
				}
				for _, user := range users {
					userRoles[user.Id] = user.Role
					if len(user.RoleId) != 0 {
						for _, role := range roles {
							if role.Id == user.RoleId {
								userRoles[user.Id] = role.Name
							}
						}
					}
					if userRoles[user.Id] == "" && defaultRole >= 0 {
						userRoles[user.Id] = roles[defaultRole].Name
					}
				}

				var proxy [][5]string
				app := &models.App{
					Name:        rand.String(rand.IntnRange(1, 20)) + "_" + strconv.Itoa(i),
					Description: rand.String(rand.Intn(20)),
					Avatar:      rand.String(rand.Intn(20)),
					GrantType:   models.AppMeta_GrantType(rand.Intn(len(models.AppMeta_GrantType_value))),
					GrantMode:   models.AppMeta_GrantMode(rand.Intn(len(models.AppMeta_GrantMode_value))),
					Status:      models.AppMeta_Status(rand.Intn(len(models.AppMeta_Status_value))),
					Users:       users,
					Roles:       roles,
				}
				if app.GrantType&models.AppMeta_proxy == models.AppMeta_proxy {
					app.Proxy = &models.AppProxy{
						Domain:   rand.String(rand.Intn(20)),
						Upstream: rand.String(rand.Intn(20)),
					}
					for j := 0; j < rand.IntnRange(1, 10); j++ {
						app.Proxy.Urls = append(app.Proxy.Urls, &models.AppProxyUrl{
							Model:  models.Model{Id: rand.String(rand.Intn(20))},
							Method: "GET",
							Url:    rand.String(10),
							Name:   rand.String(10),
						})
					}
					for _, url := range app.Proxy.Urls {
						proxy = append(proxy, [5]string{app.Proxy.Domain, app.Proxy.Upstream, url.Method, url.Url, url.Name})
						for j := 0; j < 3; j++ {
							role := app.Roles[rand.Intn(len(app.Roles))]
							role.UrlsId = append(role.UrlsId, url.Id)
						}
					}
				}
				for _, role := range app.Roles {
					if len(role.UrlsId) != 0 {
						roleURL[role.Id] = role.UrlsId
						sort.Strings(roleURL[role.Id])
					}
				}
				err := svc.CreateApp(ctx, app)
				require.NoError(t, err)
				require.NotEmpty(t, app.Id)
				info, err := svc.GetAppInfo(ctx, opts.WithAppId(app.Id))
				require.NoError(t, err)
				roleNames1 := make([]string, len(info.Roles))
				userRoles1 := make(map[string]string, len(info.Users))
				roleURL1 := make(map[string][]string)
				for j, role := range info.Roles {
					roleNames1[j] = role.Name
				}
				for _, user := range info.Users {
					userRoles1[user.Id] = user.Role
				}
				var proxy2 [][5]string
				if info.Proxy != nil {
					for _, url := range info.Proxy.Urls {
						proxy2 = append(proxy2, [5]string{info.Proxy.Domain, info.Proxy.Upstream, url.Method, url.Url, url.Id})
					}

					for _, role := range app.Roles {
						if len(role.UrlsId) != 0 {
							roleURL[role.Id] = role.UrlsId
							sort.Strings(roleURL[role.Id])
						}
					}
				}

				sort.Strings(roleNames)
				sort.Strings(roleNames1)
				require.Equal(t, roleNames, roleNames1)
				require.Equalf(t, userRoles, userRoles1, "expected app: %s, actual app: %s", w.JSONStringer(app), w.JSONStringer(info))
				require.Equalf(t, proxy, proxy2, "expected app: %s, actual app: %s", w.JSONStringer(app), w.JSONStringer(info))
				require.Equal(t, roleURL, roleURL1)
			}
		})
	}) {
		return
	}

	var appIds []string

	if !t.Run("Test List App", func(t *testing.T) {
		count, apps, err := svc.GetApps(ctx, "", nil, 1, 1024)
		require.NoError(t, err)
		require.Truef(t, len(apps) == 101, "Failed to obtain the information of all current apps: the number of users (%s) is not equal 101.", len(apps))
		require.Truef(t, count == 101, "Failed to obtain the information of all current apps: the number of users (%s) is not equal 101.", len(apps))
		for _, app := range apps {
			appIds = append(appIds, app.Id)
		}
		t.Run("Filter by keywords", func(t *testing.T) {
			_, apps, err = svc.GetApps(ctx, "AAA", nil, 1, 1024)
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
		appInfo, err := svc.GetAppInfo(ctx, opts.WithAppId(appId))
		require.NoError(t, err)

		roleNames := make([]string, len(appInfo.Roles))
		userRoles := make(map[string]string, len(appInfo.Users))
		for j, role := range appInfo.Roles {
			roleNames[j] = role.Name
		}
		for _, user := range appInfo.Users {
			userRoles[user.Id] = user.Role
		}
		newAppInfo := *appInfo
		newAppInfo.Description = rand.String(10)
		newAppInfo.Avatar = rand.String(10)
		newAppInfo.GrantMode = (newAppInfo.GrantMode + 1) % models.AppMeta_GrantMode(len(models.AppMeta_GrantMode_name))
		newAppInfo.GrantType = (newAppInfo.GrantType + 1) % models.AppMeta_GrantType(len(models.AppMeta_GrantType_name))
		newAppInfo.Status = (newAppInfo.Status + 1) % models.AppMeta_Status(len(models.AppMeta_Status_name))
		aUserIds := sets.New[string](userIds...)
		for _, user := range newAppInfo.Users {
			if aUserIds.Has(user.Id) {
				aUserIds.Delete(user.Id)
			}
		}
		appRole := models.AppRole{
			Name: "X_1asl",
		}
		roleNames = append(roleNames, appRole.Name)
		userId, ok := aUserIds.PopAny()
		if ok {
			newAppInfo.Users = append(newAppInfo.Users, &models.User{
				Model: models.Model{Id: userId},
				Role:  appRole.Name,
			})
			userRoles[userId] = appRole.Name
		}
		newAppInfo.Roles = append(newAppInfo.Roles, &appRole)
		tmpAppInfo := newAppInfo
		err = svc.UpdateApp(ctx, &tmpAppInfo)
		require.NoError(t, err)

		newAppInfo1, err := svc.GetAppInfo(ctx, opts.WithAppId(appId))
		require.NoError(t, err)
		require.Equal(t, newAppInfo.Name, newAppInfo1.Name)
		require.Equal(t, newAppInfo.Description, newAppInfo1.Description)
		require.Equal(t, newAppInfo.Avatar, newAppInfo1.Avatar)
		require.Equal(t, newAppInfo.GrantType, newAppInfo1.GrantType)
		require.Equal(t, newAppInfo.GrantMode, newAppInfo1.GrantMode)
		require.Equal(t, newAppInfo.Status, newAppInfo1.Status)

		roleNames1 := make([]string, len(newAppInfo1.Roles))
		userRoles1 := make(map[string]string, len(newAppInfo1.Users))
		for j, role := range newAppInfo1.Roles {
			roleNames1[j] = role.Name
		}
		for _, user := range newAppInfo1.Users {
			userRoles1[user.Id] = user.Role
		}
		sort.Strings(roleNames)
		sort.Strings(roleNames1)
		require.Equal(t, roleNames, roleNames1)
		require.Equal(t, userRoles, userRoles1)
	})

	t.Run("Test Delete App", func(t *testing.T) {
		err := svc.DeleteApp(ctx, appIds[rand.Intn(len(appIds))])
		require.NoError(t, err)

		t.Run("Test List App", func(t *testing.T) {
			count, apps, err := svc.GetApps(ctx, "", nil, 1, 1024)
			require.NoError(t, err)
			require.Truef(t, len(apps) == len(appIds)-1, "Failed to obtain the information of all current apps: the number of users (%d) is not equal %d.", len(apps), len(appIds)-1)
			require.Truef(t, int(count) == len(appIds)-1, "Failed to obtain the information of all current apps: the number of users (%d) is not equal %d.", len(apps), len(appIds)-1)
		})
	})
}
