package models

type Session struct {
	Player *Player
}

type PrepareOrakki struct {
	GameId int64
}

type ProvisionInfo struct {
	OrakkiId string
	Pack     *Pack
}

type GameInfo struct {
	GameId     int64
	MaxPlayers int
}

type PlayerParticipation struct {
	GameId   int64
	PlayerId int64
}

type Orakki struct {
	Id    string
	State int
}

type PeerInfo struct {
	Token    string
	GameId   int64
	PlayerId int64
}

type SdpInfo struct {
	Peer             PeerInfo
	SdpBase64Encoded string
}

type IceCandidate struct {
	Peer             PeerInfo
	IceBase64Encoded string
	Seq              int64
}

type TurnAuth struct {
	Username string
	Password string
}
