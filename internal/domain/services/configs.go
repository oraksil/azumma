package services

import "time"

type ServiceConfig struct {
	MqRpcUri       string
	MqRpcNamespace string

	DbUri string

	PeerName string

	UseStaticOrakki      bool
	StaticOrakkiId       string
	StaticOrakkiPeerName string

	OrakkiContainerImage string
	GipanContainerImage  string
	ProvisionMaxWait     time.Duration
}
