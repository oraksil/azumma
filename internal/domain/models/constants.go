package models

const (
	PackStatusReady = iota + 1
	PackStatusPreparing
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
	OrakkiStateInit = iota + 1
	OrakkiStateReady
	OrakkiStatePause
	OrakkiStatePlay
	OrakkiStateExit
	OrakkiStatePanic
)
