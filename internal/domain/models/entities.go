package models

import "time"

type Game struct {
	Id          int
	Title       string
	Maker       string
	Description string
	MaxPlayers  int
}

type Player struct {
	Id         int64
	Name       string
	TotalCoins int
}

const (
	ORAKKI_STATE_INIT = iota
	ORAKKI_STATE_READY
	ORAKKI_STATE_SIGNALING
	ORAKKI_STATE_SIGNALING_SDP
	ORAKKI_STATE_SIGNALING_ICE
	ORAKKI_STATE_PLAYING
)

type Orakki struct {
	Id       string
	State    int
	PeerName string
}

type RunningGame struct {
	Id        int64
	Orakki    *Orakki
	Game      *Game
	Players   []*Player
	CreatedAt time.Time
}

type GameRepository interface {
	GetGameById(id int) (*Game, error)

	FindAvailableGames(offset, limit int) []*Game
	FindRunningGames(offset, limit int) []*RunningGame

	SaveRunningGame(game *RunningGame) (*RunningGame, error)
}
