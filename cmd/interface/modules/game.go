package modules

import (
	"github.com/polpettone/streamdeckd/cmd/models"
	"github.com/unix-streamdeck/api"
	"image"
	"image/draw"
	"log"
	"strconv"
)

type GameHandler struct {
	Running      bool
	Callback     func(image image.Image)
	CurrentImage image.Image
}

func NewGameHandler() *GameHandler {
	return &GameHandler{
		Running: true,
	}
}

type GameState struct {
	solution       int
	pressedNumbers []int
}

func NewGameState(solution int) *GameState {
	return &GameState{
		pressedNumbers: []int{},
		solution:       solution,
	}
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

	text, v := k.IconHandlerFields["number"]
	if !v {
		text = "the game"
	}

	imgParsed, err := api.DrawText(
		g.CurrentImage, text, k.TextSize, k.TextAlignment)

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

type GameKeyHandler struct {
	Action    models.Action
	GameState *GameState
}

func NewGameKeyHandler(action models.Action, state *GameState) *GameKeyHandler {
	return &GameKeyHandler{
		Action:    action,
		GameState: state,
	}
}
func (g GameKeyHandler) Key(key api.Key, info api.StreamDeckInfo) {
	handler := key.IconHandlerStruct.(*GameHandler)

	imgParsed, _ := api.DrawText(handler.CurrentImage, "foo", key.TextSize, key.TextAlignment)

	numberText := key.IconHandlerFields["number"]
	number, err := strconv.Atoi(numberText)
	if err == nil {
		log.Printf("Game Number %d pressed", number)
		g.GameState.pressedNumbers = append(g.GameState.pressedNumbers, number)
		if number == g.GameState.solution {
			g.Action.SetImage(imgParsed, 14, 5)
		}
		log.Printf("%v", g.GameState.pressedNumbers)
	} else {
		log.Println(err)
	}

	if handler.Callback != nil {
		handler.Start(key, info, handler.Callback)
	}
}
