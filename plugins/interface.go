/*
Created by guoxin in 2023/6/5 20:49
*/
package plugins

type PluginInfos []*PluginInfo

type Plugin interface {
	validate() error
	replace() error
	Execute() error
}

type NewPluginFunc func(*PluginInfo) error

type PluginInfo struct {
	JobName         string
	Parameters      map[string]string
	StepName        string
	PluginName      string
	PluginParameter interface{}
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
