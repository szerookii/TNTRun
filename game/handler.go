package game

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"time"
)

// PlayerHandler ...
type PlayerHandler struct {
	game   *TNTRun
	player *player.Player

	player.NopHandler
}

// NewPlayerHandler ...
func NewPlayerHandler(game *TNTRun, player *player.Player) *PlayerHandler {
	h := &PlayerHandler{
		game:   game,
		player: player,
	}
	return h
}

// HandleMove ...
func (h *PlayerHandler) HandleMove(_ *event.Context, _ mgl64.Vec3, _ float64, _ float64) {
	if h.game.state == StateRunning && h.game.IsPlayer(h.player) {
		pos := h.player.Position()
		pos = mgl64.Vec3{math.RoundToEven(pos.X()), math.RoundToEven(pos.Y()), math.RoundToEven(pos.Z())}
		b := h.player.World().Block(cube.Pos{int(pos.X()), int(pos.Y()) - 1, int(pos.Z())})

		go func() {
			<-time.After(300 * time.Millisecond)

			if _, ok := b.(block.Air); !ok {
				h.player.World().SetBlock(cube.Pos{int(pos.X()), int(pos.Y()) - 1, int(pos.Z())}, block.Air{}, &world.SetOpts{})
			}
		}()
	}
}

// HandleHurt ...
func (h *PlayerHandler) HandleHurt(ctx *event.Context, _ *float64, _ *time.Duration, src world.DamageSource) {
	if _, ok := src.(entity.VoidDamageSource); ok {
		h.game.AddSpectator(h.player)
	}

	ctx.Cancel()
}

// HandleBlockBreak ...
func (h *PlayerHandler) HandleBlockBreak(ctx *event.Context, _ cube.Pos, _ *[]item.Stack) {
	if h.game.config.Enabled {
		ctx.Cancel()
	}
}

// HandleBlockPlace ...
func (h *PlayerHandler) HandleBlockPlace(ctx *event.Context, _ cube.Pos, _ world.Block) {
	if h.game.config.Enabled {
		ctx.Cancel()
	}
}

// HandleFoodLoss ...
func (h *PlayerHandler) HandleFoodLoss(ctx *event.Context, _ int, _ int) {
	if h.game.config.Enabled {
		ctx.Cancel()
	}
}

// HandleQuit ...
func (h *PlayerHandler) HandleQuit() {
	if h.game.IsPlayer(h.player) && h.game.state == StateIdle || h.game.state == StateStarting || h.game.state == StateRunning {
		h.game.BroadcastMessage(fmt.Sprintf("§e%s §7left the game. §7(§e%d§7/§e%d§7)", h.player.Name(), len(h.game.players), MaxPlayers), TypeMessage)

		if h.game.state == StateRunning && h.game.IsPlayer(h.player) {
			h.game.AddSpectator(h.player)
		} else {
			h.game.RemovePlayer(h.player)
		}
	}
}
