package models

import (
	"github.com/unix-streamdeck/api"
	streamdeck "github.com/unix-streamdeck/driver"
)

type VirtualDev struct {
	Deck   streamdeck.Device
	Page   int
	IsOpen bool
	Config []api.Page
}

type Module struct {
	Name       string
	NewIcon    func() api.IconHandler
	NewKey     func() api.KeyHandler
	IconFields []api.Field
	KeyFields  []api.Field
}
