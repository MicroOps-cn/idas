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

package models

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToken_To(t *testing.T) {
	type data struct {
		Name string `json:"name"`
	}
	type args interface{}
	tests := []struct {
		name    string
		token   Token
		args    args
		wantErr bool
		want    interface{}
	}{
		{name: "test struct to array", token: Token{Type: TokenTypeToken, Data: sql.RawBytes(`{"name":"lion"}`)}, args: &[]data{}, want: &[]data{{Name: "lion"}}, wantErr: false},
		{name: "test array to array", token: Token{Type: TokenTypeParent, Childrens: []*Token{{Data: sql.RawBytes(`{"name":"lion"}`)}}}, args: &[]data{}, want: &[]data{{Name: "lion"}}, wantErr: false},
		{name: "test array to struct", token: Token{Type: TokenTypeParent, Childrens: []*Token{{Data: sql.RawBytes(`{"name":"lion"}`)}}}, args: &data{}, want: &data{}, wantErr: true},
		{name: "test struct to struct", token: Token{Type: TokenTypeCode, Data: sql.RawBytes(`{"name":"lion"}`)}, args: &data{}, want: &data{Name: "lion"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.token.To(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.Equal(t, tt.args, tt.want)
		})
	}
}
