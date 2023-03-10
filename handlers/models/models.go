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
