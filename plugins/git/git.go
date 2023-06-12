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
	"lenkins/err"
	"lenkins/home"
	"os"
	"path"
)

const pluginName = "git"

func Execute(job lenkins.Job, stepIndex int) error {
	step, parameter, ok := lenkins.GetConf(job, stepIndex, pluginName)
	if !ok {
		return errors.NoPluginUsed
	}
	g := &Git{step: step}
	err := mapstructure.Decode(parameter, g)
	if err != nil {
		return fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	err = g.Clone()
	if err != nil {
		return fmt.Errorf("git clone failed. err: %v", err)
	}
	return nil
}

type Git struct {
	Repo   string `mapstructure:"repo"`
	Branch string `mapstructure:"branch"`
	step   lenkins.Step
}

func (g *Git) Clone() error {
	_, err := git.PlainClone(path.Join(home.HomeDeploy, g.step.Name), false, &git.CloneOptions{
		URL:      g.Repo,
		Progress: os.Stdout,
	})
	return err
}
