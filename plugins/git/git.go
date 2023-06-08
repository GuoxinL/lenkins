/*
Created by guoxin in 2023/6/2 14:43

https://github.com/go-git/go-git/tree/master/_examples
*/
package git

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/mitchellh/mapstructure"
	"lenkins"
	"lenkins/home"
	"os"
	"path"
)

func Execute(cfg lenkins.Config, parameter interface{}) error {
	g := &Git{cfg: cfg}
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
	cfg    lenkins.Config
	Repo   string `mapstructure:"repo"`
	Branch string `mapstructure:"branch"`
}

func (g *Git) Clone() error {
	_, err := git.PlainClone(path.Join(home.HomeDeploy, g.cfg.Jobs), false, &git.CloneOptions{
		URL:      g.Repo,
		Progress: os.Stdout,
	})
	return err
}
