package _interface

import (
	"context"
	"image"
	"image/gif"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/unix-streamdeck/api"
	"golang.org/x/sync/semaphore"
)

type GifIconHandler struct {
	Running bool
	Lock    *semaphore.Weighted
	Quit    chan bool
}

func (s *GifIconHandler) Start(key api.Key, info api.StreamDeckInfo, callback func(image image.Image)) {
	if s.Quit == nil {
		s.Quit = make(chan bool)
	}
	if s.Lock == nil {
		s.Lock = semaphore.NewWeighted(1)
	}
	s.Running = true
	icon, ok := key.IconHandlerFields["icon"]
	if !ok {
		return
	}
	f, err := os.Open(icon)
	if err != nil {
		log.Println(err)
		return
	}
	gifs, err := gif.DecodeAll(f)
	if err != nil {
		log.Println(err)
		return
	}
	timeDelay := gifs.Delay[0]
	frames := make([]image.Image, len(gifs.Image))
	for i, frame := range gifs.Image {
		img := api.ResizeImage(frame, info.IconSize)
		if key.IconHandlerFields["text"] != "" {
			size, _ := strconv.ParseInt(key.IconHandlerFields["text_size"], 10, 0)
			img, err = api.DrawText(img, key.IconHandlerFields["text"], int(size), key.IconHandlerFields["text_alignment"])
			if err != nil {
				log.Println(err)
			}
		}
		frames[i] = img
	}
	go s.loop(frames, timeDelay, callback)
}

func (s *GifIconHandler) IsRunning() bool {
	return s.Running
}

func (s *GifIconHandler) SetRunning(running bool) {
	s.Running = running
}

func (s *GifIconHandler) Stop() {
	s.Running = false
	s.Quit <- true
}

func (s *GifIconHandler) loop(frames []image.Image, timeDelay int, callback func(image image.Image)) {
	ctx := context.Background()
	err := s.Lock.Acquire(ctx, 1)
	if err != nil {
		return
	}
	defer s.Lock.Release(1)
	gifIndex := 0
	for {
		select {
		case <-s.Quit:
			return
		default:
			img := frames[gifIndex]
			callback(img)
			gifIndex++
			if gifIndex >= len(frames) {
				gifIndex = 0
			}
			time.Sleep(time.Duration(timeDelay * 10000000))
		}
	}
}

func RegisterGif() Module {
	return Module{NewIcon: func() api.IconHandler {
		return &GifIconHandler{Running: true, Lock: semaphore.NewWeighted(1)}
	}, Name: "Gif", IconFields: []api.Field{{Title: "Icon", Name: "icon", Type: "File", FileTypes: []string{".gif"}}, {Title: "Text", Name: "text", Type: "Text"}, {Title: "Text Size", Name: "text_size", Type: "Number"}, {Title: "Text Alignment", Name: "text_alignment", Type: "TextAlignment"}}}
}
