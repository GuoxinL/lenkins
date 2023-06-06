/*
Created by guoxin in 2023/6/2 14:43

https://github.com/go-git/go-git/tree/master/_examples
*/
package git

import (
	git "github.com/go-git/go-git/v5"
	"lenkins/home"
	"os"
)

func Execute(parameter map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

type Git struct {
	Repo   string `mapstructure:"repo"`
	Branch string `mapstructure:"branch"`
}

func (g *Git) Clone() error {
	_, err := git.PlainClone(home.Deploy, false, &git.CloneOptions{
		URL:      g.Repo,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}
