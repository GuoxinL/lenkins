/*
Created by guoxin in 2023/6/2 14:15
*/
package config

import (
	"fmt"
)
import "github.com/spf13/viper"

type Config struct {
	Version string `mapstructure:"version"`
	Jobs    []Job  `mapstructure:"jobs"`
}

type Job struct {
	Name       string            `mapstructure:"name"`
	Parameters map[string]string `mapstructure:"parameters"`
	Steps      []Step            `mapstructure:"steps"`
}

type Step struct {
	Name   string                 `mapstructure:"name"`
	Plugin map[string]interface{} `mapstructure:",remain"`
}

func LoadConfig(path string) (*Config, *viper.Viper, error) {
	var (
		err       error
		confViper *viper.Viper
	)

	if confViper, err = InitViper(path); err != nil {
		return nil, nil, fmt.Errorf("load sdk config failed, %s", err)
	}

	rootConf := &Config{}
	if err = confViper.Unmarshal(rootConf); err != nil {
		return nil, nil, fmt.Errorf("unmarshal config file failed, %s", err)
	}

	return rootConf, confViper, nil

}

func InitViper(confPath string) (*viper.Viper, error) {
	cmViper := viper.New()
	cmViper.SetConfigFile(confPath)
	if err := cmViper.ReadInConfig(); err != nil {
		return nil, err
	}
	cmViper.WatchConfig()
	return cmViper, nil
}

func GetConf(job Job, stepIndex int, pluginName string) (step Step, parameter interface{}, ok bool) {
	step = job.Steps[stepIndex]
	parameter, ok = step.Plugin[pluginName]
	return
}
