/*
Created by guoxin in 2023/6/2 11:47
*/
package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	gossh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"lenkins/plugins/ssh"
	"log"
)

func Execute(parameter map[string]interface{}) error {
	g := &ssh.Cmd{}
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

func PublicKeyPathAuthFunc(publicKeyPath string) gossh.AuthMethod {
	keyPath, err := homedir.Expand(publicKeyPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := gossh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return gossh.PublicKeys(signer)
}

func PublicKeyAuthFunc(privateKey string) gossh.AuthMethod {
	signer, err := gossh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return gossh.PublicKeys(signer)
}
