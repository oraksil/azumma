package models

const (
	PackStatusPreparing = iota
	PackStatusReady
)

const (
	MsgPrepareOrakki      = "MsgPrepareOrakki"
	MsgSetupWithNewOffer  = "MsgSetupWithNewOffer"
	MsgRemoteIceCandidate = "MsgRemoteIceCandidate"
	MsgStartGame          = "MsgStartGame"
	MsgPlayerJoined       = "MsgPlayerJoined"
	MsgPlayerJoinFailed   = "MsgPlayerJoinFailed"
	MsgPlayerLeft         = "MsgPlayerLeft"
)

const (
	OrakkiStateInit = iota
	OrakkiStateReady
	OrakkiStatePause
	OrakkiStatePlay
	OrakkiStateExit
	OrakkiStatePanic
)

const (
	INITIAL_COINS          = 10
	MAX_COINS              = 10
	TIME_TO_A_COIN_IN_SECS = 10 * 60
)
