/*
 Copyright © 2024 MicroOps-cn.

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
	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/idas/config"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
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

type testServiceConfig struct {
	name      string
	generator testServiceGenerate
}

var testServiceConfigs = []testServiceConfig{
	{name: "Test Sqlite", generator: newSqliteTestService},
}

func TestService(t *testing.T) {
	logs.SetDefaultLogger(logs.New(logs.WithConfig(logs.MustNewConfig("debug", "logfmt"))))

	for _, tt := range testServiceConfigs {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
			defer cancel()
			tt.generator(ctx, t, func(storage string, svc Service) {
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
				testUserService(ctx, t, svc)
				testAppService(ctx, t, svc)
			})
		})
	}
}
