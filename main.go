package main

import (
	_interface "github.com/polpettone/streamdeckd/cmd/interface"
)

var engine *_interface.Engine

func main() {
	engine = _interface.NewEngine()
	engine.Run()
}
