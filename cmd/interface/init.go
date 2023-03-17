package _interface

func RegisterBaseModules() {
	RegisterModule(RegisterGif())
	RegisterModule(RegisterTime())
	RegisterModule(RegisterCounter())
	RegisterModule(RegisterIconState())
	RegisterModule(RegisterGame())
}
