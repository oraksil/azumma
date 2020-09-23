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
	PeerId int64  `json:"peer_id"`
	Sdp    string `mapstructure:"SdpBase64Encoded" json:"sdp"`
}

type IceCandidateDto struct {
	PeerId int64  `json:"peer_id"`
	Ice    string `mapstructure:"IceBase64Encoded" json:"ice"`
	IsLast bool   `json:"is_last"`
	Seq    int64  `json:"seq"`
}
