package models

type Game struct {
	Id          int64
	Title       string
	Description string
	MaxPlayers  int
}

type Player struct {
	Id         int64
	Name       string
	TotalCoins int
}

type RunningGame struct {
	Id        int64
	Game      *Game
	Players   []*Player
	CreatedAt int64
}

type GameRepository interface {
	GetAllAvailableGames(offset, limit int) []*Game
	GetAllRunningGames(offset, limit int) []*RunningGame
}
