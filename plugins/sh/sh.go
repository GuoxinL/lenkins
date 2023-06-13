package sh

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lenkins"
	errors "lenkins/err"
	"lenkins/plugins"
	"os/exec"
	"strings"
)

const pluginName = "sh"

type Plugin struct {
}

func New(info plugins.PluginInfo) error {
	return nil
}

func (p Plugin) validate() error {
	//TODO implement me
	panic("implement me")
}

func (p Plugin) replace() error {
	//TODO implement me
	panic("implement me")
}

func (p Plugin) Execute() error {
	//TODO implement me
	panic("implement me")
}

func Execute(job lenkins.Job, stepIndex int) error {
	cmds := []string{}
	step, parameter, ok := lenkins.GetConf(job, stepIndex, pluginName)
	if !ok {
		return errors.NoPluginUsed
	}

	err := mapstructure.Decode(parameter, cmds)
	if err != nil {
		return fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	for _, cmd := range cmds {
		err := Command(step, cmd)
		if err != nil {
			return fmt.Errorf("sh execute command. err: %v", err)
		}
	}
	if err != nil {
		return fmt.Errorf("git clone failed. err: %v", err)
	}
	return nil
}

// 这里为了简化，我省去了stderr和其他信息
func Command(_ lenkins.Step, cmd string) error {
	fmt.Println("execute command. ", cmd)
	cmd = "-c " + cmd
	cmds := strings.Split(" ", cmd)
	c := exec.Command("sh", cmds...)
	// 此处是windows版本
	// c := exec.Command("cmd", "/C", cmd)
	output, err := c.CombinedOutput()
	fmt.Println(string(output))
	return err
}
