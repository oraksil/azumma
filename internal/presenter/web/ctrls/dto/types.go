package dto

type PackDto struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Maker      string `json:"maker"`
	Desc       string `mapstructure:"Description" json:"description"`
	MaxPlayers int    `json:"max_players"`
}

type GameDto struct {
	Id int64 `json:"id"`
	// Orakki    *Orakki
	// Game      *Game
	CreatedAt int64 `json:"created_at"`
}

type SdpDto struct {
}

type IceCandidateDto struct {
}