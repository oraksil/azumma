package dto

type PlayerDto struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Hash       string `json:"hash"`
	TotalCoins int    `json:"total_coins"`
}

type PackDto struct {
	Id         int    `json:"id"`
	Status     string `json:"status"`
	Title      string `json:"title"`
	Maker      string `json:"maker"`
	Desc       string `mapstructure:"Description" json:"description"`
	MaxPlayers int    `json:"max_players"`
	PosterUrl  string `json:"poster_url"`
	RomName    string `json:"rom_name"`
}

type GameDto struct {
	Id int64 `json:"id"`
	// Orakki    *Orakki
	// Game      *Game
}

type SdpDto struct {
	Token         string `json:"token"`
	Base64Encoded string `json:"encoded"`
}

type IceCandidateDto struct {
	Token         string `json:"token"`
	Base64Encoded string `json:"encoded"`
	Seq           int64  `json:"seq"`
}

type JoinableDto struct {
	Token string `json:"token"`
}
