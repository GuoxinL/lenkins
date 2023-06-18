/*
Created by guoxin in 2023/6/2 17:37
*/
package ssh

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/GuoxinL/lenkins/plugins"
	"github.com/eleztian/go-scp"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	gossh "golang.org/x/crypto/ssh"
)

type Server struct {
	Host               string `mapstructure:"host"`
	User               string `mapstructure:"user"`
	Port               string `mapstructure:"port"`
	Password           string `mapstructure:"password"`
	PrivateKey         string `mapstructure:"privateKey"`
	PrivateKeyPathAuth string `mapstructure:"privateKeyPathAuth"`
}

func (s Server) Validate() error {
	if len(s.User) == 0 {
		return errors.New("the cmd server user parameter cannot be empty")
	}
	if len(s.Host) == 0 {
		return errors.New("the cmd server host parameter cannot be empty")
	}
	var ok bool
	if len(s.Password) != 0 {
		ok = true
	}
	if len(s.PrivateKey) != 0 {
		ok = true
	}
	if len(s.PrivateKeyPathAuth) != 0 {
		ok = true
	}
	if !ok {
		return errors.New("the cmd server Password or PrivateKey key or PrivateKeyPathAuth only one is not empty")
	}
	return nil
}

func (s *Server) Replace(key, value string) {
	s.Port = plugins.Replace(s.Port, key, value)
	s.User = plugins.Replace(s.User, key, value)
	s.Host = plugins.Replace(s.Host, key, value)
	s.Password = plugins.Replace(s.Password, key, value)
	s.PrivateKey = plugins.Replace(s.PrivateKey, key, value)
	s.PrivateKeyPathAuth = plugins.Replace(s.PrivateKeyPathAuth, key, value)
}

func (s *Server) GetConfig() (*gossh.ClientConfig, error) {
	config := &gossh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            s.User,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
	}
	if len(s.Password) != 0 {
		config.Auth = []gossh.AuthMethod{gossh.Password(s.Password)}
	}
	if len(s.PrivateKey) != 0 {
		authFunc, err := PublicKeyAuthFunc(s.PrivateKey)
		if err != nil {
			return nil, err
		}
		config.Auth = []gossh.AuthMethod{authFunc}
	}
	if len(s.PrivateKeyPathAuth) != 0 {
		authFunc, err := PublicKeyPathAuthFunc(s.PrivateKeyPathAuth)
		if err != nil {
			return nil, err
		}
		config.Auth = []gossh.AuthMethod{authFunc}
	}
	return config, nil
}
func (s *Server) GetCmdClient() (*gossh.Client, error) {
	config, err := s.GetConfig()
	if err != nil {
		return nil, err
	}
	return gossh.Dial("tcp", fmt.Sprintf("%s:%s", s.Host, s.Port), config)
}

func (s *Server) GetScpClient() (*scp.SCP, error) {
	config, err := s.GetConfig()
	if err != nil {
		return nil, err
	}
	return scp.New(fmt.Sprintf("%s:%s", s.Host, s.Port), config)
}

func PublicKeyPathAuthFunc(publicKeyPath string) (gossh.AuthMethod, error) {
	keyPath, err := homedir.Expand(publicKeyPath)
	if err != nil {
		zap.S().Fatal("find key's home dir failed", err)
		return nil, err
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		zap.S().Fatal("ssh key file read failed", err)
		return nil, err
	}
	// Create the Signer for this private key.
	signer, err := gossh.ParsePrivateKey(key)
	if err != nil {
		zap.S().Fatal("ssh key signer failed", err)
		return nil, err
	}
	return gossh.PublicKeys(signer), nil
}

func PublicKeyAuthFunc(privateKey string) (gossh.AuthMethod, error) {
	signer, err := gossh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		zap.S().Fatal("ssh key signer failed", err)
		return nil, err
	}
	return gossh.PublicKeys(signer), nil
}
