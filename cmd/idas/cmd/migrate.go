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

package cmd

import (
	"context"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/fuck/signals"
	"github.com/spf13/cobra"

	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/service"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "自动迁移工具",
	Long:  `自动迁移仅仅会创建表，缺少列和索引，并且不会改变现有列的类型或删除未使用的列以保护数据。`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logs.GetDefaultLogger()
		ctx := context.WithValue(cmd.Context(), "command", cmd.Use)
		Migrate(ctx, signals.SetupSignalHandler(logger))
	},
}

func Migrate(ctx context.Context, _ *signals.Handler) {
	svc := service.New(ctx)
	if err := svc.AutoMigrate(ctx); err != nil {
		panic(err)
	}
	if err := svc.RegisterPermission(ctx, endpoint.Set{}.GetPermissionsDefine()); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
