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
	NeededPlayers = 3
	StartTimer    = 30
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
	players    []*player.Player
	spectators []*player.Player
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
		players:    []*player.Player{},
		spectators: []*player.Player{},
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
		t.players = append(t.players, p)

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
	var players []*player.Player
	players = append(players, t.players...)
	players = append(players, t.spectators...)

	for _, p := range players {
		if msgType == TypeMessage {
			p.Message(msg)
		} else if msgType == TypePopup {
			p.SendPopup(msg)
		} else if msgType == TypeTitle {
			p.SendTitle(title.New(msg))
		}
	}
}

// CheckWinner ...
func (t *TNTRun) CheckWinner() {
	if len(t.players) == 1 {
		winner := t.players[0]

		t.BroadcastMessage(fmt.Sprintf("§e%s §7won the game!", winner.Name()), TypeMessage)

		t.task.timer = RestartTimer
		t.state = StateRestarting
	} else if len(t.players) <= 0 {
		t.BroadcastMessage("§eNo winner!", TypeMessage)

		t.task.timer = RestartTimer
		t.state = StateRestarting
	}
}

// AddSpectator ...
func (t *TNTRun) AddSpectator(player *player.Player) {
	t.BroadcastMessage(fmt.Sprintf("§e%s §7has been eliminated", player.Name()), TypeMessage)
	t.RemovePlayer(player)

	t.spectators = append(t.spectators, player)
	player.SetGameMode(world.GameModeSpectator)
	player.Teleport(t.config.Lobby)

	t.CheckWinner()
}

// IsPlayer ...
func (t *TNTRun) IsPlayer(player *player.Player) bool {
	for _, p := range t.players {
		if player.Name() == p.Name() {
			return true
		}
	}

	return false
}

// RemovePlayer ...
func (t *TNTRun) RemovePlayer(player *player.Player) {
	for i, p := range t.players {
		if player.Name() == p.Name() {
			t.players = utils.RemoveIndex(t.players, i)
		}
	}
}
