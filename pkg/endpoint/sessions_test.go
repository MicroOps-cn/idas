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

package endpoint

import "testing"

func Test_maskEmail(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "len 1",
		args: args{addr: "1@1"},
		want: "*@*",
	}, {
		name: "len 3",
		args: args{addr: "123@456"},
		want: "1**@**6",
	}, {
		name: "len 8",
		args: args{addr: "12345678@12345678"},
		want: "1234****@****5678",
	}, {
		name: "len 20",
		args: args{addr: "12345678901234567890@12345678901234567890"},
		want: "12345***************@***************67890",
	}, {
		name: "domain",
		args: args{addr: "abc@example.com"},
		want: "a**@****ple.com",
	}, {
		name: "domain",
		args: args{addr: "abc@example."},
		want: "a**@****ple.",
	}, {
		name: "error",
		args: args{addr: "aaaaaaaaaaa"},
		want: "***",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maskEmail(tt.args.addr); got != tt.want {
				t.Errorf("maskEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
