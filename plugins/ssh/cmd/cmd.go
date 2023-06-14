/*
Created by guoxin in 2023/6/2 11:47
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lenkins"
	"lenkins/plugins"
	"lenkins/plugins/ssh"
	"log"
	"os/exec"
)

const pluginName = "cmd"

type Plugin struct {
	plugins.PluginInfo
	cmd Cmd
}

func New(info *plugins.PluginInfo) (plugins.Plugin, error) {
	var plugin = new(Plugin)
	plugin.PluginInfo = *info
	err := plugin.Unmarshal(plugin.cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	return plugin, nil
}

func (p *Plugin) Validate() error {
	if len(p.cmd.Cmd) == 0 {
		return errors.New("the commands cannot be empty")
	}
	for _, cmd := range p.cmd.Cmd {
		if len(cmd) == 0 {
			return errors.New("the command cannot be empty")
		}
	}
	for i := range p.cmd.Servers {
		if len(p.cmd.Servers[i].User) == 0 {
			return errors.New("the git repo parameter cannot be empty")
		}
	}
	return nil
}

func (p *Plugin) Replace() error {
	for key, val := range p.Parameters {
		for i := range p.cmds {
			p.cmds[i] = plugins.Replace(p.cmds[i], key, val)
		}
	}
	for i := range p.cmds {
		p.cmds[i] = "-c " + p.cmds[i]
	}
	return nil
}

func (p *Plugin) Execute() error {
	for i := range p.cmds {
		fmt.Println("execute command. ", p.cmds[i])
		c := exec.Command("sh", p.cmds[i])
		// 此处是windows版本
		// c := exec.Command("cmd", "/C", cmd)
		output, err := c.CombinedOutput()
		fmt.Println(string(output))
		return err
	}
	return nil
}

func Execute(job lenkins.Job, stepIndex int) error {
	step, parameter, ok := lenkins.GetConf(job, stepIndex, pluginName)
	g := &Cmd{step: step}
	err := mapstructure.Decode(parameter, g)
	if err != nil {
		return fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	for _, server := range g.Servers {
		err := RemoteCmds(server, g.Cmd)
		if err != nil {
			return err
		}
	}
	return err
}

type Cmd struct {
	Servers []ssh.Server `mapstructure:"servers"`
	Cmd     []string     `mapstructure:"cmd"`
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
		log.Println("cmd:", cmd)
		output, err := session.CombinedOutput(cmd)
		fmt.Println(fmt.Sprintf("execute command. %v", cmd))
		if err != nil {
			fmt.Println(fmt.Sprintf("execute command failed. output: %v", output))
			return err
		}
		fmt.Println(fmt.Sprintf("execute command success. output: %v", output))
	}
	return nil
}
