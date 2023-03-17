package _interface

import (
	modules2 "github.com/polpettone/streamdeckd/cmd/interface/modules"
	"github.com/polpettone/streamdeckd/cmd/models"
	"github.com/unix-streamdeck/api"
	"golang.org/x/sync/semaphore"
)

func RegisterBaseModules(action models.Action) {
	RegisterModule(RegisterGif())
	RegisterModule(RegisterTime())
	RegisterModule(RegisterCounter())
	RegisterModule(RegisterIconState())
	RegisterModule(RegisterGame(action))
}

func RegisterIconState() models.Module {
	return models.Module{

		NewIcon: func() api.IconHandler {
			return &modules2.IconStateHandler{Running: true}
		},
		NewKey: func() api.KeyHandler {
			return &modules2.IconStateKeyHandler{}
		},
		Name: "IconState"}
}

func RegisterGif() models.Module {
	return models.Module{NewIcon: func() api.IconHandler {
		return &modules2.GifIconHandler{Running: true, Lock: semaphore.NewWeighted(1)}
	}, Name: "Gif", IconFields: []api.Field{{Title: "Icon", Name: "icon", Type: "File", FileTypes: []string{".gif"}}, {Title: "Text", Name: "text", Type: "Text"}, {Title: "Text Size", Name: "text_size", Type: "Number"}, {Title: "Text Alignment", Name: "text_alignment", Type: "TextAlignment"}}}
}

func RegisterTime() models.Module {
	return models.Module{NewIcon: func() api.IconHandler {
		return &modules2.TimeIconHandler{Running: true}
	}, Name: "Time"}
}

func RegisterCounter() models.Module {
	return models.Module{NewIcon: func() api.IconHandler {
		return &modules2.CounterIconHandler{Running: true, Count: 0}
	}, NewKey: func() api.KeyHandler {
		return &modules2.CounterKeyHandler{}
	}, Name: "Counter"}
}

func RegisterGame(action models.Action) models.Module {
	return models.Module{
		NewIcon: func() api.IconHandler {
			return &modules2.GameHandler{Running: true}
		},
		NewKey: func() api.KeyHandler {
			return &modules2.GameKeyHandler{Action: action}
		},
		Name: "Game",
	}
}
