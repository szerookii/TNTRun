package game

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
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
		b := h.player.World().Block(cube.Pos{int(pos.X()), int(pos.Y()) - 1, int(pos.Z())})

		go func() {
			<-time.After(300 * time.Millisecond)

			if _, ok := b.(block.Air); !ok {
				h.player.World().SetBlock(cube.Pos{int(pos.X()), int(pos.Y()) - 1, int(pos.Z())}, block.Air{}, &world.SetOpts{})
			}
		}()
	}
}

func (h *PlayerHandler) HandleHurt(ctx *event.Context, dmg *float64, attackImmunity *time.Duration, src damage.Source) {
	if _, ok := src.(damage.SourceVoid); ok {
		h.game.AddSpectator(h.player)
	}

	ctx.Cancel()
}

func (h *PlayerHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	if h.game.config.Enabled {
		ctx.Cancel()
	}
}

func (h *PlayerHandler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	if h.game.config.Enabled {
		ctx.Cancel()
	}
}

func (h *PlayerHandler) HandleFoodLoss(ctx *event.Context, _ int, _ int) {
	if h.game.config.Enabled {
		ctx.Cancel()
	}
}

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
