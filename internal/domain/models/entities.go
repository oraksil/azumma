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

// ConnectionState : record info of connection for each player
type ConnectionInfo struct {
	Id       int64
	OrakkiId string
	PlayerId int64
	State    int
}
type GameRepository interface {
	GetGameById(id int) (*Game, error)

	FindAvailableGames(offset, limit int) []*Game
	FindRunningGames(offset, limit int) []*RunningGame

	SaveRunningGame(game *RunningGame) (*RunningGame, error)

	GetPlayerById(id int64) (*Player, error)
	SaveConnectionInfo(connectionState *ConnectionInfo) (*ConnectionInfo, error)

	GetOrakkiById(id string) (*Orakki, error)
}
