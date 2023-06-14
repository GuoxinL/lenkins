/*
Created by guoxin in 2023/6/2 15:00
*/
package home

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
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
	fmt.Println("home: ", Home)

	HomeDeploy = path.Join(Home, "deploy")
	err = Mkdir(HomeDeploy)
	if err != nil {
		panic(err)
	}
	fmt.Println("mkdir ", HomeDeploy)

	HomeLogs = path.Join(Home, "logs")
	err = Mkdir(HomeLogs)
	if err != nil {
		panic(err)
	}
	fmt.Println("mkdir ", HomeLogs)
}

func Mkdir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func Join(elem ...string) string {
	strings := []string{HomeDeploy}
	strings = append(strings, elem...)
	return path.Join(strings...)
}
