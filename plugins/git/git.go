/*
Created by guoxin in 2023/6/2 14:43

https://github.com/go-git/go-git/tree/master/_examples
*/
package git

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"path"
	"strings"

	"github.com/GuoxinL/lenkins/module/home"
	"github.com/GuoxinL/lenkins/module/logger"
	"github.com/GuoxinL/lenkins/plugins"
	"github.com/go-git/go-git/v5"
)

const (
	pluginName  = "git"
	goGitV5User = pluginName
	Dir         = pluginName
	Scheme      = pluginName + "://"
)

type authType string

const (
	basicAuth          authType = "basic"
	privateKeyAuth     authType = "privateKey"
	privateKeyPathAuth authType = "privateKeyPath"
)

type Plugin struct {
	plugins.PluginInfo
	git *Git
}

func New(info *plugins.PluginInfo) (plugins.Plugin, error) {
	var plugin = new(Plugin)
	plugin.PluginInfo = *info
	g := &Git{}
	err := plugin.Unmarshal(g)
	if err != nil {
		return nil, fmt.Errorf("failed to configure object mapping. err: %v", err)
	}
	plugin.git = g
	return plugin, nil
}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) Validate() error {
	if len(p.git.Repo) == 0 {
		return errors.New("the git repo parameter cannot be empty")
	}
	if len(p.git.Branch) == 0 {
		return errors.New("the git repo parameter cannot be empty")
	}
	switch p.git.AuthType {
	case basicAuth:
		if len(p.git.Username) == 0 {
			return errors.New("the git username parameter cannot be empty")
		}
		if len(p.git.Password) == 0 {
			return errors.New("the git password parameter cannot be empty")
		}
	case privateKeyAuth:
		if len(p.git.PrivateKey) == 0 {
			return errors.New("the git privateKey parameter cannot be empty")
		}
	case privateKeyPathAuth:
		if len(p.git.PrivateKeyPath) == 0 {
			return errors.New("the git privateKeyPath parameter cannot be empty")
		}
	default:
		return fmt.Errorf("git auth type %v not support", p.git.AuthType)
	}
	return nil
}

func (p *Plugin) Replace() error {
	for key, val := range p.Parameters {
		p.git.Repo = plugins.Replace(p.git.Repo, key, val)
		p.git.Branch = plugins.Replace(p.git.Branch, key, val)
		p.git.AuthType = authType(plugins.Replace(string(p.git.AuthType), key, val))
		p.git.PrivateKey = plugins.Replace(p.git.PrivateKey, key, val)
		p.git.PrivateKeyPath = plugins.Replace(p.git.PrivateKeyPath, key, val)
		p.git.Username = plugins.Replace(p.git.Username, key, val)
		p.git.Password = plugins.Replace(p.git.Password, key, val)
	}
	return nil
}

func (p *Plugin) Execute() error {
	err := p.git.Clone(home.Join(p.JobName, Dir))
	if err != nil {
		return err
	}
	return nil
}

type Git struct {
	Repo           string   `mapstructure:"repo"`
	Branch         string   `mapstructure:"branch"`
	AuthType       authType `mapstructure:"authType"`
	Username       string   `mapstructure:"username"`
	Password       string   `mapstructure:"password"`
	PrivateKey     string   `mapstructure:"privateKey"`
	PrivateKeyPath string   `mapstructure:"privateKeyPath"`
}

func (g *Git) Clone(filepath string) error {
	var (
		auth transport.AuthMethod
		err  error
	)
	if len(g.AuthType) == 0 {
		g.AuthType = privateKeyPathAuth
	}

	switch g.AuthType {
	case basicAuth:
		auth = &http.BasicAuth{
			Username: g.Username,
			Password: g.Password,
		}
	case privateKeyAuth:
		// Username must be "git" for SSH auth to work, not your real username.
		// See https://github.com/src-d/go-git/issues/637
		auth, err = ssh.NewPublicKeys(goGitV5User, []byte(g.PrivateKey), g.Password)
		if err != nil {
			return fmt.Errorf("git auth type %v, obtain the public key from the private key failed", err)
		}
	case privateKeyPathAuth:
		// Username must be "git" for SSH auth to work, not your real username.
		// See https://github.com/src-d/go-git/issues/637
		if len(g.PrivateKeyPath) == 0 {
			sshIdRsa, err := home.CurrentSshIdRSA()
			if err != nil {
				return err
			}
			g.PrivateKeyPath = sshIdRsa
		}

		auth, err = ssh.NewPublicKeysFromFile(goGitV5User, g.PrivateKeyPath, g.Password)
		if err != nil {
			return fmt.Errorf("git auth type %v, obtain the public key path from the private key failed", err)
		}
	default:
		return fmt.Errorf("git auth type %v not support", g.AuthType)
	}
	_, err = git.PlainClone(filepath, false, &git.CloneOptions{
		Auth:     auth,
		URL:      g.Repo,
		Progress: logger.GetWriter(path.Join(home.HomeLogs, "lenkins.log")),
	})
	return err
}

func ReplaceScheme(localPath, jobName string) string {
	if strings.Contains(localPath, Scheme) {
		localPath = strings.Replace(localPath, Scheme, "", -1)
		localPath = home.Join(jobName, Dir, localPath)
	}
	return localPath

}
