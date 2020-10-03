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

func (p *Player) Hash() string {
	return ""
}

type Game struct {
	Id        int64
	Orakki    *Orakki
	Pack      *Pack
	Players   []*Player
	CreatedAt time.Time
}

func (g *Game) Join(p *Player) {
	for _, pl := range g.Players {
		if pl.Id == p.Id {
			return
		}
	}
	g.Players = append(g.Players, p)
}

func (g *Game) Leave(p *Player) {
	for i, pl := range g.Players {
		if pl.Id == p.Id {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			break
		}
	}
}

type Signaling struct {
	Id       int64
	Token    string
	GameId   int64
	PlayerId int64
	Data     string
}

type PlayerRepository interface {
	GetById(id int64) (*Player, error)
	Save(player *Player) (*Player, error)
}

type PackRepository interface {
	GetById(id int) (*Pack, error)
	Find(offset, limit int) []*Pack
}

type GameRepository interface {
	GetById(id int64) (*Game, error)
	Find(offset, limit int) []*Game
	Save(game *Game) (*Game, error)
}

type SignalingRepository interface {
	Find(token string, sinceId int64) ([]*Signaling, error)
	Save(signaling *Signaling) (*Signaling, error)
}
