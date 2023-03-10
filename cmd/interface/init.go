package _interface

import "github.com/polpettone/streamdeckd/cmd/handlers"

func RegisterBaseModules() {
	RegisterModule(handlers.RegisterGif())
	RegisterModule(handlers.RegisterTime())
	RegisterModule(handlers.RegisterCounter())
	RegisterModule(handlers.RegisterIconState())
}
