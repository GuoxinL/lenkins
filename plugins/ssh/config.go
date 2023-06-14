/*
Created by guoxin in 2023/6/2 17:37
*/
package ssh

import (
	"fmt"
	"github.com/eleztian/go-scp"
	"github.com/mitchellh/go-homedir"
	gossh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"time"
)

type Server struct {
	Host               string `mapstructure:"host"`
	User               string `mapstructure:"user"`
	Port               uint32 `mapstructure:"port"`
	Password           string `mapstructure:"password"`
	PrivateKey         string `mapstructure:"privateKey"`
	PrivateKeyPathAuth string `mapstructure:"privateKeyPathAuth"`
}

func (s *Server) GetConfig() *gossh.ClientConfig {
	config := &gossh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            s.User,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if len(s.Password) != 0 {
		config.Auth = []gossh.AuthMethod{gossh.Password(s.Password)}
	}
	if len(s.PrivateKey) != 0 {
		config.Auth = []gossh.AuthMethod{PublicKeyAuthFunc(s.PrivateKey)}
	}
	if len(s.PrivateKeyPathAuth) != 0 {
		config.Auth = []gossh.AuthMethod{PublicKeyPathAuthFunc(s.PrivateKeyPathAuth)}
	}
	return config
}
func (s *Server) GetCmdClient() (*gossh.Client, error) {
	config := s.GetConfig()
	return gossh.Dial("tcp", fmt.Sprintf("%s:%s", s.Host, s.Port), config)
}

func (s *Server) GetScpClient() (*scp.SCP, error) {
	config := s.GetConfig()
	return scp.New(fmt.Sprintf("%s:%s", s.Host, s.Port), config)
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
