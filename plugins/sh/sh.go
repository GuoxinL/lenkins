package sh

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"lenkins/plugins"
	"lenkins/plugins/git"
	"os/exec"
	"strings"
)

const pluginName = "sh"

type Plugin struct {
	plugins.PluginInfo
	cmds []string
}

func New(info *plugins.PluginInfo) (plugins.Plugin, error) {
	var plugin = new(Plugin)
	plugin.PluginInfo = *info
	err := plugin.Unmarshal(&plugin.cmds)
	if err != nil {
		return nil, fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	return plugin, nil
}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) Validate() error {
	if len(p.cmds) == 0 {
		return errors.New("the commands cannot be empty")
	}
	for _, cmd := range p.cmds {
		if len(cmd) == 0 {
			return errors.New("the command cannot be empty")
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
	//for i := range p.cmds {
	//	p.cmds[i] = "-c " + p.cmds[i]
	//}
	return nil
}

func (p *Plugin) Execute() error {
	for i := range p.cmds {
		zap.S().Infof("[%v] execute command: %v", pluginName, p.cmds[i])
		cmdList := strings.Split(p.cmds[i], " ")
		for j := range cmdList {
			cmdList[j] = git.ReplaceScheme(cmdList[j], p.JobName)
		}
		cmd := strings.Join(cmdList, " ")
		c := exec.Command("/bin/bash", "-c", cmd)
		// 此处是windows版本
		// c := exec.Command("cmd", "/C", cmd)
		output, err := c.CombinedOutput()
		zap.S().Infof("[%v] command output: %s", pluginName, string(output))
		if err != nil {
			return err
		}
	}
	return nil
}
