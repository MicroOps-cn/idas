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
	"fmt"

	logs "github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/go-kit/log/level"
	"github.com/spf13/cobra"

	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/utils/signals"
)

var adminUsername string

// migrateCmd represents the migrate command
var initDataCmd = &cobra.Command{
	Use:   "init",
	Short: "Data initialization tool",
	Long:  `The data initialization tool will create a table with missing columns and indexes. And create the required user and application data.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logs.GetDefaultLogger()
		ctx := context.WithValue(cmd.Context(), "command", cmd.Use)
		InitData(ctx, signals.SetupSignalHandler(logger))
	},
}

func InitData(ctx context.Context, _ *signals.StopChan) {
	svc := service.New(ctx)
	if err := svc.AutoMigrate(ctx); err != nil {
		panic(err)
	}
	if err := svc.RegisterPermission(ctx, endpoint.Set{}.GetPermissionsDefine()); err != nil {
		panic(err)
	}
	if err := svc.InitData(ctx, adminUsername); err != nil {
		logger := logs.GetContextLogger(ctx)
		level.Error(logs.WithPrint(w.NewStringer(func() string {
			return fmt.Sprintf("%+v", err)
		}))(logger)).Log("msg", "failed to http request", "err", err)
		panic(err)
	}
}

func init() {
	initDataCmd.PersistentFlags().StringVar(&adminUsername, "admin", "admin", "admin username.")
	rootCmd.AddCommand(initDataCmd)
}
