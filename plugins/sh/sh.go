package sh

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lenkins"
	"os/exec"
	"strings"
)

func Execute(cfg lenkins.Config, parameter interface{}) error {
	cmds := []string{}
	err := mapstructure.Decode(parameter, cmds)
	if err != nil {
		return fmt.Errorf("failed to configure object mapping. error: %v", err)
	}
	for _, cmd := range cmds {
		err := Command(cmd)
		if err != nil {
			return fmt.Errorf("sh execute command. error: %v", err)
		}
	}
	if err != nil {
		return fmt.Errorf("git clone failed. error: %v", err)
	}
	return nil
}

// 这里为了简化，我省去了stderr和其他信息
func Command(cmd string) error {
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
