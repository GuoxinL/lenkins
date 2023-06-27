/*
Created by guoxin in 2023/6/2 15:00
*/
package home

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
)

const (
	HomeDir = "~/.lenkins"
	SSHDir  = "~/.ssh"
)

var (
	Home             = ""
	HomeDeploy       = ""
	HomeLogs         = ""
	HomeRemoteConfig = ""
)

func init() {
	var err error
	Home, err = homedir.Expand(HomeDir)
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
	HomeRemoteConfig = path.Join(Home, "remote")
	err = Mkdir(HomeRemoteConfig)
	if err != nil {
		panic(err)
	}
	zap.S().Infof("mkdir ", HomeLogs)
}

func Mkdir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func DeployJoin(elem ...string) string {
	strings := []string{HomeDeploy}
	strings = append(strings, elem...)
	return path.Join(strings...)
}

func CurrentSshIdRSA() (string, error) {
	ssh, err := homedir.Expand(SSHDir)
	if err != nil {
		return "", err
	}
	sshIdRsa := path.Join(ssh, "id_rsa")
	return sshIdRsa, nil
}
