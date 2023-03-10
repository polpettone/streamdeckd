package examples

import "github.com/polpettone/streamdeckd/handlers"

func RegisterBaseModules() {
	handlers.RegisterModule(RegisterGif())
	handlers.RegisterModule(RegisterTime())
	handlers.RegisterModule(RegisterCounter())
	handlers.RegisterModule(RegisterSpotify())
	handlers.RegisterModule(RegisterIconState())
}
