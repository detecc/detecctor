package plugin

import (
	"fmt"
	"github.com/detecc/detecctor/config"
	"log"
	"plugin"
	"sync"
)

var pluginManager *PluginManager

func init() {
	once := sync.Once{}
	once.Do(func() {
		GetPluginManager()
	})
}

// PluginManager is a manager for the plugins. It stores and maps the plugins to the command.
type PluginManager struct {
	plugins sync.Map
}

// Register a plugin to the manager.
func Register(name string, plugin Handler) {
	GetPluginManager().AddPlugin(name, plugin)
}

// GetPluginManager gets the plugin manager instance (singleton).
func GetPluginManager() *PluginManager {
	if pluginManager == nil {
		pluginManager = &PluginManager{plugins: sync.Map{}}
	}
	return pluginManager
}

// HasPlugin Check if the plugin exists in the manager.
func (pm *PluginManager) HasPlugin(name string) bool {
	_, exists := pm.plugins.Load(name)
	return exists
}

// AddPlugin Add a plugin to the manager.
func (pm *PluginManager) AddPlugin(name string, plugin Handler) {
	log.Println("Adding plugin to manager", name, plugin)
	if !pm.HasPlugin(name) {
		pm.plugins.Store(name, plugin)
	}
}

// GetPlugin returns the plugin, if found.
func (pm *PluginManager) GetPlugin(name string) (Handler, error) {
	mPlugin, exists := pm.plugins.Load(name)
	if exists {
		return mPlugin.(Handler), nil
	}
	return nil, fmt.Errorf("plugin doesnt exist")
}

// LoadPlugins Load the plugins from the folder, specified in the configuration file.
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
