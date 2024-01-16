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
	"fmt"
	"os"
	"strings"

	"github.com/MicroOps-cn/fuck/safe"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a string",
	Long:  `Encrypt strings using internal encryption algorithms.`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := gopass.GetPasswdPrompt("The string that needs to be encrypted (original string): ", true, os.Stdin, os.Stderr)
		if err != nil {
			fmt.Printf("System error: %s\n", err)
			os.Exit(1)
			return
		}
		key, err := gopass.GetPasswdPrompt("Please enter the key: ", true, os.Stdin, os.Stderr)
		if err != nil {
			fmt.Printf("System error: %s\n", err)
			os.Exit(1)
		}
		encrypt, err := safe.Encrypt(data, string(key), nil)
		if err != nil {
			fmt.Printf("Encrypt failed: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(encrypt)
		line := strings.Repeat("-", len(encrypt)+5)
		fmt.Printf("+%s+\n", line)
		fmt.Printf(fmt.Sprintf("| %%-%ds |\n", len(encrypt)+3), encrypt)
		fmt.Printf("+%s+\n", line)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
}
