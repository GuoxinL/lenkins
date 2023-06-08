/*
Created by guoxin in 2023/6/2 17:34
*/
package scp

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"lenkins/home"
	"lenkins/plugins/ssh"
)

func Execute(parameter map[string]interface{}) error {
	g := &ssh.Scp{}
	err := mapstructure.Decode(parameter, g)
	if err != nil {
		return fmt.Errorf("failed to configure object mapping. error: %v", err)
	}
	for _, server := range g.Servers {
		err := ScpUpload(server, g.Remote)
		if err != nil {
			return err
		}
	}
	return nil
}

func ScpUpload(server ssh.Server, remote string) error {
	client, err := server.GetScpClient()
	if err != nil {
		return fmt.Errorf("create ssh scp %v client failed. error: %v", "uoload", err)
	}
	defer client.Close()

	err = client.Upload(home.HomeDeploy, remote)
	if err != nil {
		return fmt.Errorf("%v failed. error: %v", "uoload", err)
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
