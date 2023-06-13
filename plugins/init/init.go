package init

import (
	"lenkins/plugins"
	"lenkins/plugins/git"
	"lenkins/plugins/sh"
	"lenkins/plugins/ssh/cmd"
	"lenkins/plugins/ssh/scp"
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
