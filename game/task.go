package game

import (
	"fmt"
	"time"
)

// TNTRunTask ...
type TNTRunTask struct {
	game  *TNTRun
	timer int
}

// NewTNTRunTask ...
func NewTNTRunTask(game *TNTRun) *TNTRunTask {
	return &TNTRunTask{
		game:  game,
		timer: StartTimer,
	}
}

// Start ...
func (t *TNTRunTask) Start() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			switch t.game.state {
			case StateIdle:
				if len(t.game.players) < NeededPlayers {
					t.game.BroadcastMessage(fmt.Sprintf("§eWaiting for %d players", NeededPlayers-len(t.game.players)), TypePopup)
				}

			case StateStarting:
				if t.timer > 0 {
					if t.timer <= 5 {
						t.game.BroadcastMessage(fmt.Sprintf("§e%d", t.timer), TypeTitle)
						// TODO: Add sound support later
					}

					t.game.BroadcastMessage(fmt.Sprintf("§eStarting in %d...", t.timer), TypePopup)
					t.timer--
				} else {
					t.game.BroadcastMessage("§eGame Started!", TypeTitle)
					t.game.state = StateRunning
				}

			case StateRunning:
				t.game.BroadcastMessage(fmt.Sprintf("§e%d players remaining", len(t.game.players)), TypePopup)

			case StateRestarting:
				if t.timer > 0 {
					t.game.BroadcastMessage(fmt.Sprintf("§eRestarting in %d...", t.timer), TypePopup)
					t.timer--
				} else {
					_ = t.game.srv.Close()
				}

			}
		}
	}()
}
