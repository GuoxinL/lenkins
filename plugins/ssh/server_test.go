/*
Created by guoxin in 2023/6/21 15:51
*/
package ssh

import (
	gossh "golang.org/x/crypto/ssh"
	"reflect"
	"testing"
)

func TestServer_GetCmdClient(t *testing.T) {
	type fields struct {
		Host           string
		User           string
		Port           string
		Type           sshType
		Password       string
		PrivateKey     string
		PrivateKeyPath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *gossh.Client
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				Host:           "172.16.12.93",
				User:           "root",
				Port:           "22",
				Type:           "privateKeyPath",
				Password:       "",
				PrivateKey:     "",
				PrivateKeyPath: "~/.ssh/id_ras",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Host:           tt.fields.Host,
				User:           tt.fields.User,
				Port:           tt.fields.Port,
				Type:           tt.fields.Type,
				Password:       tt.fields.Password,
				PrivateKey:     tt.fields.PrivateKey,
				PrivateKeyPath: tt.fields.PrivateKeyPath,
			}
			got, err := s.GetCmdClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCmdClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCmdClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}
