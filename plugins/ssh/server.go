/*
Created by guoxin in 2023/6/2 17:37
*/
package ssh

import (
	"errors"
	"fmt"
	"github.com/GuoxinL/lenkins/module/home"
	"os"
	"time"

	"github.com/GuoxinL/lenkins/plugins"
	"github.com/eleztian/go-scp"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	gossh "golang.org/x/crypto/ssh"
)

type authType string

const (
	passwordAuth       authType = "password"
	privateKeyAuth     authType = "privateKey"
	privateKeyPathAuth authType = "privateKeyPath"
)

type Server struct {
	Host           string   `mapstructure:"host"`
	Port           string   `mapstructure:"port"`
	AuthType       authType `mapstructure:"authType"`
	Username       string   `mapstructure:"username"`
	Password       string   `mapstructure:"password"`
	PrivateKey     string   `mapstructure:"privateKey"`
	PrivateKeyPath string   `mapstructure:"privateKeyPath"`
}

func (s Server) Validate() error {
	if len(s.Username) == 0 {
		return errors.New("the username parameter cannot be empty")
	}
	if len(s.Host) == 0 {
		return errors.New("the host parameter cannot be empty")
	}
	switch s.AuthType {
	case passwordAuth:
		if len(s.Password) != 0 {
			return errors.New("the password parameter cannot be empty")
		}
	case privateKeyAuth:
		if len(s.PrivateKey) != 0 {
			return errors.New("the privateKey parameter cannot be empty")
		}
	case privateKeyPathAuth:
		// privateKeyPathAuth can be empty
	default:
		return fmt.Errorf("auth type %v not support", s.AuthType)
	}
	return nil
}

func (s *Server) Replace(key, value string) {
	s.Port = plugins.Replace(s.Port, key, value)
	s.Host = plugins.Replace(s.Host, key, value)
	s.Username = plugins.Replace(s.Username, key, value)
	s.Password = plugins.Replace(s.Password, key, value)
	s.PrivateKey = plugins.Replace(s.PrivateKey, key, value)
	s.PrivateKeyPath = plugins.Replace(s.PrivateKeyPath, key, value)
}

func (s *Server) GetConfig() (*gossh.ClientConfig, error) {
	var (
		config *gossh.ClientConfig
	)

	config = &gossh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            s.Username,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
	}
	if len(s.Password) != 0 {
		config.Auth = []gossh.AuthMethod{gossh.Password(s.Password)}
	}
	if len(s.AuthType) == 0 {
		s.AuthType = privateKeyPathAuth
	}

	switch s.AuthType {
	case passwordAuth:
		config.Auth = []gossh.AuthMethod{gossh.Password(s.Password)}
	case privateKeyAuth:
		authFunc, err := PublicKeyAuthFunc(s.PrivateKey)
		if err != nil {
			return nil, err
		}
		config.Auth = []gossh.AuthMethod{authFunc}
	case privateKeyPathAuth:
		if len(s.PrivateKeyPath) == 0 {
			sshIdRsa, err := home.CurrentSshIdRSA()
			if err != nil {
				return nil, err
			}
			s.PrivateKeyPath = sshIdRsa
		}
		authFunc, err := PublicKeyPathAuthFunc(s.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
		config.Auth = []gossh.AuthMethod{authFunc}
	default:
		return nil, fmt.Errorf("ssh auth type %v not support", s.AuthType)
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
		zap.S().Errorf("find key's home dir failed", err)
		return nil, err
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		zap.S().Errorf("ssh key file read failed", err)
		return nil, err
	}
	// Create the Signer for this private key.
	signer, err := gossh.ParsePrivateKey(key)
	if err != nil {
		zap.S().Errorf("ssh key signer failed", err)
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
