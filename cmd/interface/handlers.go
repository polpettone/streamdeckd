package _interface

import (
	"github.com/polpettone/streamdeckd/cmd/models"
	"log"
	"plugin"
)

var modules []models.Module

func AvailableModules() []models.Module {
	return modules
}

func RegisterModule(m models.Module) {
	for _, module := range modules {
		if module.Name == m.Name {
			log.Println("Module already loaded: " + m.Name)
			return
		}
	}
	log.Println("Loaded module " + m.Name)
	modules = append(modules, m)
}

func LoadModule(path string) {
	plug, err := plugin.Open(path)
	if err != nil {
		//log.Println("Failed to load module: " + path)
		log.Println(err)
		return
	}
	mod, err := plug.Lookup("GetModule")
	if err != nil {
		log.Println(err)
		return
	}
	var modMethod func() models.Module
	modMethod, ok := mod.(func() models.Module)
	if !ok {
		log.Println("Failed to load module: " + path)
		return
	}
	RegisterModule(modMethod())
}
