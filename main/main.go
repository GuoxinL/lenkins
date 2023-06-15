/*
Created by guoxin in 2023/6/2 11:44
*/
package main

import (
	"encoding/json"
	"fmt"
	"lenkins"
	"lenkins/home"
	_ "lenkins/home"
	"lenkins/log"
	"lenkins/plugins"
	"lenkins/plugins/plninit"
	"os"
)

func main() {
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
		fmt.Println(string(marshal))
		fmt.Println("构建名称：", job.Name)
		fmt.Println("构建参数：", job.Parameters)
		clearJobCache(job.Name)
		for _, step := range job.Steps {
			fmt.Println("步骤名称：", step.Name)
			for pluginName, pluginParameter := range step.Plugin {
				pluginInfos = append(pluginInfos,
					plugins.Build(job.Name, step.Name, job.Parameters, pluginName, pluginParameter))
			}
		}
	}
	var pluginInstance []plugins.Plugin
	// 初始化插件
	for _, info := range pluginInfos {
		newPlugin, ok := plninit.Plugins[info.PluginName]
		if !ok {
			log.Errorf("plugin %v not support", info.PluginName)
			return
		}
		plugin, err := newPlugin(info)
		if err != nil {
			log.Errorf("new plugin failed.", err)
			return
		}
		pluginInstance = append(pluginInstance, plugin)
	}
	for i := range pluginInstance {
		err = pluginInstance[i].Replace()
		if err != nil {
			log.Errorf("replace failed. error: %v", err)
			return
		}
		err = pluginInstance[i].Validate()
		if err != nil {
			log.Errorf("validate failed. error: %v", err)
			return
		}
		err = pluginInstance[i].Execute()
		if err != nil {
			log.Errorf("execute failed. error: %v", err)
			return
		}
	}
}

func clearJobCache(name string) {
	cachePath := home.Join(name)
	err := os.RemoveAll(cachePath)
	log.Infof("remove %v cache success. path: %v, error: %v", name, cachePath, err)
}

func prettyJson(plugin interface{}) {
	res, _ := json.Marshal(plugin)
	fmt.Println(string(res))
	//var out bytes.Buffer
	//_ = json.Indent(&out, res, "", "\t")
	//_, err := out.WriteTo(os.Stdout)
	//if err != nil {
	//	return
	//}
	//fmt.Printf("\n")
}
