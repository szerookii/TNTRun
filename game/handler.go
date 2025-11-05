package game

import (
	"fmt"
	"math"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)

// PlayerHandler ...
type PlayerHandler struct {
	game   *TNTRun
	player *world.EntityHandle

	player.NopHandler
}

// NewPlayerHandler ...
func NewPlayerHandler(game *TNTRun, player *player.Player) *PlayerHandler {
	h := &PlayerHandler{
		game:   game,
		player: player.H(),
	}
	return h
}

// HandleMove ...
func (h *PlayerHandler) HandleMove(ctx *player.Context, newPos mgl64.Vec3, rot cube.Rotation) {
	if h.game.state != StateRunning {
		return
	}

	playerUUID := h.player.UUID()
	isInGame := false
	for _, handle := range h.game.players {
		if handle.UUID() == playerUUID {
			isInGame = true
			break
		}
	}

	if !isInGame {
		return
	}

	pos := newPos
	pos = mgl64.Vec3{math.RoundToEven(pos.X()), math.RoundToEven(pos.Y()), math.RoundToEven(pos.Z())}
	blockPos := cube.Pos{int(pos.X()), int(pos.Y()) - 1, int(pos.Z())}

	var b world.Block
	h.game.srv.World().Exec(func(tx *world.Tx) {
		b = tx.Block(blockPos)
	})

	go func() {
		time.Sleep(300 * time.Millisecond)

		if _, ok := b.(block.Air); ok {
			return
		}

		fallingBlock := entity.NewFallingBlock(world.EntitySpawnOpts{
			Position: mgl64.Vec3{float64(blockPos.X()) + 0.5, float64(blockPos.Y()) + 0.5, float64(blockPos.Z()) + 0.5},
		}, b)

		h.game.srv.World().Exec(func(gameTx *world.Tx) {
			gameTx.AddEntity(fallingBlock)
			gameTx.SetBlock(blockPos, block.Air{}, &world.SetOpts{})
		})

		go func() {
			ticker := time.NewTicker(300 * time.Millisecond)
			defer ticker.Stop()

			timeout := time.After(3 * time.Second)
			lastY := float64(blockPos.Y()) + 0.5

			for {
				select {
				case <-ticker.C:
					h.game.srv.World().Exec(func(checkTx *world.Tx) {
						if entity, exists := fallingBlock.Entity(checkTx); exists {
							currentPos := entity.Position()
							currentY := currentPos.Y()

							if math.Abs(currentY-lastY) > 0.05 {
								nextPos := currentPos.Add(mgl64.Vec3{0, -0.1, 0})
								nextBlockPos := cube.PosFromVec3(nextPos)
								belowBlock := checkTx.Block(nextBlockPos)

								shouldStop := currentY <= float64(blockPos.Y())-0.5 ||
									(!isAirBlock(belowBlock) && !isLiquidBlock(belowBlock))

								if shouldStop {
									checkTx.AddParticle(currentPos, particle.BlockBreak{Block: b})
									checkTx.RemoveEntity(entity)
									ticker.Stop()
								}
							} else if currentY <= float64(blockPos.Y())-0.5 {
								checkTx.AddParticle(currentPos, particle.BlockBreak{Block: b})
								checkTx.RemoveEntity(entity)
								ticker.Stop()
							}
							lastY = currentY
						} else {
							ticker.Stop()
						}
					})

				case <-timeout:
					h.game.srv.World().Exec(func(cleanTx *world.Tx) {
						if entity, exists := fallingBlock.Entity(cleanTx); exists {
							cleanTx.RemoveEntity(entity)
						}
					})
					return
				}
			}
		}()
	}()
}

// HandleHurt ...
func (h *PlayerHandler) HandleHurt(ctx *player.Context, _ *float64, _ bool, _ *time.Duration, src world.DamageSource) {
	if _, ok := src.(entity.VoidDamageSource); ok {
		p := ctx.Val()
		if h.game.IsPlayer(p) {
			h.game.AddSpectator(p)
		}
	}

	ctx.Cancel()
}

// HandleBlockBreak ...
func (h *PlayerHandler) HandleBlockBreak(ctx *player.Context, _ cube.Pos, _ *[]item.Stack, _ *int) {
	if h.game.config.Enabled {
		ctx.Cancel()
	}
}

// HandleBlockPlace ...
func (h *PlayerHandler) HandleBlockPlace(ctx *player.Context, _ cube.Pos, _ world.Block) {
	if h.game.config.Enabled {
		ctx.Cancel()
	}
}

// HandleFoodLoss ...
func (h *PlayerHandler) HandleFoodLoss(ctx *player.Context, _ int, _ *int) {
	if h.game.config.Enabled {
		ctx.Cancel()
	}
}

// HandleQuit ...
func (h *PlayerHandler) HandleQuit(p *player.Player) {
	if h.game.IsPlayer(p) && (h.game.state == StateIdle || h.game.state == StateStarting || h.game.state == StateRunning) {
		h.game.BroadcastMessage(fmt.Sprintf("§e%s §7left the game. §7(§e%d§7/§e%d§7)", p.Name(), len(h.game.players), MaxPlayers), TypeMessage)

		if h.game.state == StateRunning && h.game.IsPlayer(p) {
			h.game.AddSpectator(p)
		} else {
			h.game.RemovePlayer(p)
		}
	}
}

// isAirBlock checks if the block is air
func isAirBlock(b world.Block) bool {
	_, ok := b.(block.Air)
	return ok
}

// isLiquidBlock checks if the block is a liquid (water or lava)
func isLiquidBlock(b world.Block) bool {
	switch b.(type) {
	case block.Water, block.Lava:
		return true
	default:
		return false
	}
}
