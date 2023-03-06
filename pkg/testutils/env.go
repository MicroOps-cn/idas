/*
 Copyright © 2023 MicroOps-cn.

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

package testutils

import (
	"context"
	"fmt"
	"net"
	"os"
	"path"
	"testing"
	"time"

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
)

func RunWithOpenLDAPContainer(ctx context.Context, t *testing.T, f func(host, rootPassword string)) {
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
	//nolint:errcheck
	defer cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true})

	err = cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	//nolint:errcheck
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

func RunWithMySQLContainer(ctx context.Context, t *testing.T, f func(host, rootPassword string)) {
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
	//nolint:errcheck
	defer cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true})

	err = cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
	//nolint:errcheck
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
