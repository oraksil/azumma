package services

import "time"

type ServiceConfig struct {
	MqRpcUri        string
	MqRpcNamespace  string
	MqRpcIdentifier string

	StaticOrakkiId       string
	OrakkiContainerImage string
	GipanContainerImage  string
	ProvisionMaxWait     time.Duration
}
