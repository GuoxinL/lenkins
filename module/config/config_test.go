package config

import (
	"github.com/GuoxinL/lenkins/module/tools"
	"github.com/spf13/viper"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tools.SetCurrentProjectWorkingDir()

	c := &Config{
		Version: "v1.0.0",
		Jobs: []Job{{
			Name: "Hello Lenkins!",
			Steps: []Step{
				{
					Name:   "Hello world",
					Plugin: map[string]interface{}{"sh": []string{"echo 'Hello Lenkins!'"}},
				},
			},
		}},
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		want1   *viper.Viper
		wantErr bool
	}{
		{
			name:    "remote config",
			args:    args{path: "https://gitee.com/guoxinliu/lenkins/raw/master/example/primary.yaml"},
			want:    c,
			want1:   nil,
			wantErr: false,
		},
		{
			name:    "local config",
			args:    args{"./example/primary.yaml"},
			want:    c,
			want1:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := LoadYamlConfig(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadYamlConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Jobs[0].Name != tt.want.Jobs[0].Name {
				t.Errorf("LoadYamlConfig() JobsName got = %v, want %v", got, tt.want)
			}
			if got.Jobs[0].Steps[0].Name != tt.want.Jobs[0].Steps[0].Name {
				t.Errorf("LoadYamlConfig() StepName got = %v, want %v", got, tt.want)
			}
			if len(got.Jobs[0].Steps[0].Plugin) != len(tt.want.Jobs[0].Steps[0].Plugin) {
				t.Errorf("LoadYamlConfig() StepName got = %v, want %v", got, tt.want)
			}
			if got1 == tt.want1 {
				t.Errorf("LoadYamlConfig() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
