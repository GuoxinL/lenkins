package sh

import (
	"errors"
	"fmt"
	"lenkins/plugins"
	"os/exec"
)

const pluginName = "sh"

type Plugin struct {
	plugins.PluginInfo
	cmds []string
}

func New(info *plugins.PluginInfo) (plugins.Plugin, error) {
	var plugin = new(Plugin)
	plugin.PluginInfo = *info
	err := plugin.Unmarshal(plugin.cmds)
	if err != nil {
		return nil, fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	return plugin, nil
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
