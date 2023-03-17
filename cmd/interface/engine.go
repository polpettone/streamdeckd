package _interface

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"image"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/polpettone/streamdeckd/cmd/models"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/unix-streamdeck/api"
	streamdeck "github.com/unix-streamdeck/driver"
	"golang.org/x/sync/semaphore"
)

type Engine struct {
	devs          map[string]*models.VirtualDev
	config        *api.Config
	configPath    string
	disconnectSem *semaphore.Weighted
	connectSem    *semaphore.Weighted
	basicConfig   api.Config
	isRunning     bool
}

func (engine *Engine) SetImage(img image.Image, i int, page int) {
	SetImage(engine, engine.devs["CL33L2A02177"], img, i, page)
}

func NewEngine() *Engine {

	return &Engine{

		config:     nil,
		configPath: "",

		devs:      make(map[string]*models.VirtualDev),
		isRunning: true,

		disconnectSem: semaphore.NewWeighted(1),
		connectSem:    semaphore.NewWeighted(2),

		basicConfig: api.Config{
			Modules: []string{},
			Decks: []api.Deck{
				{},
			},
		},
	}

}

func (engine *Engine) Run() {
	checkOtherRunningInstances()

	configPtr := flag.String("config", engine.configPath, "Path to config file")

	flag.Parse()

	if *configPtr != "" {
		engine.configPath = *configPtr
	} else {

		basePath := os.Getenv("HOME") + string(os.PathSeparator) + ".config"
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			basePath = os.Getenv("XDG_CONFIG_HOME")
		}

		engine.configPath = basePath + string(os.PathSeparator) + ".streamdeck-config.json"
	}

	engine.cleanupHook()

	go InitDBUS()

	RegisterBaseModules(engine)

	engine.loadConfig()
	engine.attemptConnection()
}

func checkOtherRunningInstances() {
	processes, err := process.Processes()
	if err != nil {
		log.Println("Could not check for other instances of streamdeckd, assuming no others running: %s", err)
	}
	for _, proc := range processes {
		name, err := proc.Name()
		if err == nil &&
			name == "streamdeckd" &&
			int(proc.Pid) != os.Getpid() {
			log.Fatalln("Another instance of streamdeckd is already running, exiting...")
		}
	}
}

func (engine *Engine) attemptConnection() {

	for engine.isRunning {

		dev := &models.VirtualDev{}
		dev, _ = engine.openDevice()

		if dev.IsOpen {

			SetPage(engine, dev, dev.Page)
			found := false

			for i := range sDInfo {
				if sDInfo[i].Serial == dev.Deck.Serial {
					found = true
				}
			}

			if !found {
				sDInfo = append(sDInfo, api.StreamDeckInfo{
					Cols:     int(dev.Deck.Columns),
					Rows:     int(dev.Deck.Rows),
					IconSize: int(dev.Deck.Pixels),
					Page:     0,
					Serial:   dev.Deck.Serial,
				})
			}

			go Listen(dev, engine)

		}
		time.Sleep(250 * time.Millisecond)
	}
}

func (engine *Engine) Disconnect(dev *models.VirtualDev) {
	ctx := context.Background()
	err := engine.disconnectSem.Acquire(ctx, 1)
	if err != nil {
		return
	}
	defer engine.disconnectSem.Release(1)
	if !dev.IsOpen {
		return
	}
	log.Println("Device (" + dev.Deck.Serial + ") disconnected")
	_ = dev.Deck.Close()
	dev.IsOpen = false
	unmountDevHandlers(dev)
}

func (engine *Engine) openDevice() (*models.VirtualDev, error) {
	ctx := context.Background()
	err := engine.connectSem.Acquire(ctx, 1)
	if err != nil {
		return &models.VirtualDev{}, err
	}
	defer engine.connectSem.Release(1)
	d, err := streamdeck.Devices()
	if err != nil {
		return &models.VirtualDev{}, err
	}
	if len(d) == 0 {
		return &models.VirtualDev{}, errors.New("No streamdeck devices found")
	}
	device := streamdeck.Device{Serial: ""}
	for i := range d {
		found := false
		for s := range engine.devs {
			if d[i].ID == engine.devs[s].Deck.ID && engine.devs[s].IsOpen {
				found = true
				break
			} else if d[i].Serial == s && !engine.devs[s].IsOpen {
				err = d[i].Open()
				if err != nil {
					return &models.VirtualDev{}, err
				}
				engine.devs[s].Deck = d[i]
				engine.devs[s].IsOpen = true
				return engine.devs[s], nil
			}
		}
		if !found {
			device = d[i]
		}
	}
	if len(device.Serial) != 12 {
		return &models.VirtualDev{}, errors.New("No streamdeck devices found")
	}
	err = device.Open()
	if err != nil {
		return &models.VirtualDev{}, err
	}
	devNo := -1

	for i := range engine.config.Decks {
		if engine.config.Decks[i].Serial == device.Serial {
			devNo = i
		}
	}
	if devNo == -1 {
		var pages []api.Page
		page := api.Page{}
		for i := 0; i < int(device.Rows)*int(device.Columns); i++ {
			page = append(page, api.Key{})
		}
		pages = append(pages, page)
		engine.config.Decks = append(engine.config.Decks, api.Deck{Serial: device.Serial, Pages: pages})
		devNo = len(engine.config.Decks) - 1
	}
	dev := &models.VirtualDev{Deck: device, Page: 0, IsOpen: true, Config: engine.config.Decks[devNo].Pages}
	engine.devs[device.Serial] = dev
	log.Println("Device (" + device.Serial + ") connected")

	return dev, nil
}

func (engine *Engine) loadConfig() {

	var err error

	engine.config, err = engine.ReadConfig()

	if err != nil && !os.IsNotExist(err) {
		log.Println(err)
	} else if os.IsNotExist(err) {

		file, err := os.Create(engine.configPath)
		if err != nil {
			log.Println(err)
		}

		err = file.Close()
		if err != nil {
			log.Println(err)
		}

		engine.config = &engine.basicConfig
		err = engine.SaveConfig()
		if err != nil {
			log.Println(err)
		}

	}
	if len(engine.config.Modules) > 0 {

		for _, module := range engine.config.Modules {
			LoadModule(module)
		}
	}
}

func (engine *Engine) ReadConfig() (*api.Config, error) {
	data, err := os.ReadFile(engine.configPath)
	if err != nil {
		return &api.Config{}, err
	}
	var config api.Config
	err = json.Unmarshal(data, &config)
	return &config, nil
}

func (engine *Engine) cleanupHook() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs,
		syscall.SIGSTOP,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT)

	go func() {
		<-sigs

		log.Println("Cleaning up")

		engine.isRunning = false

		engine.unmountHandlers()

		var err error

		for _, v := range engine.devs {

			if v.IsOpen {
				err = v.Deck.Reset()
				if err != nil {
					log.Println(err)
				}

				err = v.Deck.Close()
				if err != nil {
					log.Println(err)
				}
			}
		}
		os.Exit(0)
	}()
}

func (engine *Engine) SetConfig(configString string) error {
	engine.unmountHandlers()
	var err error
	engine.config = nil
	err = json.Unmarshal([]byte(configString), &engine.config)
	if err != nil {
		return err
	}
	for _, dev := range engine.devs {
		for i := range engine.config.Decks {
			if dev.Deck.Serial == engine.config.Decks[i].Serial {
				dev.Config = engine.config.Decks[i].Pages
			}
		}
		SetPage(engine, dev, dev.Page)
	}
	return nil
}

func (engine *Engine) ReloadConfig() error {
	engine.unmountHandlers()
	engine.loadConfig()
	for _, dev := range engine.devs {
		for i := range engine.config.Decks {
			if dev.Deck.Serial == engine.config.Decks[i].Serial {
				dev.Config = engine.config.Decks[i].Pages
			}
		}
		SetPage(engine, dev, dev.Page)
	}
	return nil
}

func (engine *Engine) SaveConfig() error {
	f, err := os.OpenFile(engine.configPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	var configString []byte
	configString, err = json.Marshal(engine.config)
	if err != nil {
		return err
	}
	_, err = f.Write(configString)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}
func (engine *Engine) unmountHandlers() {
	for _, dev := range engine.devs {
		unmountDevHandlers(dev)
	}
}

func unmountDevHandlers(dev *models.VirtualDev) {
	for i := range dev.Config {
		UnmountPageHandlers(dev.Config[i])
	}
}

func UnmountPageHandlers(page api.Page) {
	for i2 := 0; i2 < len(page); i2++ {
		key := &page[i2]
		if key.IconHandlerStruct != nil {
			log.Printf("Stopping %s\n", key.IconHandler)
			if key.IconHandlerStruct.IsRunning() {
				go func() {
					key.IconHandlerStruct.Stop()
					log.Printf("Stopped %s\n", key.IconHandler)
				}()
			}
		}
	}
}
