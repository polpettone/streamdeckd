package _interface

import (
	"github.com/polpettone/streamdeckd/cmd/models"
	"github.com/unix-streamdeck/api"
	"image"
	"image/draw"
	"log"
)

type GameHandler struct {
	Running      bool
	Callback     func(image image.Image)
	CurrentImage image.Image
}

func (g *GameHandler) Start(
	k api.Key,
	info api.StreamDeckInfo,
	callback func(image image.Image)) {

	if g.Callback == nil {
		g.Callback = callback
	}

	img := image.NewRGBA(
		image.Rect(0, 0, info.IconSize, info.IconSize))
	g.CurrentImage = img

	draw.Draw(
		img,
		img.Bounds(),
		image.Black,
		image.ZP,
		draw.Src)

	imgParsed, err := api.DrawText(
		g.CurrentImage, "The Game", k.TextSize, k.TextAlignment)

	if err != nil {
		log.Println(err)
	} else {
		callback(imgParsed)
	}
}

func (g *GameHandler) IsRunning() bool {
	return g.Running
}

func (g *GameHandler) SetRunning(running bool) {
	g.Running = running
}

func (g *GameHandler) Stop() {
	g.Running = false
}

type GameKeyHandler struct{}

func (GameKeyHandler) Key(key api.Key, info api.StreamDeckInfo) {
	handler := key.IconHandlerStruct.(*GameHandler)

	imgParsed, _ := api.DrawText(handler.CurrentImage, "foo", key.TextSize, key.TextAlignment)

	SetImage(currentEngine,
		currentEngine.devs["CL33L2A02177"],
		imgParsed,
		10,
		5)

	if handler.Callback != nil {
		handler.Start(key, info, handler.Callback)
	}
}

func RegisterGame() models.Module {
	return models.Module{
		NewIcon: func() api.IconHandler {
			return &GameHandler{Running: true}
		},
		NewKey: func() api.KeyHandler {
			return &GameKeyHandler{}
		},
		Name: "Game",
	}
}
