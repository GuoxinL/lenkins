/*
Created by guoxin in 2023/6/5 20:49
*/
package plugins

import (
	"lenkins/plugins/git"
	"lenkins/plugins/ssh/cmd"
	"lenkins/plugins/ssh/scp"
)

type PluginFunc func(map[string]interface{}) error

var (
	PluginMap = make(map[string]PluginFunc)
)

func init() {
	PluginMap["git"] = git.Execute
	PluginMap["cmd"] = cmd.Execute
	PluginMap["scp"] = scp.Execute
}
