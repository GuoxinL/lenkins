package git

import (
	"github.com/GuoxinL/lenkins/module/home"
	"os"
	"testing"
)

func TestGit_Clone(t *testing.T) {
	type fields struct {
		Repo           string
		Branch         string
		AuthType       authType
		Username       string
		Password       string
		PrivateKey     string
		PrivateKeyPath string
	}
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "git clone basic",
			fields: fields{
				Repo:     "https://github.com/git-fixtures/basic.git",
				Branch:   "master",
				AuthType: basicAuth,
				Username: "username",
				Password: "password",
			},
			args:    args{filepath: home.DeployJoin("example", Dir)},
			wantErr: true,
		},
		{
			name: "git clone privateKeyPath is empty",
			fields: fields{
				Repo:     "git@gitee.com:guoxinliu/test-private-key.git",
				Branch:   "master",
				AuthType: privateKeyPathAuth,
			},
			args:    args{filepath: home.DeployJoin("example", Dir)},
			wantErr: false,
		},
		{
			name: "git clone privateKeyPath",
			fields: fields{
				Repo:     "git@gitee.com:guoxinliu/test-private-key.git",
				Branch:   "master",
				AuthType: privateKeyPathAuth,
				PrivateKeyPath: func() string {
					rsa, err := home.CurrentSshIdRSA()
					if err != nil {
						panic(err)
					}
					return rsa
				}(),
			},
			args:    args{filepath: home.DeployJoin("example", Dir)},
			wantErr: false,
		},
		{
			name: "git clone privateKey",
			fields: fields{
				Repo:     "git@gitee.com:guoxinliu/test-private-key.git",
				Branch:   "master",
				AuthType: privateKeyPathAuth,
				PrivateKey: func() string {
					rsa, err := home.CurrentSshIdRSA()
					if err != nil {
						panic(err)
					}
					key, err := os.ReadFile(rsa)
					if err != nil {
						panic(err)
					}
					return string(key)
				}(),
			},
			args:    args{filepath: home.DeployJoin("example", Dir)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.RemoveAll(tt.args.filepath)
			if err != nil {
				panic(err)
			}
			g := &Git{
				Repo:           tt.fields.Repo,
				Branch:         tt.fields.Branch,
				AuthType:       tt.fields.AuthType,
				Username:       tt.fields.Username,
				Password:       tt.fields.Password,
				PrivateKey:     tt.fields.PrivateKey,
				PrivateKeyPath: tt.fields.PrivateKeyPath,
			}
			if err := g.Clone("", tt.args.filepath); (err != nil) != tt.wantErr {
				t.Errorf("Clone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
