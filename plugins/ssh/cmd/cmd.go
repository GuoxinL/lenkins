/*
Created by guoxin in 2023/6/2 11:47
*/
package cmd

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lenkins"
	errors "lenkins/err"
	"lenkins/plugins"
	"lenkins/plugins/ssh"
	"log"
)

const pluginName = "cmd"

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
	step, parameter, ok := lenkins.GetConf(job, stepIndex, pluginName)
	if !ok {
		return errors.NoPluginUsed
	}
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
	step    lenkins.Step
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
