package _interface

var engine *Engine

func StartEngine(configPath string) {
	engine = NewEngine(configPath)
	engine.Run()
}
