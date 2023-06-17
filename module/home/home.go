/*
Created by guoxin in 2023/6/2 15:00
*/
package home

import (
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"os"
	"path"
)

const (
	HomeDir = "~/.lenkins"
)

var (
	Home       = ""
	HomeDeploy = ""
	HomeLogs   = ""
)

func init() {
	Home, err := homedir.Expand(HomeDir)
	if err != nil {
		panic(err)
	}
	zap.S().Infof("home: ", Home)

	HomeDeploy = path.Join(Home, "deploy")
	err = Mkdir(HomeDeploy)
	if err != nil {
		panic(err)
	}
	zap.S().Infof("mkdir ", HomeDeploy)

	HomeLogs = path.Join(Home, "logs")
	err = Mkdir(HomeLogs)
	if err != nil {
		panic(err)
	}
	zap.S().Infof("mkdir ", HomeLogs)
}

func Mkdir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func Join(elem ...string) string {
	strings := []string{HomeDeploy}
	strings = append(strings, elem...)
	return path.Join(strings...)
}
