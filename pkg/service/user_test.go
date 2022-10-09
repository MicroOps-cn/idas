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
	"bytes"
	"context"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/sign"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/rand"
	"net/http"
	"testing"
	"time"
)

func testUserService(t *testing.T, ctx context.Context, storage string, svc Service) {

	var userId string

	oriUser := models.User{
		Username:    "lion",
		Email:       "lion@idas.local",
		PhoneNumber: "+0112345678",
		FullName:    "Lion",
		Avatar:      "xxxxxxxxxxx",
		Status:      models.UserMeta_inactive,
	}

	if !t.Run("Test Create User", func(t *testing.T) {
		cUser := oriUser
		count, users, err := svc.GetUsers(ctx, storage, "", models.UserMeta_status_all, "", 1, 1024)
		require.NoError(t, err)
		require.Len(t, users, 0)
		require.Equal(t, count, int64(0))
		t.Run("Test Create Null User", func(t *testing.T) {
			_, err = svc.CreateUser(ctx, storage, &models.User{})
			require.Error(t, err)
		})
		for i := 0; i < 5; i++ {
			_, err = svc.CreateUser(ctx, storage, &models.User{
				Username:    rand.String(5),
				Email:       rand.String(7),
				PhoneNumber: rand.String(9),
				FullName:    rand.String(5),
				Avatar:      rand.String(20),
				Status:      models.UserMeta_UserStatus(rand.Intn(4)),
			})
			require.NoError(t, err)
		}

		user, err := svc.CreateUser(ctx, storage, &cUser)
		require.NoError(t, err)
		_, err = uuid.FromString(user.Id)
		require.NoError(t, err)
		userId = user.Id

		require.True(t, time.Now().Sub(user.CreateTime) < time.Second*3 && time.Now().Sub(user.CreateTime) > -time.Second)
		require.Equal(t, user.Username, "lion")
		require.Equal(t, user.FullName, "Lion")
		require.Equal(t, user.Email, "lion@idas.local")
		require.Equal(t, user.PhoneNumber, "+0112345678")
		require.Equal(t, user.Avatar, "xxxxxxxxxxx")
		require.Equal(t, user.Status, models.UserMeta_inactive)

		t.Run("Test Create Duplicate User", func(t *testing.T) {
			cUser = oriUser
			_, err = svc.CreateUser(ctx, storage, &cUser)
			require.Error(t, err)
		})
		for i := 0; i < 5; i++ {
			_, err = svc.CreateUser(ctx, storage, &models.User{
				Username:    rand.String(5),
				Email:       rand.String(7),
				PhoneNumber: rand.String(9),
				FullName:    rand.String(5),
				Avatar:      rand.String(20),
				Status:      models.UserMeta_UserStatus(rand.Intn(4)),
			})
			require.NoError(t, err)
		}

		count, users, err = svc.GetUsers(ctx, storage, "", models.UserMeta_status_all, "", 1, 20)
		require.NoError(t, err)
		require.Len(t, users, 11)
		require.Equal(t, count, int64(11))

		for _, u := range users {
			if u.Id == userId {
				_, err = uuid.FromString(u.Id)
				require.NoError(t, err)
				require.Truef(t, time.Now().Sub(u.CreateTime) <= time.Minute && time.Now().Sub(u.CreateTime) >= -time.Second, "now=%s, createTime=%s,sub=%s", time.Now(), u.CreateTime, time.Now().Sub(u.CreateTime).String())
				require.Equal(t, u.Username, "lion")
				require.Equal(t, u.FullName, "Lion")
				require.Equal(t, u.Email, "lion@idas.local")
				require.Equal(t, u.PhoneNumber, "+0112345678")
				require.Equal(t, u.Avatar, "xxxxxxxxxxx")
				require.Equal(t, u.Status, models.UserMeta_inactive)
			}
		}
	}) {
		return
	}

	t.Run("Test Get Users", func(t *testing.T) {
		count, users, err := svc.GetUsers(ctx, storage, "Asdooa299shdoiasgd8269bw3i7y9fdsahigf", models.UserMeta_status_all, "", 1, 20)
		require.NoError(t, err)
		require.Len(t, users, 0)
		require.Equal(t, count, int64(0))

		count, users, err = svc.GetUsers(ctx, storage, "", models.UserMeta_inactive, "", 1, 20)
		require.NoError(t, err)
		for _, user := range users {
			require.Equal(t, user.Status, models.UserMeta_inactive)
		}

		count, users, err = svc.GetUsers(ctx, storage, "", models.UserMeta_normal, "", 1, 20)
		require.NoError(t, err)
		for _, user := range users {
			require.Equal(t, user.Status, models.UserMeta_normal)
		}

		count, users, err = svc.GetUsers(ctx, "", "", oriUser.Status, "", 1, 20)
		require.NoError(t, err)
		found := false
		for _, user := range users {
			require.Equal(t, user.Status, oriUser.Status)
			if user.Id == userId {
				found = true
			}
		}
		require.Equal(t, found, true)
	})

	if !t.Run("Test Update User", func(t *testing.T) {
		user, err := svc.UpdateUser(ctx, storage, &models.User{
			Model:       models.Model{Id: userId},
			Username:    "lion_u",
			Email:       "lion_u@idas.local",
			PhoneNumber: "+01123456789",
			FullName:    "Lion_u",
			Avatar:      "xxxxxxxxxxx_u",
			Status:      models.UserMeta_normal,
		})
		require.NoError(t, err)
		_, err = uuid.FromString(user.Id)
		require.NoError(t, err)
		require.Truef(t, time.Now().Sub(user.CreateTime) <= time.Minute && time.Now().Sub(user.CreateTime) >= -time.Second, "now=%s, createTime=%s,sub=%s", time.Now(), user.CreateTime, time.Now().Sub(user.CreateTime).String())
		require.Equal(t, user.Username, "lion")
		require.Equal(t, user.FullName, "Lion_u")
		require.Equal(t, user.Email, "lion_u@idas.local")
		require.Equal(t, user.PhoneNumber, "+01123456789")
		require.Equal(t, user.Avatar, "xxxxxxxxxxx_u")
		require.Equal(t, user.Status, models.UserMeta_normal)
		count, users, err := svc.GetUsers(ctx, storage, "", models.UserMeta_status_all, "", 1, 20)
		require.NoError(t, err)
		require.Len(t, users, 11)
		require.Equal(t, count, int64(11))

		for _, u := range users {
			if u.Id == userId {
				_, err = uuid.FromString(u.Id)
				require.NoError(t, err)
				require.True(t, time.Now().Sub(u.CreateTime) < time.Second*3 && time.Now().Sub(u.CreateTime) > -time.Second)
				require.Equal(t, u.Username, "lion")
				require.Equal(t, u.FullName, "Lion_u")
				require.Equal(t, u.Email, "lion_u@idas.local")
				require.Equal(t, u.PhoneNumber, "+01123456789")
				require.Equal(t, u.Avatar, "xxxxxxxxxxx_u")
				require.Equal(t, u.Status, models.UserMeta_normal)
			}
		}
	}) {
		return
	}

	t.Run("Test Update some fields of users", func(t *testing.T) {
		user, err := svc.UpdateUser(ctx, storage, &models.User{
			Model:       models.Model{Id: userId},
			Username:    "lion_u2",
			Email:       "lion_u2@idas.local",
			PhoneNumber: "+011234567890",
			FullName:    "Lion_u2",
			Avatar:      "xxxxxxxxxxx_u2",
			Status:      models.UserMeta_inactive,
		}, "email", "avatar")
		require.NoError(t, err)
		_, err = uuid.FromString(user.Id)
		require.NoError(t, err)
		require.True(t, time.Now().Sub(user.CreateTime) < time.Second*3 && time.Now().Sub(user.CreateTime) > -time.Second)
		require.Equal(t, user.Username, "lion")
		require.Equal(t, user.FullName, "Lion_u")
		require.Equal(t, user.Email, "lion_u2@idas.local")
		require.Equal(t, user.PhoneNumber, "+01123456789")
		require.Equal(t, user.Avatar, "xxxxxxxxxxx_u2")
		require.Equal(t, user.Status, models.UserMeta_normal)
		count, users, err := svc.GetUsers(ctx, storage, "", models.UserMeta_status_all, "", 1, 20)
		require.NoError(t, err)
		require.Len(t, users, 11)
		require.Equal(t, count, int64(11))

		for _, u := range users {
			if u.Id == userId {
				_, err = uuid.FromString(u.Id)
				require.NoError(t, err)
				require.True(t, time.Now().Sub(u.CreateTime) < time.Second*3 && time.Now().Sub(u.CreateTime) > -time.Second)
				require.Equal(t, u.Username, "lion")
				require.Equal(t, u.FullName, "Lion_u")
				require.Equal(t, u.Email, "lion_u2@idas.local")
				require.Equal(t, u.PhoneNumber, "+01123456789")
				require.Equal(t, u.Avatar, "xxxxxxxxxxx_u2")
				require.Equal(t, u.Status, models.UserMeta_normal)
			}
		}
	})
	t.Run("Test Patch User", func(t *testing.T) {

		user, err := svc.PatchUser(ctx, storage, map[string]interface{}{"id": userId, "status": models.UserMeta_disable})
		require.NoError(t, err)
		require.Equal(t, user.Status, models.UserMeta_disable)

		_, users, err := svc.GetUsers(ctx, storage, "", models.UserMeta_status_all, "", 1, 20)
		require.NoError(t, err)

		require.Len(t, users, 11)
		for _, u := range users {
			if u.Id == userId {
				require.Equal(t, u.Status, models.UserMeta_disable)
			}
		}
	})
	var keyPairName = rand.String(rand.Intn(20))
	t.Run("Test Create User Keypair", func(t *testing.T) {
		keypair, err := svc.CreateUserKey(ctx, userId, keyPairName)
		require.NoError(t, err)
		require.NotEmpty(t, keypair.Secret)
		require.NotEmpty(t, keypair.Key)
		require.NotEmpty(t, keypair.Private)

		req, err := http.NewRequest("POST", "https://example.com/api/users", bytes.NewBuffer([]byte(`{"username":"lion"}`)))
		require.NoError(t, err)
		req.Header.Set("content-type", sign.MimeJSON)
		signStr, err := sign.GetSignFromHTTPRequest(req, keypair.Key, keypair.Secret, keypair.Private, sign.ECDSA)
		require.NoError(t, err)
		payload, err := sign.GetPayloadFromHTTPRequest(req)
		require.NoError(t, err)
		users, err := svc.Authentication(ctx, models.AuthMeta_signature, sign.ECDSA, keypair.Key, "", payload, signStr)
		require.NoError(t, err)
		require.Len(t, users, 1)
	})
	t.Run("Test List User Keypair", func(t *testing.T) {
		count, pairs, err := svc.GetUserKeys(ctx, userId, 1, 100)
		require.NoError(t, err)
		require.Len(t, pairs, 1)
		require.Equal(t, count, int64(1))
		require.Equal(t, keyPairName, pairs[0].Name)
		require.NotEmpty(t, pairs[0].Id)
		require.NotEmpty(t, pairs[0].CreateTime)
		require.NotEmpty(t, pairs[0].Key)
		require.Empty(t, pairs[0].Secret)
		require.Empty(t, pairs[0].Private)
	})

	t.Run("Test Delete User", func(t *testing.T) {
		err := svc.DeleteUser(ctx, storage, userId)
		require.NoError(t, err)
		count, users, err := svc.GetUsers(ctx, storage, "", models.UserMeta_status_all, "", 1, 20)
		require.NoError(t, err)
		require.Len(t, users, 10)
		require.Equal(t, count, int64(10))
	})
}
