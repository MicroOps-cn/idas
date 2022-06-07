/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/go-kit/log/level"
	"os"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"

	"idas/pkg/logs"
	"idas/pkg/service"
	"idas/pkg/service/models"
)

var (
	username string
	password string
	email    string
	fullName string
	role     string
	storage  string
)

// migrateCmd represents the migrate command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User manager",
	Long:  `User manager tools.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "create user",
	Long:  `create user.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logs.New(&logConfig)
		logs.SetRootLogger(logger)
		svc := service.New(cmd.Context())
		if password == "-" {
			p, err := gopass.GetPasswdPrompt("please input password: ", true, os.Stdin, os.Stderr)
			if err != nil {
				level.Error(logger).Log("msg", "failed to create user", "err", err)
				os.Exit(1)
			}
			password = string(p)

		}
		if len(username) == 0 {
			level.Error(logger).Log("msg", "username is null")
			os.Exit(1)
		}
		if len(password) == 0 {
			level.Error(logger).Log("msg", "password is null")
			os.Exit(1)
		}
		var us []string
		if len(storage) > 0 {
			us = append(us, storage)
		} else if uss, _, err := svc.GetUserSource(cmd.Context()); err != nil {
			level.Error(logger).Log("msg", "failed to get user storage source")
			os.Exit(1)
		} else if len(uss) == 0 {
			level.Error(logger).Log("msg", "can't get user storage source config")
			os.Exit(1)
		} else {
			for _, val := range uss {
				us = append(us, val)
			}
		}

		if len(fullName) == 0 {
			fullName = username
		}

		for _, s := range us {
			_, err := svc.CreateUser(cmd.Context(), s, &models.User{
				Username: username,
				Password: []byte(password),
				Email:    email,
				FullName: fullName,
				Role:     models.UserRole(role),
				Status:   models.UserStatusNormal,
			})
			if err != nil {
				level.Error(logger).Log("msg", "failed to create user", "err", err, "storage", s)
			}
		}
	},
}

func init() {
	userCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username (login name).")
	userCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "user password.")
	userCmd.PersistentFlags().StringVarP(&email, "email", "e", "", "user email.")
	userCmd.PersistentFlags().StringVarP(&fullName, "fullname", "f", "", "user full name.")
	userCmd.PersistentFlags().StringVarP(&role, "role", "r", "user", "user/admin")
	userCmd.PersistentFlags().StringVarP(&storage, "storage", "s", "", "user storage source")

	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userAddCmd)
}
