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

package models

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
	"sort"
	"testing"
	"time"
)

func getSeed(t *testing.T) []string {
	var ret []string
	rand.Seed(uint64(time.Now().UnixMilli()))
	for i := 0; i < rand.Intn(16); i++ {
		maxRead := rand.Intn(1024)
		buf := make([]byte, maxRead)
		_, err := rand.Read(buf)
		require.NoError(t, err)
		ret = append(ret, string(buf))
	}
	return ret
}

func TestNewId(t *testing.T) {
	type args struct {
		seed []string
	}
	tests := []struct {
		name string
		args args
	}{{
		name: "Success",
		args: struct{ seed []string }{seed: getSeed(t)},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var idList []string
			for i := 0; i < 1024; i++ {
				time.Sleep(time.Microsecond)
				idList = append(idList, NewId(tt.args.seed...))
			}
			time.Sleep(time.Second * time.Duration(rand.Intn(3)))
			for _, s := range idList {
				fmt.Println(s)
			}
			require.True(t, sort.StringsAreSorted(idList))
		})
	}
}
