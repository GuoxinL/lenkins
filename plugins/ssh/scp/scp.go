/*
Created by guoxin in 2023/6/2 17:34
*/
package scp

import (
	"fmt"
	"lenkins/plugins/ssh"
)

func Execute(parameter map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func ScpUpload(server ssh.Server) error {
	client, err := server.GetScpClient()
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.Upload("testdata", "/root/scp")
	if err != nil {
		return err
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
