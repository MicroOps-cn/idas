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
			//ff, err := os.OpenFile("/tmp/xxxxxxx.png", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			//if err != nil {
			//	t.Errorf("failed to create file: %s", err)
			//}
			//defer ff.Close()
			//if _, err = io.Copy(ff, gotReader); err != nil {
			//	t.Errorf("failed to write file: %s", err)
			//	return
			//}
		})
	}
}
