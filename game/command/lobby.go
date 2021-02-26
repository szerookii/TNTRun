package command

import (
	"github.com/Seyz123/tntrun/game/config"
	"github.com/df-mc/dragonfly/dragonfly/cmd"
	"github.com/df-mc/dragonfly/dragonfly/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math"
	"reflect"
)

type LobbbyRunnable struct{
	Sub lobby
}

func (r *LobbbyRunnable) Run(source cmd.Source, output *cmd.Output) {
	if p, ok := source.(*player.Player); ok {
		pos :=  p.Position()
		pos = mgl64.Vec3{math.Round(pos.X()), math.Round(pos.Y()), math.Round(pos.Z())}

		err := config.UpdateConfig(true, pos)

		if err != nil {
			output.Print(text.Colourf("<red>Cannot set the lobby position!</red>"))
		} else {
			output.Printf(text.Colourf("<green>Lobby position set to X: %d, Y: %d, Z: %d</green>", int(pos.X()), int(pos.Y()), int(pos.Z())))
		}
	}
}

type lobby string

func (lobby) Type() string {
	return "lobby"
}

func (lobby) Options() []string {
	return []string{"lobby"}
}

func (lobby) SetOption(string, reflect.Value) {}
