/*
Created by guoxin in 2023/6/2 11:47
*/
package cmd

import (
	"github.com/mitchellh/go-homedir"
	gossh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"lenkins/plugins/ssh"
	"log"
)

func Execute(parameter map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func RemoteCmds(server ssh.Server, cmds []string) error {
	client, err := server.GetCmdClient()
	if err != nil {
		//log.Fatal("创建ssh client 失败", err)
		return err
	}
	defer client.Close()

	//创建ssh-session
	session, err := client.NewSession()
	if err != nil {
		//log.Fatal("创建ssh session 失败", err)
		return err
	}
	defer session.Close()
	//执行远程命令
	for _, cmd := range cmds {
		log.Println("cmd:", cmd)
		//output, err := session.CombinedOutput(cmd)
		if err != nil {
			//log.Fatal("cmd failed. error: ", err)
			return err
		}
		//log.Println("cmd output:", string(output))
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
