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
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	logs "github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/go-sql-driver/mysql"
	"github.com/moby/term"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/MicroOps-cn/idas/config"
)

type testServiceGenerate func(ctx context.Context, t *testing.T, testFunc func(name string, svc Service))

func newSqliteTestService(ctx context.Context, t *testing.T, testFunc func(name string, svc Service)) {
	const dsName = "sqlite"
	const sqliteYamlConfig = `
storage:
 default:
   sqlite:
     path: 'file:testdatabase?mode=memory&cache=shared'
   name: "sqlite"
`
	logger := logs.GetContextLogger(ctx)
	err := config.ReloadConfigFromYamlReader(logger, config.NewConverter("", bytes.NewBuffer([]byte(sqliteYamlConfig))))
	require.NoError(t, err)
	testFunc(dsName, New(ctx))
}

func runWithOpenLDAPContainer(ctx context.Context, t *testing.T, f func(host, rootPassword string)) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	customSchemaDir := path.Join(wd, "../../resource/openldap/schema")
	fmt.Println(customSchemaDir)
	rootPassword := rand.String(10)
	cli, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)
	t.Log("start openldap container...")

	// pull image
	images, err := cli.ImageList(ctx, types.ImageListOptions{Filters: filters.NewArgs(filters.Arg("reference", "osixia/openldap:1.5.0"))})
	require.NoError(t, err)
	if len(images) == 0 {
		resp, err := cli.ImagePull(ctx, "osixia/openldap:1.5.0", types.ImagePullOptions{})
		require.NoError(t, err)
		defer resp.Close()
		fd, isTerminal := term.GetFdInfo(os.Stdout)
		err = jsonmessage.DisplayJSONMessagesStream(resp, os.Stdout, fd, isTerminal, nil)
		require.NoError(t, err)
	}

	// start container
	c, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "osixia/openldap:1.5.0",
		Env: []string{
			`LDAP_ADMIN_PASSWORD=` + rootPassword,
			`LDAP_ORGANISATION=ops`,
			`LDAP_DOMAIN=microops.com`,
			`LDAP_REMOVE_CONFIG_AFTER_SETUP=false`,
			`DISABLE_CHOWN=true`,
			`LDAP_OPENLDAP_UID=0`,
			`LDAP_OPENLDAP_GID=0`,
			`LDAP_SEED_INTERNAL_SCHEMA_PATH=/etc/ldap/custom_schema/*.schema`,
		},
	}, &container.HostConfig{
		Binds: []string{customSchemaDir + ":/etc/ldap/custom_schema"},
	}, nil, nil, "")
	require.NoError(t, err)
	defer cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true})

	err = cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	defer cli.ContainerStop(ctx, c.ID, w.P[time.Duration](time.Minute))
	require.NoError(t, err)

	inspect, err := cli.ContainerInspect(ctx, c.ID)
	require.NoError(t, err)
	ticker := time.NewTicker(time.Second)
	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
loop:
	for {
		select {
		case <-timeoutCtx.Done():
			t.Error(timeoutCtx.Err())
		case <-ticker.C:
			_, err = net.Dial("tcp", inspect.NetworkSettings.IPAddress+":389")
			if err != nil {
				inspect, err = cli.ContainerInspect(ctx, c.ID)
				require.NoError(t, err)
				continue
			}
			time.Sleep(time.Second * 3)
			break loop
		}
	}
	t.Logf("openldap container <%s> started", c.ID)
	time.Sleep(time.Second * 20)
	f(inspect.NetworkSettings.IPAddress, rootPassword)
}

func runWithMySQLContainer(ctx context.Context, t *testing.T, f func(host, rootPassword string)) {
	rootPassword := rand.String(10)
	cli, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)
	t.Log("start mysql container...")

	// pull image
	images, err := cli.ImageList(ctx, types.ImageListOptions{Filters: filters.NewArgs(filters.Arg("reference", "mysql:5.7"))})
	require.NoError(t, err)
	if len(images) == 0 {
		resp, err := cli.ImagePull(ctx, "mysql:5.7", types.ImagePullOptions{})
		require.NoError(t, err)
		defer resp.Close()
		fd, isTerminal := term.GetFdInfo(os.Stdout)
		err = jsonmessage.DisplayJSONMessagesStream(resp, os.Stdout, fd, isTerminal, nil)
		require.NoError(t, err)
	}
	// start container
	c, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "mysql:5.7",
		Env:   []string{`MYSQL_ROOT_PASSWORD=` + rootPassword},
	}, &container.HostConfig{}, nil, nil, "")
	require.NoError(t, err)
	defer cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true})

	err = cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	defer cli.ContainerStop(ctx, c.ID, w.P[time.Duration](time.Minute))
	require.NoError(t, err)

	inspect, err := cli.ContainerInspect(ctx, c.ID)
	require.NoError(t, err)
	ticker := time.NewTicker(time.Second)
	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
loop:
	for {
		select {
		case <-timeoutCtx.Done():
			t.Error(timeoutCtx.Err())
		case <-ticker.C:
			_, err = net.Dial("tcp", inspect.NetworkSettings.IPAddress+":3306")
			if err != nil {
				inspect, err = cli.ContainerInspect(ctx, c.ID)
				require.NoError(t, err)
				continue
			}
			time.Sleep(time.Second * 1)
			break loop
		}
	}
	t.Logf("mysql container <%s> started", c.ID)

	// create schema
	connector, err := mysql.NewConnector(&mysql.Config{
		User:                 "root",
		Passwd:               rootPassword,
		Addr:                 inspect.NetworkSettings.IPAddress,
		AllowNativePasswords: true,
		Params:               map[string]string{"charset": "utf8mb4"},
		Collation:            "utf8mb4_general_ci",
	})
	require.NoError(t, err)
	conn, err := connector.Connect(ctx)
	require.NoError(t, err)
	prepare, err := conn.Prepare("CREATE SCHEMA `idas` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
	require.NoError(t, err)
	_, err = prepare.Exec(nil)
	require.NoError(t, err)
	f(inspect.NetworkSettings.IPAddress, rootPassword)
}

func newMySQLTestService(ctx context.Context, t *testing.T, testFunc func(name string, svc Service)) {
	runWithMySQLContainer(ctx, t, func(host, rootPassword string) {
		const dsName = "mysql"
		sqliteYamlConfig := fmt.Sprintf(`
storage:
 default:
   name: "mysql"
   mysql:
     host: "%s"
     schema: "idas"
     username: "root"
     password: "%s"
`, host, rootPassword)
		logger := logs.GetContextLogger(ctx)
		err := config.ReloadConfigFromYamlReader(logger, config.NewConverter("", bytes.NewBuffer([]byte(sqliteYamlConfig))))
		require.NoError(t, err)
		testFunc(dsName, New(ctx))
	})
}

func newLDAPTestService(ctx context.Context, t *testing.T, testFunc func(name string, svc Service)) {
	runWithOpenLDAPContainer(ctx, t, func(ldapHost, ldapPassword string) {
		runWithMySQLContainer(ctx, t, func(host, rootPassword string) {
			const dsName = "LDAP"
			ldapYamlConfig := fmt.Sprintf(`
storage:
  user:
  - name: "LDAP"
    ldap:
      host: "%s:389"
      manager_dn: "cn=admin,dc=microops,dc=com"
      manager_password: "%s"
      user_search_base: "ou=users,dc=microops,dc=com"
      app_search_base: "ou=groups,dc=microops,dc=com"
      attr_email: mail
      attr_user_display_name: cn
      attr_username: uid
  default:
    name: "mysql"
    mysql: 
      host: "%s"
      schema: "idas"
      username: "root"
      password: "%s"
`, ldapHost, ldapPassword, host, rootPassword)
			logger := logs.GetContextLogger(ctx)
			err := config.ReloadConfigFromYamlReader(logger, config.NewConverter("", bytes.NewBuffer([]byte(ldapYamlConfig))))
			require.NoError(t, err)
			testFunc(dsName, New(ctx))
		})
	})
}

func TestService(t *testing.T) {
	logs.SetDefaultLogger(logs.New())
	tests := []struct {
		name string
		sg   testServiceGenerate
	}{
		{name: "Test Sqlite", sg: newSqliteTestService},
		{name: "Test MySQL", sg: newMySQLTestService},
		{name: "Test LDAP", sg: newLDAPTestService},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
			defer cancel()
			tt.sg(ctx, t, func(storage string, svc Service) {
				if svc == nil {
					t.Logf("[%s] service is null, ignore testing...", tt.name)
				}
				var err error
				config.Get().Global.UploadPath, err = os.MkdirTemp("", strings.ReplaceAll(tt.name, " ", "_")+".")
				require.NoError(t, err)
				defer os.RemoveAll(config.Get().Global.UploadPath)
				if !t.Run("Test Auto migrate", func(t *testing.T) {
					require.NoError(t, svc.AutoMigrate(ctx))
				}) {
					return
				}
				testUserService(ctx, t, storage, svc)
				testAppService(ctx, t, storage, svc)
			})
		})
	}
}

func Test111(t *testing.T) {
	fmt.Println("\xe4\xbd\xa0")
}
