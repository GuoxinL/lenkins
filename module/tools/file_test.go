package tools

import "testing"

func TestWriteFile(t *testing.T) {
	type args struct {
		fileName string
		content  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "write",
			args: args{
				fileName: "xxx.txt",
				content:  []byte("xxx"),
			},
			wantErr: false,
		},
		{
			name: "fugai",
			args: args{
				fileName: "xxx.txt",
				content:  []byte("aaa"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteFile(tt.args.fileName, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
