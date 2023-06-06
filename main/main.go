/*
Created by guoxin in 2023/6/2 11:44
*/
package main

import (
	"encoding/json"
	"fmt"
	"lenkins"
	_ "lenkins/home"
)

type typeFunc func(string, interface{}) error

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
		fmt.Println("构建名称：", job.Parameters)
		for _, step := range job.Steps {
			fmt.Println("步骤名称：", step.Name)
			for pluginKey, pluginParameter := range step.Plugin {
				fmt.Println("插件名称：", pluginKey)
				pluginFields := pluginParameter.(map[string]interface{})
				for fieldName, fieldVal := range pluginFields {
					fmt.Println("插件字段: ", fieldName, fieldVal)
				}
				fmt.Println("插件名称：", pluginParameter)
			}
		}
		//fmt.Print("SourceBefore cmd：")
		//for _, cmd := range deploy.SourceBefore {
		//	if len(cmd) == 0 {
		//		continue
		//	}
		//	fmt.Println(cmd, " ")
		//}
		//fmt.Println()
		//
		//fmt.Println("SourceType:", deploy.Source.Type)
		//fmt.Println("SourceGitRepo:", deploy.Source.Git.Repo)
		//fmt.Println("SourceGitBranch:", deploy.Source.Git.Branch)
		//
		//fmt.Print("SourceAfter cmd：")
		//for _, cmd := range deploy.SourceAfter {
		//	if len(cmd) == 0 {
		//		continue
		//	}
		//	fmt.Println(cmd, " ")
		//}
		//fmt.Println()

		//deploy.TargetBefore
		//deploy.Target
		//deploy.TargetAfter
	}
}
