/*
Created by guoxin in 2023/6/2 17:34
*/
package scp

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lenkins"
	errors "lenkins/err"
	"lenkins/home"
	"lenkins/plugins/ssh"
	"path"
)

const pluginName = "scp"

func Execute(job lenkins.Job, stepIndex int) error {
	step, parameter, ok := lenkins.GetConf(job, stepIndex, pluginName)
	if !ok {
		return errors.NoPluginUsed
	}

	g := &Scp{step: step}
	err := mapstructure.Decode(parameter, g)
	if err != nil {
		return fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	for _, server := range g.Servers {
		err := g.upload(server, g.Remote)
		if err != nil {
			return err
		}
	}
	return nil
}

type Scp struct {
	Servers []ssh.Server `mapstructure:"servers"`
	Remote  string       `mapstructure:"remote"`
	step    lenkins.Step
}

func (s *Scp) upload(server ssh.Server, remote string) error {
	client, err := server.GetScpClient()
	if err != nil {
		return fmt.Errorf("create ssh scp %v client failed. err: %v", "uoload", err)
	}
	defer client.Close()

	err = client.Upload(path.Join(home.HomeDeploy, s.step.Name), remote)
	if err != nil {
		return fmt.Errorf("%v failed. err: %v", "uoload", err)
	}
	fmt.Printf("File %s upload successfully!\n", "remotefile.txt")
	return nil
}

func ScpDownload(server ssh.Server) error {
	client, err := server.GetScpClient()
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Download("/root/scp", "testdata")
	if err != nil {
		return err
	}
	fmt.Printf("File %s downloaded successfully!\n", "remotefile.txt")
	return nil
}
