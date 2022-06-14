package main

import (
	"fmt"
	"github.com/Seyz123/tntrun/game"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	serverConfig, err := readConfig()
	if err != nil {
		log.Fatalf("error reading server config file: %v", err)
	}

	serv := server.New(&serverConfig, log)
	serv.CloseOnProgramEnd()
	if err := serv.Start(); err != nil {
		log.Fatalln(err)
	}

	tntrun := game.NewTNTRun(serv)

	for serv.Accept(tntrun.OnJoin) {

	}
}

func readConfig() (server.Config, error) {
	c := server.DefaultConfig()
	c.Server.JoinMessage = ""
	c.Server.QuitMessage = ""

	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}
