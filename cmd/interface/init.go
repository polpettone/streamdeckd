package _interface

import modules2 "github.com/polpettone/streamdeckd/cmd/interface/modules"

func RegisterBaseModules() {
	RegisterModule(modules2.RegisterGif())
	RegisterModule(modules2.RegisterTime())
	RegisterModule(modules2.RegisterCounter())
	RegisterModule(modules2.RegisterIconState())
	RegisterModule(RegisterGame())
}
