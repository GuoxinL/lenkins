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
	Home = "~/.lenkins"
)

var (
	Deploy = ""
	Logs   = ""
)

func init() {
	homeDir, err := homedir.Expand(Home)
	if err != nil {
		panic(err)
	}
	fmt.Println("home: ", homeDir)

	Deploy = path.Join(homeDir, "deploy")
	err = Mkdir(Deploy)
	if err != nil {
		panic(err)
	}
	fmt.Println("mkdir ", Deploy)

	Logs = path.Join(homeDir, "logs")
	err = Mkdir(Logs)
	if err != nil {
		panic(err)
	}
	fmt.Println("mkdir ", Logs)
}

func Mkdir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}
