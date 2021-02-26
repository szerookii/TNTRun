package game

import (
	"fmt"
	"github.com/df-mc/dragonfly/dragonfly/player/title"
	"github.com/df-mc/dragonfly/dragonfly/world/sound"
	"time"
)

type TNTRunTask struct {
	game  *TNTRun
	timer int
}

func NewTNTRunTask(game *TNTRun) *TNTRunTask {
	return &TNTRunTask{
		game:  game,
		timer: StartTimer,
	}
}

func (t *TNTRunTask) Start() {
	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				switch t.game.state {
				case StateIdle:
					for _, p := range t.game.players {
						p.SendPopup(fmt.Sprintf("§eWaiting for %d players", NeededPlayers-len(t.game.players)))
					}

				case StateStarting:
					if t.timer > 0 {
						for _, p := range t.game.players {
							if t.timer <= 5 {
								p.SendTitle(title.New(fmt.Sprintf("§e%d", t.timer)))
								p.PlaySound(sound.Click{})
							}

							p.SendPopup(fmt.Sprintf("§eStarting in %d...", t.timer))
						}

						t.timer--
					} else {
						for _, p := range t.game.players {
							p.SendTitle(title.New("§eGame Started!"))
							p.PlaySound(sound.Explosion{})
						}

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
		}
	}()
}
