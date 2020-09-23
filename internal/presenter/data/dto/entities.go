package dto

import "time"

type PackData struct {
	Id          int    `db:"id"`
	Title       string `db:"title"`
	Maker       string `db:"maker"`
	Description string `db:"description"`
	MaxPlayers  int    `db:"max_players"`
}
type GameData struct {
	Id              int64     `db:"id"`
	PackId          int       `db:"pack_id"`
	OrakkiId        string    `db:"orakki_id"`
	OrakkiState     int       `db:"orakki_state"`
	FirstPlayerId   int64     `db:"first_player_id"`
	JoinedPlayerIds string    `db:"joined_player_ids"`
	CreatedAt       time.Time `db:"created_at"`
}

// ConnectionDescriptionData : record info of connection for each player, such as signaling state
type SignalingData struct {
	Id        int64     `db:"id"`
	GameId    int64     `db:"game_id"`
	Data      string    `db:"data"`
	CreatedAt time.Time `db:"created_at"`
}
