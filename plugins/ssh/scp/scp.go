/*
Created by guoxin in 2023/6/2 17:34
*/
package scp

import (
	"errors"
	"fmt"

	"github.com/GuoxinL/lenkins/plugins"
	"github.com/GuoxinL/lenkins/plugins/git"
	"github.com/GuoxinL/lenkins/plugins/ssh"
	"go.uber.org/zap"
)

const pluginName = "scp"

type ScpType string

const (
	TypeUpload   ScpType = "upload"
	TypeDownload ScpType = "download"
)

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
		err := p.scp.scp(server, localPath, p.scp.Remote)
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
	Type    ScpType      `mapstructure:"type"`
}

func (s *Scp) scp(server ssh.Server, local string, remote string) error {
	client, err := server.GetScpClient()
	if err != nil {
		return fmt.Errorf("create ssh scp %v client failed. err: %v", s.Type, err)
	}
	defer client.Close()
	switch s.Type {
	case TypeUpload:
		err = client.Upload(local, remote)
		if err != nil {
			return fmt.Errorf("%s failed. error: %v", s.Type, err)
		}
	case TypeDownload:
		err = client.Download(local, remote)
		if err != nil {
			return fmt.Errorf("%s failed. error: %v", s.Type, err)
		}
	default:
		return fmt.Errorf("scp type %v not support", s.Type)
	}
	zap.S().Infof("[%s] File %s %s successfully!", pluginName, local, s.Type)
	return nil
}
