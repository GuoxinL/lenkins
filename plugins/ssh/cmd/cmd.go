/*
Created by guoxin in 2023/6/2 11:47
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/GuoxinL/lenkins/plugins"
	"github.com/GuoxinL/lenkins/plugins/ssh"
	"go.uber.org/zap"
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

	for _, cmd := range cmds {
		session, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("create ssh session failed. err: %v", err)
		}
		zap.S().Infof("[%s] [%v] execute command.", pluginName, cmd)
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
		)

		session.Stdout = &stdout
		session.Stderr = &stderr

		err = session.Start("bash -c " + cmd)
		if err != nil {
			zap.S().Errorf("[%s] [%v] execute command failed. stdout: %v, stderr: %v, error: %v", pluginName, cmd, stdout.String(), stderr.String(), err)
			return err
		}
		zap.S().Infof("[%s] [%v] execute command success. stdout: %v, stderr: %v", pluginName, cmd, stdout.String(), stderr.String())
		session.Close()
	}
	return nil
}
