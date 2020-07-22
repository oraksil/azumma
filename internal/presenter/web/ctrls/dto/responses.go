package dto

type AvailableGame struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Desc       string `mapstructure:"Description" json:"description"`
	Maker      string `json:"maker"`
	MaxPlayers int    `json:"max_players"`
}
