package myConfig

import (
	"github.com/polpettone/streamdeckd/cmd/models"
	"github.com/unix-streamdeck/api"
)

var MyConfig *api.Config
var MyDevs map[string]*models.VirtualDev

func SetConfig(config *api.Config) {
	MyConfig = config
}

func SetDevs(devs map[string]*models.VirtualDev) {
	MyDevs = devs
}
