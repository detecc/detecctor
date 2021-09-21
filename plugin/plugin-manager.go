package plugin

import (
	"fmt"
	"log"
	"plugin"
	"github.com/detecc/detecctor/config"
	"sync"
)

var pluginManager *PluginManager

func init() {
	once := sync.Once{}
	once.Do(func() {
		GetPluginManager()
	})
}

type PluginManager struct {
	plugins sync.Map
}

func Register(name string, plugin Handler) {
	GetPluginManager().AddPlugin(name, plugin)
}

func GetPluginManager() *PluginManager {
	if pluginManager == nil {
		pluginManager = &PluginManager{plugins: sync.Map{}}
	}
	return pluginManager
}

func (pm *PluginManager) HasPlugin(name string) bool {
	_, exists := pm.plugins.Load(name)
	return exists
}

func (pm *PluginManager) AddPlugin(name string, plugin Handler) {
	log.Println("Adding plugin to manager", name, plugin)
	if !pm.HasPlugin(name) {
		pm.plugins.Store(name, plugin)
	}
}

func (pm *PluginManager) GetPlugin(name string) (Handler, error) {
	mPlugin, exists := pm.plugins.Load(name)
	if exists {
		return mPlugin.(Handler), nil
	}
	return nil, fmt.Errorf("plugin doesnt exist")
}

func (pm *PluginManager) LoadPlugins() {
	log.Println("Loading plugins..")
	server := config.GetServerConfiguration().Server

	for _, pluginFromList := range server.Plugins {
		log.Println("Loading plugin:", pluginFromList)
		_, err := plugin.Open(fmt.Sprintf("%s/%s.so", server.PluginDir, pluginFromList))
		if err != nil {
			log.Println("error loading plugin", err)
			continue
		}
	}
}
