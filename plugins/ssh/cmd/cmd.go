/*
Created by guoxin in 2023/6/2 11:47
*/
package cmd

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"lenkins/plugins"
	"lenkins/plugins/ssh"
)

const pluginName = "cmd"

type Plugin struct {
	plugins.PluginInfo
	cmd Cmd
}

type Cmd struct {
	Servers []ssh.Server `mapstructure:"servers"`
	Cmds    []string     `mapstructure:"cmd"`
}

func New(info *plugins.PluginInfo) (plugins.Plugin, error) {
	var plugin = new(Plugin)
	plugin.PluginInfo = *info
	err := plugin.Unmarshal(&plugin.cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	return plugin, nil
}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) Validate() error {
	if len(p.cmd.Cmds) == 0 {
		return errors.New("the commands cannot be empty")
	}
	for _, cmd := range p.cmd.Cmds {
		if len(cmd) == 0 {
			return errors.New("the command cannot be empty")
		}
	}
	for i := range p.cmd.Servers {
		err := p.cmd.Servers[i].Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) Replace() error {
	for key, val := range p.Parameters {
		for i := range p.cmd.Cmds {
			p.cmd.Cmds[i] = plugins.Replace(p.cmd.Cmds[i], key, val)
		}
		for i := range p.cmd.Servers {
			p.cmd.Servers[i].Replace(key, val)
		}
	}
	for i := range p.cmd.Cmds {
		p.cmd.Cmds[i] = "-c " + p.cmd.Cmds[i]
	}
	return nil
}

func (p *Plugin) Execute() error {
	for i := range p.cmd.Servers {
		_ = RemoteCmds(p.cmd.Servers[i], p.cmd.Cmds)
	}
	return nil
}

func RemoteCmds(server ssh.Server, cmds []string) error {
	client, err := server.GetCmdClient()
	if err != nil {
		return fmt.Errorf("create ssh cmd client failed. err: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("create ssh session failed. err: %v", err)
	}
	defer session.Close()
	for _, cmd := range cmds {
		zap.S().Infof("[%s] cmd: %v", pluginName, cmd)
		output, err := session.CombinedOutput(cmd)
		zap.S().Infof("[%s] execute command. %v", pluginName, cmd)
		if err != nil {
			zap.S().Infof("[%s] execute command failed. output: %v", pluginName, output)
			return err
		}
		zap.S().Infof("[%s] execute command success. output: %v", pluginName, output)
	}
	return nil
}
