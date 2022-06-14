package command

import (
	"github.com/Seyz123/tntrun/game/config"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math"
)

type LobbyRunnable struct{}

func (LobbyRunnable) Run(src cmd.Source, o *cmd.Output) {
	if p, ok := src.(*player.Player); ok {
		pos := p.Position()
		pos = mgl64.Vec3{math.Round(pos.X()), math.Round(pos.Y()), math.Round(pos.Z())}

		err := config.UpdateConfig(true, pos)

		if err != nil {
			o.Print(text.Colourf("<red>Cannot set the lobby position!</red>"))
		} else {
			o.Printf(text.Colourf("<green>Lobby position set to X: %d, Y: %d, Z: %d</green>", int(pos.X()), int(pos.Y()), int(pos.Z())))
		}
	}
}

func (LobbyRunnable) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
