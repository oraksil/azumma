package models

type Game struct {
	Id          int64  `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	MaxPlayers  int    `db:"max_players" json:"max_players"`
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
	FindAvailableGames(offset, limit int) []*Game
	FindRunningGames(offset, limit int) []*RunningGame
}
