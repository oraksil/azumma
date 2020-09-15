package models

import "time"

type Pack struct {
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

type Game struct {
	Id        int64
	Orakki    *Orakki
	Pack      *Pack
	PeerName  string
	Players   []*Player
	CreatedAt time.Time
}

type Signaling struct {
	Id     int64
	Game   *Game
	Data   string
	IsLast bool
}

type PackRepository interface {
	GetById(id int) (*Pack, error)
	Find(offset, limit int) []*Pack
}

type GameRepository interface {
	Find(offset, limit int) []*Game
	FindById(id int64) (*Game, error)
	Save(game *Game) (*Game, error)
}

type SignalingRepository interface {
	Save(signaling *Signaling) (*Signaling, error)
	FindByGameId(gameId int64, sinceId int64) (*Signaling, error)
}
