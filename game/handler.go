package game

import (
	"fmt"
	"github.com/df-mc/dragonfly/dragonfly/block"
	"github.com/df-mc/dragonfly/dragonfly/entity/damage"
	"github.com/df-mc/dragonfly/dragonfly/event"
	"github.com/df-mc/dragonfly/dragonfly/player"
	"github.com/df-mc/dragonfly/dragonfly/world"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"time"
)

type PlayerHandler struct {
	player.NopHandler
	game   *TNTRun
	player *player.Player
}

func NewPlayerHandler(game *TNTRun, player *player.Player) *PlayerHandler {
	h := &PlayerHandler{
		game:   game,
		player: player,
	}

	return h
}

func (h *PlayerHandler) HandleMove(ctx *event.Context, _ mgl64.Vec3, _ float64, _ float64) {
	if h.game.state == StateRunning && h.game.IsPlayer(h.player) {
		pos := h.player.Position()
		pos = mgl64.Vec3{math.RoundToEven(pos.X()), math.RoundToEven(pos.Y()), math.RoundToEven(pos.Z())}
		b := h.player.World().Block(world.BlockPos{int(pos.X()), int(pos.Y()) - 1, int(pos.Z())})

		go func() {
			<-time.After(300 * time.Millisecond)

			if _, ok := b.(block.Air); !ok {
				h.player.World().SetBlock(world.BlockPos{int(pos.X()), int(pos.Y()) - 1, int(pos.Z())}, block.Air{})
			}
		}()
	}
}

func (h *PlayerHandler) HandleHurt(ctx *event.Context, _ *float64, source damage.Source) {
	if _, ok := source.(damage.SourceVoid); ok {
		h.game.AddSpectator(h.player)
	}

	ctx.Cancel()
}

func (h *PlayerHandler) HandleQuit() {
	if h.game.state == StateIdle || h.game.state == StateStarting || h.game.state == StateRunning {
		h.game.BroadcastMessage(fmt.Sprintf("§e%s §7left the game. §7(§e%d§7/§e%d§7)", h.player.Name(), len(h.game.players), MaxPlayers), TypeMessage)

		if h.game.state == StateRunning && h.game.IsPlayer(h.player) {
			h.game.AddSpectator(h.player)
		}
	}
}
