/*
Created by guoxin in 2023/6/2 11:44
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/GuoxinL/lenkins/module/config"
	"github.com/GuoxinL/lenkins/module/home"
	_ "github.com/GuoxinL/lenkins/module/home"
	"github.com/GuoxinL/lenkins/module/logger"
	"github.com/GuoxinL/lenkins/module/plugin"
	"github.com/GuoxinL/lenkins/plugins"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	DeployPath     = "config"
	ParametersPath = "parameters-path"

	SystemPath = "lenkins.log"
)

var (
	root = &cobra.Command{
		Use:   "lenkins",
		Short: "Lenkins CLI",
		Long: strings.TrimSpace(`Lenkins is a lightweight deployment tool. Lenkins can automatically execute scripts, deploy applicat
ions, and remotely execute commands through a configuration file; it supports git plug-ins, sh plug-
ins (local execution commands), cmd plug-ins (remote execution commands), scp plugins (upload or dow
nload) etc.`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return initMain(configPath)
		},
	}

	configPath     string
	parametersPath string
)

func init() {
	root.PersistentFlags().StringVarP(&configPath, DeployPath, "c", "", "Configuration file.")
	root.PersistentFlags().StringVarP(&parametersPath, ParametersPath, "p", "", "Parameters file path.")
}

func main() {
	err := root.Execute()
	if err != nil {
		panic(err)
	}
}

func initMain(deploy string) error {
	logger.InitLog(path.Join(home.HomeLogs, SystemPath), zap.DebugLevel)
	if len(deploy) == 0 {
		return errors.New("the configuration file cannot be empty")
	}
	conf, _, err := config.LoadYamlConfig(deploy)
	if err != nil {
		return err
	}

	// 构建PluginInfo
	var pluginInfos plugins.PluginInfos
	for _, job := range conf.Jobs {
		marshal, err := json.Marshal(job)
		if err != nil {
			zap.S().Errorf("[%v] json marshal failed. error: %v", job.Name, err)
			return fmt.Errorf("[%v] json marshal failed. error: %v", job.Name, err)
		}
		zap.S().Info(string(marshal))
		zap.S().Info("构建名称：", job.Name)
		zap.S().Debug("构建参数：", job.Parameters)
		//clearJobCache(job.Name)
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
		newPlugin, ok := plugin.Plugins[info.PluginName]
		if !ok {
			zap.S().Errorf("[%v] new plugin failed. plugin not support.", info.PluginName)
			return fmt.Errorf("[%v] new plugin failed. plugin not support", info.PluginName)
		}
		pluginIns, err := newPlugin(info)
		if err != nil {
			zap.S().Errorf("[%s] new plugin failed. error: %v", info.PluginName, err)
			return fmt.Errorf("[%s] new plugin failed. error: %v", info.PluginName, err)
		}
		pluginInstance = append(pluginInstance, pluginIns)
		zap.S().Infof("[%v] new plugin success.", info.PluginName)
	}

	for i := range pluginInstance {
		err = pluginInstance[i].Replace()
		if err != nil {
			zap.S().Errorf("[%v] replace parameter failed. error: %v", pluginInstance[i].Name(), err)
			return fmt.Errorf("[%v] replace parameter failed. error: %v", pluginInstance[i].Name(), err)
		}
		zap.S().Infof("[%v] replace parameter success.", pluginInstance[i].Name())
	}

	for i := range pluginInstance {
		err = pluginInstance[i].Validate()
		if err != nil {
			zap.S().Errorf("[%v] validate parameter failed. error: %v", pluginInstance[i].Name(), err)
			return fmt.Errorf("[%v] validate parameter failed. error: %v", pluginInstance[i].Name(), err)
		}
		zap.S().Infof("[%v] validate parameter success.", pluginInstance[i].Name())
	}

	for i := range pluginInstance {
		err = pluginInstance[i].Execute()
		if err != nil {
			zap.S().Errorf("[%v] execute failed. error: %v", pluginInstance[i].Name(), err)
			return fmt.Errorf("[%v] execute failed. error: %v", pluginInstance[i].Name(), err)
		}
		zap.S().Infof("[%v] execute success.", pluginInstance[i].Name())
	}

	return nil
}

func clearJobCache(name string) {
	cachePath := home.DeployJoin(name)
	_ = os.RemoveAll(cachePath)
	zap.S().Infof("remove %v cache success. path: %v", name, cachePath)
}
