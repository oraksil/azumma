package dto

type GameData struct {
	Id          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	MaxPlayers  int    `db:"max_players"`
}
