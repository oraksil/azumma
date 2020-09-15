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
	PeerName  string
	Players   []*Player
	CreatedAt time.Time
}

type SignalingInfo struct {
	Id     int64
	Game   *RunningGame
	Data   string
	IsLast bool
}

type GameRepository interface {
	GetById(id int) (*Game, error)
	Find(offset, limit int) []*Game
}

type RunningGameRepository interface {
	Find(offset, limit int) []*RunningGame
	FindById(id int64) (*RunningGame, error)
	Save(game *RunningGame) (*RunningGame, error)
}

type SignalingRepository interface {
	Save(signalingInfo *SignalingInfo) (*SignalingInfo, error)
	FindByRunningGameId(runningGameId int64, sinceId int64) (*SignalingInfo, error)
}
