package models

import "time"

type Game struct {
	Id          int
	Title       string
	Maker       string
	Description string
	MaxPlayers  int
}

type Player struct {
	Id         int64
	Name       string
	TotalCoins int
}

type Orakki struct {
	Id       string
	State    int
	PeerName string
}

type RunningGame struct {
	Id        int64
	Orakki    *Orakki
	Game      *Game
	PeerName  string
	Players   []*Player
	CreatedAt time.Time
}

type SignalingInfo struct {
	Id       int64
	OrakkiId string
	Game     *RunningGame
	Data     string
	IsLast   bool
}

type GameRepository interface {
	GetGameById(id int) (*Game, error)
	FindAvailableGames(offset, limit int) []*Game
	FindRunningGames(offset, limit int) []*RunningGame
	FindRunningGameByOrakkiId(orakkiId string) (*RunningGame, error)
	SaveRunningGame(game *RunningGame) (*RunningGame, error)
}

type SignalingRepository interface {
	SaveSignalingInfo(signalingInfo *SignalingInfo) (*SignalingInfo, error)
	UpdateSignalingInfo(signalingInfo *SignalingInfo) (*SignalingInfo, error)
	FindSignalingInfo(orakkiId string, order string, num int) (*SignalingInfo, error)
	FindIceCandidate(orakkiId string, seqAfter int) (*SignalingInfo, error)
}
