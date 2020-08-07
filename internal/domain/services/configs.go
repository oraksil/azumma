package services

import "time"

type ServiceConfig struct {
	UseStaticOrakki      bool
	StaticOrakkiId       string
	StaticOrakkiPeerName string
	PeerName             string
	ProvisionMaxWait     time.Duration
}
