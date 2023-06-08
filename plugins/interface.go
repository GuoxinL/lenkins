/*
Created by guoxin in 2023/6/5 20:49
*/
package plugins

import "lenkins"

// TODO plugin validate parameter
type Plugin interface {
	Validate() error
}

type PluginFunc func(lenkins.Config, interface{}) error
