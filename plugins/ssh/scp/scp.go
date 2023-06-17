/*
Created by guoxin in 2023/6/2 17:34
*/
package scp

import (
	"errors"
	"fmt"

	plugins "github.com/GuoxinL/lenkins/plugins"
	"github.com/GuoxinL/lenkins/plugins/git"
	"github.com/GuoxinL/lenkins/plugins/ssh"
	"go.uber.org/zap"
)

const pluginName = "scp"

type Plugin struct {
	plugins.PluginInfo
	scp *Scp
}

func New(info *plugins.PluginInfo) (plugins.Plugin, error) {
	var plugin = new(Plugin)
	plugin.PluginInfo = *info
	plugin.scp = &Scp{}
	err := plugin.Unmarshal(plugin.scp)
	if err != nil {
		return nil, fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	return plugin, nil
}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) Validate() error {
	for i := range p.scp.Servers {
		err := p.scp.Servers[i].Validate()
		if err != nil {
			return err
		}
	}
	if len(p.scp.Local) == 0 {
		return errors.New("the scp local parameter cannot be empty")
	}
	if len(p.scp.Remote) == 0 {
		return errors.New("the scp remote parameter cannot be empty")
	}
	return nil
}

func (p *Plugin) Replace() error {
	for key, val := range p.Parameters {
		p.scp.Local = plugins.Replace(p.scp.Local, key, val)
		p.scp.Remote = plugins.Replace(p.scp.Remote, key, val)
		for i := range p.scp.Servers {
			p.scp.Servers[i].Replace(key, val)
		}
	}
	return nil
}

func (p *Plugin) Execute() error {
	for _, server := range p.scp.Servers {
		var localPath = git.ReplaceScheme(p.scp.Local, p.JobName)
		err := p.scp.upload(server, localPath, p.scp.Remote)
		if err != nil {
			return err
		}
	}
	return nil
}

type Scp struct {
	Servers []ssh.Server `mapstructure:"servers"`
	Remote  string       `mapstructure:"remote"`
	Local   string       `mapstructure:"local"`
}

func (s *Scp) upload(server ssh.Server, local string, remote string) error {
	client, err := server.GetScpClient()
	if err != nil {
		return fmt.Errorf("create ssh scp %v client failed. err: %v", "uoload", err)
	}
	defer client.Close()
	err = client.Upload(local, remote)
	if err != nil {
		return fmt.Errorf("%v failed. err: %v", "uoload", err)
	}
	zap.S().Infof("[%v] File %s upload successfully!", pluginName, local)
	return nil
}

func ScpDownload(server ssh.Server) error {
	client, err := server.GetScpClient()
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Download("/root/scp", "testdata")
	if err != nil {
		return err
	}
	zap.S().Infof("File %s downloaded successfully!", "remotefile.txt")
	return nil
}
