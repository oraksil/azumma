package services

type MessageService interface {
	Identifier() string
	Send(to, msgType string, payload interface{}) error
	SendToAny(msgType string, payload interface{}) error
	Broadcast(msgType string, payload interface{}) error
	Request(to, msgType string, payload interface{}) (interface{}, error)
}

type OrakkiDriver interface {
	RunInstance(peerName string) (string, error)
	DeleteInstance(id string) error
}
