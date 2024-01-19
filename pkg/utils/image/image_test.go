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

package image

import (
	"context"
	"io"
	"testing"
)

func TestGenerateAvatar(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name       string
		args       args
		wantReader io.Reader
		wantErr    bool
	}{
		{name: "Test Avatar Generate", args: struct{ content string }{content: "测试"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReader, err := GenerateAvatar(context.Background(), tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAvatar() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotReader == nil {
				t.Errorf("GenerateAvatar() reader = %v, want: not null", err)
			}
		})
	}
}
