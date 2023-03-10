package _interface

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/polpettone/streamdeckd/cmd"
	"github.com/polpettone/streamdeckd/cmd/myConfig"
	"github.com/unix-streamdeck/api"
)

type IconStateHandler struct {
	Count        int
	Running      bool
	Callback     func(image image.Image)
	State        bool
	Icon1        image.Image
	Icon2        image.Image
	CurrentImage image.Image
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
		c.CurrentImage = img
		c.LoadIcons(k)

		draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)

		text := k.IconHandlerFields["text_2"]
		command := k.IconHandlerFields["command_1"]

		if c.Icon1 != nil {
			c.CurrentImage = c.Icon1
		}

		if !c.State {
			text = k.IconHandlerFields["text_1"]
			command = k.IconHandlerFields["command_2"]
			if c.Icon2 != nil {
				c.CurrentImage = c.Icon2
			}
		}

		imgParsed, err := api.DrawText(c.CurrentImage, text, k.TextSize, k.TextAlignment)

		runCommand(command)

		key := myConfig.MyConfig.Decks[0].Pages[8][7]
		devs := myConfig.MyDevs

		log.Printf("key: %v", key)

		for k, _ := range devs {
			log.Printf("devs: %s", k)
		}

		SetImage(devs["DL49K1A67580"], imgParsed, 7, 8)

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

func RegisterIconState() cmd.Module {
	return cmd.Module{
		NewIcon: func() api.IconHandler {
			return &IconStateHandler{Running: true, Count: 0}

		}, NewKey: func() api.KeyHandler {
			return &IconStateKeyHandler{}
		}, Name: "IconState"}
}

// Both Icon Images must be loaded successfully
// If one of both is missing, no image will be loaded
func (c *IconStateHandler) LoadIcons(k api.Key) {

	iconPath1, ok := k.IconHandlerFields["icon_1"]
	if !ok {
		return
	}

	iconFile1, err := os.Open(iconPath1)
	if err != nil {
		log.Println(err)
		return
	}

	image1, _, err := image.Decode(iconFile1)
	if err != nil {
		log.Println(err)
		return
	}

	c.Icon1 = image1

	iconPath2, ok := k.IconHandlerFields["icon_2"]
	if !ok {
		return
	}

	iconFile2, err := os.Open(iconPath2)

	if err != nil {
		log.Println(err)
		return
	}

	image2, _, err := image.Decode(iconFile2)
	if err != nil {
		log.Println(err)
		return
	}

	c.Icon2 = image2
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
