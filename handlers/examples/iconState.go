package examples

import (
	"image"
	"image/draw"
	"log"

	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/handlers"
)

type IconStateHandler struct {
	Count    int
	Running  bool
	Callback func(image image.Image)
	State    bool
}

func (c *IconStateHandler) Start(
	k api.Key,
	info api.StreamDeckInfo,
	callback func(image image.Image)) {

	if c.Callback == nil {
		c.Callback = callback
	}

	if c.Running {
		img := image.NewRGBA(image.Rect(0, 0, info.IconSize, info.IconSize))

		draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)

		text := k.IconHandlerFields["text_1"]

		if !c.State {
			text = k.IconHandlerFields["text_2"]
		}

		imgParsed, err := api.DrawText(img, text, k.TextSize, k.TextAlignment)

		if err != nil {
			log.Println(err)
		} else {
			callback(imgParsed)
		}
	}
}

func (c *IconStateHandler) IsRunning() bool {
	return c.Running
}

func (c *IconStateHandler) SetRunning(running bool) {
	c.Running = running
}

func (c IconStateHandler) Stop() {
	c.Running = false
}

type IconStateKeyHandler struct{}

func (IconStateKeyHandler) Key(key api.Key, info api.StreamDeckInfo) {

	handler := key.IconHandlerStruct.(*IconStateHandler)
	handler.Count += 1

	if handler.State {
		handler.State = false
	} else {
		handler.State = true
	}

	if handler.Callback != nil {
		handler.Start(key, info, handler.Callback)
	}
}

func RegisterIconState() handlers.Module {
	return handlers.Module{NewIcon: func() api.IconHandler {
		return &IconStateHandler{Running: true, Count: 0}
	}, NewKey: func() api.KeyHandler {
		return &IconStateKeyHandler{}
	}, Name: "IconState"}
}
