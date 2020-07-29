package dto

import "time"

type GameData struct {
	Id          int    `db:"id"`
	Title       string `db:"title"`
	Maker       string `db:"maker"`
	Description string `db:"description"`
	MaxPlayers  int    `db:"max_players"`
}

type RunningGameData struct {
	Id              int64     `db:"id"`
	OrakkiId        string    `db:"orakki_id"`
	OrakkiState     int       `db:"orakki_state"`
	PeerName        string    `db:"peer_name"`
	GameId          int       `db:"game_id"`
	FirstPlayerId   int64     `db:"first_player_id"`
	JoinedPlayerIds string    `db:"joined_player_ids"`
	CreatedAt       time.Time `db:"created_at"`
}
