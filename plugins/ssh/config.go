/*
Created by guoxin in 2023/6/2 17:37
*/
package ssh

import (
	"fmt"
	"github.com/eleztian/go-scp"
	gossh "golang.org/x/crypto/ssh"
	"lenkins/plugins/ssh/cmd"
	"time"
)

type Server struct {
	Host               string `mapstructure:"host"`
	User               string `mapstructure:"user"`
	Port               int    `mapstructure:"port"`
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
		config.Auth = []gossh.AuthMethod{cmd.PublicKeyAuthFunc(s.PrivateKey)}
	}
	if len(s.PrivateKeyPathAuth) != 0 {
		config.Auth = []gossh.AuthMethod{cmd.PublicKeyPathAuthFunc(s.PrivateKeyPathAuth)}
	}
	return config
}

type Cmd struct {
	Servers []Server `mapstructure:"servers"`
	Cmd     []string `mapstructure:"cmd"`
}

func (s *Server) GetCmdClient() (*gossh.Client, error) {
	config := s.GetConfig()
	return gossh.Dial("tcp", fmt.Sprintf("%s:%s", s.Host, s.Port), config)
}

type Scp struct {
	Servers []Server `mapstructure:"servers"`
	Remote  string   `mapstructure:"remote"`
}

func (s *Server) GetScpClient() (*scp.SCP, error) {
	config := s.GetConfig()
	return scp.New(fmt.Sprintf("%s:%s", s.Host, s.Port), config)
}
