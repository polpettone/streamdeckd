package _interface

import "github.com/polpettone/streamdeckd/cmd"

func RegisterBaseModules() {
	cmd.RegisterModule(RegisterGif())
	cmd.RegisterModule(RegisterTime())
	cmd.RegisterModule(RegisterCounter())
	cmd.RegisterModule(RegisterIconState())
}
