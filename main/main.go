/*
Created by guoxin in 2023/6/2 11:44
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lenkins"
	_ "lenkins/home"
	"os"
)

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
			for pluginKey, pluginParameter := range step.Plugin {
				fmt.Println("\t插件名称：", pluginKey)
				fmt.Print("\t插件字段: ")
				prettyJson(pluginParameter)
			}
		}
	}
}

func prettyJson(plugin interface{}) {
	res, _ := json.Marshal(plugin)
	var out bytes.Buffer
	_ = json.Indent(&out, res, "", "\t")
	_, err := out.WriteTo(os.Stdout)
	if err != nil {
		return
	}
	fmt.Printf("\n")
}
