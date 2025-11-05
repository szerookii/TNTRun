package game

import (
	"fmt"

	"github.com/Seyz123/tntrun/game/command"
	"github.com/Seyz123/tntrun/game/config"
	"github.com/Seyz123/tntrun/game/utils"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/world"
)

const (
	MaxPlayers    = 10
	NeededPlayers = 1
	StartTimer    = 5
	RestartTimer  = 10
)

const (
	TypeMessage = 0
	TypePopup   = 1
	TypeTitle   = 2
)

// TNTRun ...
type TNTRun struct {
	srv        *server.Server
	config     *config.Config
	state      int
	task       *TNTRunTask
	players    []*world.EntityHandle
	spectators []*world.EntityHandle
}

// NewTNTRun ...
func NewTNTRun(srv *server.Server) *TNTRun {
	conf, err := config.GetConfig()

	if err != nil {
		panic(err)
	}

	w := srv.World()
	w.SetDefaultGameMode(world.GameModeAdventure)
	w.SetTime(5000)
	w.StopTime()

	cmd.Register(cmd.New("tntrun", "", []string{}, &command.LobbyRunnable{}))

	game := &TNTRun{
		srv:        srv,
		config:     conf,
		state:      StateIdle,
		players:    []*world.EntityHandle{},
		spectators: []*world.EntityHandle{},
	}

	game.task = NewTNTRunTask(game)
	game.task.Start()

	return game
}

// OnJoin ...
func (t *TNTRun) OnJoin(p *player.Player) {
	if len(t.players) >= MaxPlayers {
		p.Disconnect("§cThis game is full.")
		return
	}

	if t.state != StateIdle && t.state != StateStarting {
		p.Disconnect("§cThis game is already started.")
		return
	}

	if !t.config.Enabled {
		p.SetGameMode(world.GameModeCreative)
	} else {
		t.players = append(t.players, p.H())

		p.SetGameMode(world.GameModeAdventure)
		p.Teleport(t.config.Lobby)
		p.Handle(NewPlayerHandler(t, p))

		t.BroadcastMessage(fmt.Sprintf("§e%s §7joined the game. §7(§e%d§7/§e%d§7)", p.Name(), len(t.players), MaxPlayers), TypeMessage)
	}

	if len(t.players) >= NeededPlayers {
		t.state = StateStarting
	}
}

// BroadcastMessage ...
func (t *TNTRun) BroadcastMessage(msg string, msgType int) {
	w := t.srv.World()
	w.Exec(func(tx *world.Tx) {
		playerCount := 0
		for entity := range tx.Entities() {
			if player, ok := entity.(*player.Player); ok {
				if t.IsPlayer(player) || t.IsSpectator(player) {
					playerCount++
					switch msgType {
					case TypeMessage:
						player.Message(msg)
					case TypePopup:
						player.SendPopup(msg)
					case TypeTitle:
						player.SendTitle(title.New(msg))
					}
				}
			}
		}
	})
}

// IsSpectator ...
func (t *TNTRun) IsSpectator(player *player.Player) bool {
	playerUUID := player.UUID()
	for _, handle := range t.spectators {
		if handle.UUID() == playerUUID {
			return true
		}
	}
	return false
}

// CheckWinner ...
func (t *TNTRun) CheckWinner() {
	if len(t.players) == 1 {
		winnerHandle := t.players[0]
		winner := utils.EntityHandleToEntity[player.Player](winnerHandle)
		if winner != nil {
			t.BroadcastMessage(fmt.Sprintf("§e%s §7won the game!", winner.Name()), TypeMessage)
		}

		t.task.timer = RestartTimer
		t.state = StateRestarting
	} else if len(t.players) <= 0 {
		t.BroadcastMessage("§eNo winner!", TypeMessage)

		t.task.timer = RestartTimer
		t.state = StateRestarting
	}
}

func (t *TNTRun) AddSpectator(p *player.Player) {
	t.spectators = append(t.spectators, p.H())

	p.SetGameMode(world.GameModeSpectator)
	p.Teleport(t.config.Lobby)

	t.RemovePlayerHandle(p.H())
	t.BroadcastMessage(fmt.Sprintf("§e%s §7has been eliminated", p.Name()), TypeMessage)

	t.CheckWinner()
}

// IsPlayer ...
func (t *TNTRun) IsPlayer(player *player.Player) bool {
	playerUUID := player.UUID()
	for _, handle := range t.players {
		if handle.UUID() == playerUUID {
			return true
		}
	}

	return false
}

// IsPlayerByHandle ...
func (t *TNTRun) IsPlayerByHandle(playerHandle *world.EntityHandle) bool {
	targetUUID := playerHandle.UUID()
	for _, handle := range t.players {
		if handle.UUID() == targetUUID {
			return true
		}
	}

	return false
}

// RemovePlayer ...
func (t *TNTRun) RemovePlayer(player *player.Player) {
	playerUUID := player.UUID()
	for i, handle := range t.players {
		if handle.UUID() == playerUUID {
			t.players = utils.RemoveIndex(t.players, i)
		}
	}
}

// RemovePlayerHandle ...
func (t *TNTRun) RemovePlayerHandle(playerHandle *world.EntityHandle) {
	targetUUID := playerHandle.UUID()
	for i, handle := range t.players {
		if handle.UUID() == targetUUID {
			t.players = utils.RemoveIndex(t.players, i)
		}
	}
}
