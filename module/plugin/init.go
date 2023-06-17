package plugin

import (
	"github.com/GuoxinL/lenkins/plugins"
	"github.com/GuoxinL/lenkins/plugins/git"
	"github.com/GuoxinL/lenkins/plugins/sh"
	"github.com/GuoxinL/lenkins/plugins/ssh/cmd"
	"github.com/GuoxinL/lenkins/plugins/ssh/scp"
)

var (
	Plugins = make(map[string]plugins.NewPluginFunc)
)

func init() {
	Plugins["git"] = git.New
	Plugins["sh"] = sh.New
	Plugins["cmd"] = cmd.New
	Plugins["scp"] = scp.New
}
