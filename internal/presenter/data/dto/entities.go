package dto

import (
	"strconv"
	"strings"
	"time"

	"github.com/oraksil/azumma/internal/domain/models"
)

type PlayerData struct {
	Id                  int64     `db:"id"`
	Name                string    `db:"name"`
	TotalCoinsUsed      int       `db:"total_coins_used"`
	CoinsUsedInCharging int       `db:"coins_used_in_charging"`
	ChargingStartedAt   time.Time `db:"charging_started_at"`
}

type PackData struct {
	Id          int    `db:"id"`
	Status      int    `db:"status"`
	Title       string `db:"title"`
	Maker       string `db:"maker"`
	Description string `db:"description"`
	MaxPlayers  int    `db:"max_players"`
	PosterUrl   string `db:"poster_url"`
	RomName     string `db:"rom_name"`
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

func (d *GameData) SetJoinedPlayers(players []*models.Player) {
	playerIds := make([]string, len(players))
	for i, p := range players {
		playerIds[i] = strconv.FormatInt(p.Id, 10)
	}
	d.JoinedPlayerIds = strings.Join(playerIds, ",")
}

func (d *GameData) GetJoinedPlayers() []*models.Player {
	playerIds := strings.FieldsFunc(d.JoinedPlayerIds, func(c rune) bool { return c == ',' })
	players := make([]*models.Player, 0)
	for _, pIdString := range playerIds {
		pId, _ := strconv.ParseInt(pIdString, 10, 0)
		players = append(players, &models.Player{Id: pId})
	}

	return players
}

// ConnectionDescriptionData : record info of connection for each player, such as signaling state
type SignalingData struct {
	Id        int64     `db:"id"`
	Token     string    `db:"token"`
	GameId    int64     `db:"game_id"`
	PlayerId  int64     `db:"player_id"`
	Data      string    `db:"data"`
	CreatedAt time.Time `db:"created_at"`
}

type UserFeedbackData struct {
	Id        int64     `db:"id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}
