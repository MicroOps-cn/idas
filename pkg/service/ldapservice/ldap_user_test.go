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

package ldapservice

import "testing"

func Test_getPhoneNumberFilter(t *testing.T) {
	type args struct {
		name    string
		phoneNo string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "+86", args: args{name: "phoneNumber", phoneNo: "+8612345678"}, want: "(|(phoneNumber=+8612345678)(phoneNumber=12345678))"},
		{name: "+86 ", args: args{name: "phoneNumber", phoneNo: "+86 12345678"}, want: "(|(phoneNumber=+86 12345678)(phoneNumber=12345678))"},
		{name: "+86-", args: args{name: "phoneNumber", phoneNo: "+86-12345678"}, want: "(|(phoneNumber=+86-12345678)(phoneNumber=12345678))"},
		{name: "no", args: args{name: "phoneNumber", phoneNo: "12345678"}, want: "(phoneNumber=12345678)"},
		{name: "+86+86-+86", args: args{name: "phoneNumber", phoneNo: "+86+86-+8612345678"}, want: "(|(phoneNumber=+86+86-+8612345678)(phoneNumber=+86-+8612345678))"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPhoneNumberFilter(tt.args.name, tt.args.phoneNo); got != tt.want {
				t.Errorf("getPhoneNumberFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
