package command

import (
	"math"

	"github.com/Seyz123/tntrun/game/config"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// LobbyRunnable ...
type LobbyRunnable struct{}

// Run ...
func (LobbyRunnable) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		pos := p.Position()
		pos = mgl64.Vec3{math.Round(pos.X()), math.Round(pos.Y()), math.Round(pos.Z())}

		err := config.UpdateConfig(true, pos)

		if err != nil {
			o.Print(text.Colourf("<red>Cannot set the lobby position!</red>"))
		} else {
			o.Print(text.Colourf("<green>Lobby position set to X: %d, Y: %d, Z: %d</green>", int(pos.X()), int(pos.Y()), int(pos.Z())))
		}
	}
}

// Allow ...
func (LobbyRunnable) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
