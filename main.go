package main

import (
	"fmt"
	"github.com/Seyz123/tntrun/game"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"os"
)

// main ...
func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	serverConfig, err := readConfig(log)
	if err != nil {
		log.Fatalf("error reading server config file: %v", err)
	}

	serv := serverConfig.New()
	serv.CloseOnProgramEnd()

	serv.Listen()

	tntrun := game.NewTNTRun(serv)
	for serv.Accept(tntrun.OnJoin) {
	}
}

// readConfig ...
func readConfig(log server.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	c.Server.JoinMessage = ""
	c.Server.QuitMessage = ""

	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("failed creating config: %v", err)
		}
		return zero, nil
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("error decoding config: %v", err)
	}
	return c.Config(log)
}
