package dto

type PlayerDto struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Hash            string `json:"hash"`
	LastCoins       int    `json:"lastCoins"`
	LastCoinsUsedAt int64  `json:"lastCoinsUsedAt"`
}

type CoinDto struct {
	LastCoins       int   `json:"lastCoins"`
	LastCoinsUsedAt int64 `json:"lastCoinsUsedAt"`
}

type PackDto struct {
	Id         int    `json:"id"`
	Status     string `json:"status"`
	Title      string `json:"title"`
	Maker      string `json:"maker"`
	Desc       string `mapstructure:"Description" json:"description"`
	MaxPlayers int    `json:"maxPlayers"`
	PosterUrl  string `json:"posterUrl"`
	RomName    string `json:"romName"`
}

type GameDto struct {
	Id int64 `json:"id"`
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
	Token        string `json:"token"`
	TurnUsername string `json:"username"`
	TurnPassword string `json:"password"`
}
