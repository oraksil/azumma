package models

import (
	"time"
)

type Pack struct {
	Id          int
	Status      int
	Title       string
	Maker       string
	Description string
	MaxPlayers  int
	PosterUrl   string
	RomName     string
}

func (p *Pack) GetStatusAsString() string {
	switch p.Status {
	case PackStatusReady:
		return "ready"
	case PackStatusPreparing:
		return "prepare"
	default:
		return "invalid"
	}
}

type Player struct {
	Id              int64
	Name            string
	LastCoins       int
	LastCoinsUsedAt time.Time
}

func (p *Player) Hash() string {
	return ""
}

func (p *Player) UpdateCoins() {
	elapsedSecs := time.Since(p.LastCoinsUsedAt)
	chargedCoins := int(elapsedSecs / TIME_TO_A_COIN_IN_SECS)

	totalCoins := p.LastCoins + chargedCoins
	if p.LastCoins = totalCoins; totalCoins > MAX_COINS {
		p.LastCoins = MAX_COINS
	}
	p.LastCoinsUsedAt = time.Now().UTC()
}

func (p *Player) UseCoins(coins int) {
	totalCoins := p.LastCoins - coins
	if p.LastCoins = totalCoins; totalCoins < 0 {
		p.LastCoins = 0
	}
	p.LastCoinsUsedAt = time.Now().UTC()
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

type UserFeedback struct {
	Id        int64
	Content   string
	CreatedAt time.Time
}

type PlayerRepository interface {
	GetById(id int64) (*Player, error)
	Save(player *Player) (*Player, error)
}

type PackRepository interface {
	GetById(id int) (*Pack, error)
	FindAll(offset, limit int) []*Pack
	FindByStatus(status, offset, limit int) []*Pack
}

type GameRepository interface {
	GetById(id int64) (*Game, error)
	Save(game *Game) (*Game, error)
}

type SignalingRepository interface {
	FindByToken(token string, sinceId int64) ([]*Signaling, error)
	Save(signaling *Signaling) (*Signaling, error)
}

type UserFeedbackRepository interface {
	Save(feedback *UserFeedback) (*UserFeedback, error)
}
