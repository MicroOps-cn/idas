//go:build !make_test

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

package config

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/MicroOps-cn/idas/pkg/client/gorm"
	"github.com/MicroOps-cn/idas/pkg/testutils"
)

func TestConfig(t *testing.T) {
	logger := logs.New(logs.WithConfig(logs.MustNewConfig("debug", "logfmt")))
	logs.SetDefaultLogger(logger)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	ctx, logger = logs.NewContextLogger(ctx)
	tablePrefix := "idas_" + rand.String(10)
	schema := "idas_" + rand.String(10)
	var rawCfg string
	testutils.RunWithMySQLContainer(ctx, t, func(host, rootPassword string) {
		var testCfg Config
		t.Run("Test Marshal Config", func(t *testing.T) {
			mysqlOptions := gorm.NewMySQLOptions()
			mysqlOptions.Host = host
			mysqlOptions.Username = "root"
			mysqlOptions.Password = rootPassword
			mysqlOptions.TablePrefix = tablePrefix
			mysqlOptions.Schema = schema
			marshaler := jsonpb.Marshaler{
				Indent:   "    ",
				OrigName: true,
			}
			client := &gorm.MySQLClient{}
			client.SetOptions(mysqlOptions)
			testCfg.Storage = &Storages{
				Default: &Storage{
					Source: &Storage_Mysql{
						Mysql: client,
					},
				},
			}
			buf := bytes.NewBuffer(nil)
			err := marshaler.Marshal(buf, &testCfg)
			require.NoError(t, err, "Failed to Marshal config")
			fmt.Println(buf.String())
			rawCfg = buf.String()
		})
		t.Run("Test Unmarshal Config", func(t *testing.T) {
			err := safeCfg.ReloadConfigFromYamlReader(logger, NewConverter("./", bytes.NewReader([]byte(rawCfg))))
			require.NoError(t, err, "Failed to Unmarshal config")
			dftSource, ok := safeCfg.C.Storage.Default.Source.(*Storage_Mysql)
			require.True(t, ok)
			require.Equal(t, dftSource.Mysql.Options().Host, host)
			require.Equal(t, dftSource.Mysql.Options().Username, "root")
			require.Equal(t, dftSource.Mysql.Options().Schema, schema)
			require.Equal(t, dftSource.Mysql.Options().Password, rootPassword)
			require.Equal(t, dftSource.Mysql.Options().TablePrefix, tablePrefix)
		})
	})
}
