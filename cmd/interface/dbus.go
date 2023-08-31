package _interface

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/godbus/dbus/v5"
	"github.com/polpettone/streamdeckd/cmd/models"
	"github.com/unix-streamdeck/api"
)

var conn *dbus.Conn

var sDbus *StreamDeckDBus
var sDInfo []api.StreamDeckInfo

func NewStreamDeckBus(engine *Engine) *StreamDeckDBus {
	return &StreamDeckDBus{
		engine: engine,
	}
}

type StreamDeckDBus struct {
	engine *Engine
}

func (s StreamDeckDBus) GetDeckInfo() (string, *dbus.Error) {
	infoString, err := json.Marshal(sDInfo)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return string(infoString), nil
}

func (s StreamDeckDBus) GetConfig() (string, *dbus.Error) {
	configString, err := json.Marshal(s.engine.config)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return string(configString), nil
}

func (s StreamDeckDBus) ReloadConfig() *dbus.Error {
	err := s.engine.ReloadConfig()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (s StreamDeckDBus) SetPage(serial string, page int) *dbus.Error {
	for _, dev := range s.engine.devs {
		if dev.Deck.Serial == serial {
			SetPage(s.engine, dev, page)
			return nil
		}
	}
	return dbus.MakeFailedError(errors.New("Device with Serial: " + serial + " could not be found"))
}

func (s StreamDeckDBus) SetConfig(configString string) *dbus.Error {
	err := s.engine.SetConfig(configString)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (s StreamDeckDBus) CommitConfig() *dbus.Error {
	err := s.engine.SaveConfig()
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (StreamDeckDBus) GetModules() (string, *dbus.Error) {
	var modules []api.Module
	for _, module := range AvailableModules() {
		modules = append(modules, api.Module{Name: module.Name, IconFields: module.IconFields, KeyFields: module.KeyFields, IsIcon: module.NewIcon != nil, IsKey: module.NewKey != nil})
	}
	modulesString, err := json.Marshal(modules)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return string(modulesString), nil
}

func (s StreamDeckDBus) PressButton(serial string, keyIndex int) *dbus.Error {
	dev, ok := s.engine.devs[serial]
	if !ok || !dev.IsOpen {
		return dbus.MakeFailedError(errors.New("Can't find connected device: " + serial))
	}
	HandleInput(s.engine, dev, &dev.Config[dev.Page][keyIndex], dev.Page)
	return nil
}

func InitDBUS() error {
	var err error
	conn, err = dbus.SessionBus()
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close()

	sDbus = &StreamDeckDBus{}
	conn.ExportAll(sDbus, "/com/unixstreamdeck/streamdeckd", "com.unixstreamdeck.streamdeckd")
	reply, err := conn.RequestName("com.unixstreamdeck.streamdeckd",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Println(err)
		return err
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return errors.New("DBus: Name already taken")
	}
	select {}
}

func EmitPage(dev *models.VirtualDev, page int) {
	if conn != nil {
		conn.Emit("/com/unixstreamdeck/streamdeckd", "com.unixstreamdeck.streamdeckd.Page", dev.Deck.Serial, page)
	}
	for i := range sDInfo {
		if sDInfo[i].Serial == dev.Deck.Serial {
			sDInfo[i].Page = page
		}
	}
}
