package home

import (
	"os/user"
	"path"
	"testing"
)

func TestCurrentSshIdRSA(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		return
	}
	idRsaPath := path.Join(current.HomeDir, ".ssh", "id_rsa")
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "current user dir",
			want:    idRsaPath,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CurrentSshIdRSA()
			if (err != nil) != tt.wantErr {
				t.Errorf("CurrentSshIdRSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CurrentSshIdRSA() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		return
	}
	lenkinsDeployPath := path.Join(current.HomeDir, ".lenkins", "deploy", "deploy")

	type args struct {
		elem []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "join path",
			args: args{[]string{"deploy"}},
			want: lenkinsDeployPath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeployJoin(tt.args.elem...); got != tt.want {
				t.Errorf("DeployJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMkdir(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		return
	}
	lenkinsPath := path.Join(current.HomeDir, ".lenkins")

	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "current lenkins dir",
			args:    args{dir: lenkinsPath},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Mkdir(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("Mkdir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
