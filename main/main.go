/*
Created by guoxin in 2023/6/2 11:44
*/
package main

import (
	"encoding/json"
	"fmt"
	"lenkins"
	_ "lenkins/home"
	"lenkins/plugins"
	"lenkins/plugins/git"
	"lenkins/plugins/ssh/cmd"
	"lenkins/plugins/ssh/scp"
)

var (
	PluginMap = make(map[string]plugins.PluginFunc)
)

func init() {
	PluginMap["git"] = git.Execute
	PluginMap["cmd"] = cmd.Execute
	PluginMap["scp"] = scp.Execute
}

func main() {
	config, _, err := lenkins.LoadConfig("../config/deploy-test.yaml")
	if err != nil {
		panic(err)
		return
	}
	for _, job := range config.Jobs {
		marshal, err := json.Marshal(job)
		if err != nil {
			return
		}
		fmt.Println(string(marshal))
		fmt.Println("构建名称：", job.Name)
		fmt.Println("构建参数：", job.Parameters)
		for _, step := range job.Steps {
			fmt.Println("步骤名称：", step.Name)
			for pluginKey, parameter := range step.Plugin {
				parameterMap := parameter.(map[string]interface{})
				fmt.Println("\t插件名称：", pluginKey)
				fmt.Print("\t插件字段: ")
				prettyJson(parameter)
				pluginFunc := PluginMap[pluginKey]
				err := pluginFunc(parameterMap)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
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
