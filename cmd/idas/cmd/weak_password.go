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
	"bufio"
	"errors"
	"fmt"
	"os"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/fuck/signals"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/MicroOps-cn/idas/pkg/service"
)

var (
	weakPasswordFile string
	verifyPassword   string
)

// migrateCmd represents the migrate command
var weakPasswordCmd = &cobra.Command{
	Use:   "weak-password",
	Short: "Weak password management tooll",
	Long:  `The data initialization tool will create a table with missing columns and indexes. And create the required user and application data.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var weakPasswordImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import weak password to database.",
	Long:  `The data initialization tool will create a table with missing columns and indexes. And create the required user and application data.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(weakPasswordFile) == 0 {
			return errors.New("clear text weak password file cannot be empty")
		}
		ctx := cmd.Context()
		f, err := os.Open(weakPasswordFile)
		if err != nil {
			return err
		}
		stat, err := f.Stat()
		if err != nil {
			return err
		}
		ch := signals.SetupSignalHandler(logs.GetContextLogger(ctx))
		svc := service.New(cmd.Context())
		scanner := bufio.NewScanner(f)
		var batch []string

		bar := progressbar.Default(stat.Size())
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return nil
			case <-ch.Channel():
				return nil
			default:
				pass := scanner.Text()
				_ = bar.Add(len(pass) + 1)
				if len(pass) != 0 {
					batch = append(batch, pass)
				}
				if len(batch) > 1000 {
					if err := svc.InsertWeakPassword(ctx, batch...); err != nil {
						fmt.Printf("failed to insert weak password: %s\n", err)
					}
					batch = []string{}
				}
			}
		}
		return nil
	},
}

var weakPasswordVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify weak password.",
	Long:  `Verify if a password is a weak password.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(verifyPassword) == 0 {
			return errors.New("clear text password cannot be empty")
		}
		ctx := cmd.Context()
		svc := service.New(cmd.Context())
		if err := svc.VerifyWeakPassword(ctx, verifyPassword); err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

func init() {
	weakPasswordImportCmd.PersistentFlags().StringVarP(&weakPasswordFile, "weak-password-file", "f", "", "Clear text weak password file, one password per line.")
	weakPasswordVerifyCmd.PersistentFlags().StringVarP(&verifyPassword, "password", "p", "", "Clear text password")

	rootCmd.AddCommand(weakPasswordCmd)
	weakPasswordCmd.AddCommand(weakPasswordImportCmd)
	weakPasswordCmd.AddCommand(weakPasswordVerifyCmd)
}
