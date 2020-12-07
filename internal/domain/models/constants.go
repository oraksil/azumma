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
	MAX_COINS              = 5
	INITIAL_COINS          = MAX_COINS
	TIME_TO_A_COIN_IN_SECS = 10 * 60
)
