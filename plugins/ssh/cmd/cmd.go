/*
Created by guoxin in 2023/6/2 11:47
*/
package cmd

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lenkins"
	"lenkins/plugins/ssh"
	"log"
)

func Execute(cfg lenkins.Config, parameter interface{}) error {
	g := &Cmd{}
	err := mapstructure.Decode(parameter, g)
	if err != nil {
		return fmt.Errorf("failed to configure object mapping. error: %v", err)
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
		return fmt.Errorf("create ssh cmd client failed. error: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("create ssh session failed. error: %v", err)
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
