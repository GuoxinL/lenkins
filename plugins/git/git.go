/*
Created by guoxin in 2023/6/2 14:43

https://github.com/go-git/go-git/tree/master/_examples
*/
package git

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/mitchellh/mapstructure"
	"lenkins/home"
	"os"
)

func Execute(parameter map[string]interface{}) error {
	g := &Git{}
	err := mapstructure.Decode(parameter, g)
	if err != nil {
		return fmt.Errorf("failed to configure object mapping. error: %v", err)
	}
	err = g.Clone()
	if err != nil {
		return fmt.Errorf("git clone failed. error: %v", err)
	}
	return nil
}

type Git struct {
	Repo   string `mapstructure:"repo"`
	Branch string `mapstructure:"branch"`
}

func (g *Git) Clone() error {
	_, err := git.PlainClone(home.HomeDeploy, false, &git.CloneOptions{
		URL:      g.Repo,
		Progress: os.Stdout,
	})
	return err
}
