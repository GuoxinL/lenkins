/*
Created by guoxin in 2023/6/2 14:15
*/
package config

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/GuoxinL/lenkins/module/home"
	"github.com/GuoxinL/lenkins/module/tools"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

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
	v := viper.New()
	v.SetConfigFile(confPath)
	v.SetConfigType("yaml")
	if strings.HasPrefix(confPath, "http://") || strings.HasPrefix(confPath, "https://") {
		filename, err := save2local(confPath)
		if err != nil {
			return nil, err
		}
		v.SetConfigFile(filename)

		zap.S().Infof("reading the remote configuration file succeeded.")
		zap.S().Infof("save to local. filename: %s", filename)
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	v.WatchConfig()
	return v, nil
}

func save2local(configUrl string) (string, error) {
	resp, err := http.Get(configUrl)
	if err != nil {
		return "", fmt.Errorf("read remote file failed. error: %s, url: %s", err, configUrl)
	}
	defer resp.Body.Close()

	parsedUrl, err := url.Parse(configUrl)
	if err != nil {
		return "", fmt.Errorf("failed to read the file name in the URL. error: %s, url: %s", err, configUrl)
	}
	filename := path.Base(parsedUrl.Path)

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response reader failed. error: %s, url: %s", err, configUrl)
	}

	filename = path.Join(home.HomeRemoteConfig, filename)

	if err = tools.WriteFile(filename, content); err != nil {
		return "", fmt.Errorf("write file failed. error: %s, filename: %v, content: %s", err, filename, string(content))
	}
	return filename, nil
}
