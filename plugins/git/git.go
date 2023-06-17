/*
Created by guoxin in 2023/6/2 14:43

https://github.com/go-git/go-git/tree/master/_examples
*/
package git

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/GuoxinL/lenkins/module/home"
	"github.com/GuoxinL/lenkins/module/logger"
	"github.com/GuoxinL/lenkins/plugins"
	git "github.com/go-git/go-git/v5"
)

const (
	pluginName = "git"
	Dir        = pluginName
	Scheme     = pluginName + "://"
)

type Plugin struct {
	plugins.PluginInfo
	git *Git
}

func New(info *plugins.PluginInfo) (plugins.Plugin, error) {
	var plugin = new(Plugin)
	plugin.PluginInfo = *info
	g := &Git{}
	err := plugin.Unmarshal(g)
	if err != nil {
		return nil, fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	plugin.git = g
	return plugin, nil
}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) Validate() error {
	if len(p.git.Repo) == 0 {
		return errors.New("the git repo parameter cannot be empty")
	}
	if len(p.git.Branch) == 0 {
		return errors.New("the git repo parameter cannot be empty")
	}
	return nil
}

func (p *Plugin) Replace() error {
	for key, val := range p.Parameters {
		p.git.Repo = plugins.Replace(p.git.Repo, key, val)
		p.git.Branch = plugins.Replace(p.git.Branch, key, val)
	}
	return nil
}

func (p *Plugin) Execute() error {
	err := p.git.Clone(home.Join(p.JobName, Dir))
	if err != nil {
		return err
	}
	return nil
}

type Git struct {
	Repo   string `mapstructure:"repo"`
	Branch string `mapstructure:"branch"`
}

func (g *Git) Clone(filepath string) error {
	_, err := git.PlainClone(filepath, false, &git.CloneOptions{
		URL:      g.Repo,
		Progress: logger.GetWriter(path.Join(home.HomeLogs, "lenkins.log")),
	})
	return err
}

func ReplaceScheme(localPath, jobName string) string {
	if strings.Contains(localPath, Scheme) {
		localPath = strings.Replace(localPath, Scheme, "", -1)
		localPath = home.Join(jobName, Dir, localPath)
	}
	return localPath

}
