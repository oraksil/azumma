package models

type OrakkiState struct {
	OrakkiId string
	State    int
}

type SetupAnswer struct {
	PeerId string
	Answer string
}

type Icecandidate struct {
	PlayerId  int64
	OrakkiId  string
	IceString string
}
