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
	OrakkiId        string    `db:"orakki_id"`
	OrakkiState     int       `db:"orakki_state"`
	PeerName        string    `db:"peer_name"`
	PackId          int       `db:"pack_id"`
	FirstPlayerId   int64     `db:"first_player_id"`
	JoinedPlayerIds string    `db:"joined_player_ids"`
	CreatedAt       time.Time `db:"created_at"`
}

// ConnectionDescriptionData : record info of connection for each player, such as signaling state
type SignalingData struct {
	Id        int64     `db:"id"`
	OrakkiId  string    `db:"orakki_id"`
	Data      string    `db:"data"`
	IsLast    bool      `db:"is_last"`
	CreatedAt time.Time `db:"created_at"`
}
