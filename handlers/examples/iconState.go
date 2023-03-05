package examples

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/handlers"
)

type IconStateHandler struct {
	Count    int
	Running  bool
	Callback func(image image.Image)
	State    bool
	Command  string
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
		command := k.IconHandlerFields["command_1"]

		icon, ok := k.IconHandlerFields["icon_1"]
		if !ok {
			return
		}

		f, err := os.Open(icon)

		if err != nil {
			log.Println(err)
			return
		}

		if !c.State {
			text = k.IconHandlerFields["text_2"]
			command = k.IconHandlerFields["command_2"]

			icon, ok := k.IconHandlerFields["icon_2"]
			if !ok {
				return
			}

			f, err = os.Open(icon)
			if err != nil {
				log.Println(err)
				return
			}

		}

		c.Command = command

		i, _, err := image.Decode(f)

		imgParsed, err := api.DrawText(i, text, k.TextSize, k.TextAlignment)

		runCommand(command)

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

func runCommand(command string) {
	go func() {
		cmd := exec.Command("/bin/sh", "-c", "/usr/bin/nohup "+command)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid:   true,
			Pgid:      0,
			Pdeathsig: syscall.SIGHUP,
		}
		if err := cmd.Start(); err != nil {
			fmt.Println("There was a problem running ", command, ":", err)
		} else {
			pid := cmd.Process.Pid
			cmd.Process.Release()
			fmt.Println(command, " has been started with pid", pid)
		}
	}()
}
