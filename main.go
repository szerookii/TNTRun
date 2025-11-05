package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Seyz123/tntrun/game"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
)

// main ...
func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	serverConfig, err := readConfig(log)
	if err != nil {
		slog.Error("error reading server config file", "error", err)
		os.Exit(1)
	}

	serv := serverConfig.New()
	serv.CloseOnProgramEnd()

	serv.Listen()

	tntrun := game.NewTNTRun(serv)
	for p := range serv.Accept() {
		tntrun.OnJoin(p)
	}
}

// readConfig ...
func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	c.Server.DisableJoinQuitMessages = true
	c.Players.SaveData = false
	c.Resources.Required = false

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
