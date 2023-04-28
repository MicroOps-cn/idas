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

package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"os"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/spf13/cobra"

	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/transport"
	"github.com/MicroOps-cn/idas/pkg/utils/signals"
)

var openapiOutputFile string

// migrateCmd represents the migrate command
var rootCmd = &cobra.Command{
	Use:   "openapi",
	Short: "OpenAPI generator",
	Long:  `OpenAPI generator`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logs.GetDefaultLogger()
		ch := signals.SetupSignalHandler(logger)
		ctx, cancelFunc := context.WithCancel(cmd.Context())
		go func() {
			<-ch.Channel()
			cancelFunc()
		}()
		ctx = context.WithValue(ctx, global.HTTPWebPrefixKey, "/")
		handler := transport.NewHTTPHandler(ctx, logger, endpoint.Set{}, nil, nil, "/apidocs.json")

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/apidocs.json", nil)
		handler.ServeHTTP(w, req)
		if len(openapiOutputFile) == 0 {
			_, _ = io.Copy(os.Stdout, w.Body)
		} else if err := os.WriteFile(openapiOutputFile, w.Body.Bytes(), 0600); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&openapiOutputFile, "output", "o", "", "Output openAPI to the specified file")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
