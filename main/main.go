/*
Created by guoxin in 2023/6/2 11:44
*/
package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"lenkins"
	"lenkins/module/home"
	_ "lenkins/module/home"
	"lenkins/module/logger"
	"lenkins/plugins"
	"lenkins/plugins/plninit"
	"os"
	"path"
)

func main() {
	logger.InitLog(path.Join(home.HomeLogs, "lenkins.log"), zap.DebugLevel)
	config, _, err := lenkins.LoadConfig("../config/deploy-test.yaml")
	if err != nil {
		panic(err)
		return
	}
	var pluginInfos plugins.PluginInfos

	// 构建PluginInfo
	for _, job := range config.Jobs {
		marshal, err := json.Marshal(job)
		if err != nil {
			return
		}
		zap.S().Info(string(marshal))
		zap.S().Info("构建名称：", job.Name)
		zap.S().Info("构建参数：", job.Parameters)
		clearJobCache(job.Name)
		for _, step := range job.Steps {
			zap.S().Info("步骤名称：", step.Name)
			for pluginName, pluginParameter := range step.Plugin {
				pluginInfos = append(pluginInfos,
					plugins.Build(job.Name, step.Name, job.Parameters, pluginName, pluginParameter))
			}
		}
	}
	var pluginInstance []plugins.Plugin
	for _, info := range pluginInfos {
		newPlugin, ok := plninit.Plugins[info.PluginName]
		if !ok {
			zap.S().Errorf("[%v] new plugin failed. plugin not support", info.PluginName)
			return
		}
		plugin, err := newPlugin(info)
		if err != nil {
			zap.S().Errorf("[%s] new plugin failed. error: %v", info.PluginName, err)
			return
		}
		pluginInstance = append(pluginInstance, plugin)
		zap.S().Infof("[%v] new plugin success.", info.PluginName)
	}
	for i := range pluginInstance {
		err = pluginInstance[i].Replace()
		if err != nil {
			zap.S().Errorf("[%v] replace parameter failed. error: %v", err)
			return
		}
		zap.S().Infof("[%v] replace parameter success.", pluginInstance[i].Name())
	}
	for i := range pluginInstance {
		err = pluginInstance[i].Validate()
		if err != nil {
			zap.S().Errorf("[%v] validate parameter failed. error: %v", pluginInstance[i].Name(), err)
			return
		}
		zap.S().Infof("[%v] validate parameter success.", pluginInstance[i].Name())
	}
	for i := range pluginInstance {
		err = pluginInstance[i].Execute()
		if err != nil {
			zap.S().Errorf("[%v] execute failed. error: %v", pluginInstance[i].Name(), err)
			return
		}
		zap.S().Infof("[%v] execute success.", pluginInstance[i].Name())
	}
}

func clearJobCache(name string) {
	cachePath := home.Join(name)
	err := os.RemoveAll(cachePath)
	zap.S().Infof("remove %v cache success. path: %v, error: %v", name, cachePath, err)
}

func prettyJson(plugin interface{}) {
	res, _ := json.Marshal(plugin)
	zap.S().Infof(string(res))
	//var out bytes.Buffer
	//_ = json.Indent(&out, res, "", "\t")
	//_, err := out.WriteTo(os.Stdout)
	//if err != nil {
	//	return
	//}
	//fmt.Printf("\n")
}
