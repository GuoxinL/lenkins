/*
Created by guoxin in 2023/6/5 20:49
*/
package plugins

import (
	"github.com/mitchellh/mapstructure"
)

type PluginInfos []*PluginInfo

type Plugin interface {
	Validate() error
	Replace() error
	Execute() error
}

type NewPluginFunc func(*PluginInfo) (Plugin, error)

type PluginInfo struct {
	JobName         string
	Parameters      map[string]string
	StepName        string
	PluginName      string
	PluginParameter interface{}
}

func (i *PluginInfo) Unmarshal(output interface{}) error {
	return mapstructure.Decode(i.Parameters, output)
}

func Build(jobName, stepName string, parameters map[string]string,
	pluginName string, pluginParameter interface{}) *PluginInfo {
	return &PluginInfo{
		JobName:         jobName,
		Parameters:      parameters,
		StepName:        stepName,
		PluginName:      pluginName,
		PluginParameter: pluginParameter,
	}
}
