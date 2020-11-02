package services

import "time"

type ServiceConfig struct {
	DbUri string

	MqRpcUri        string
	MqRpcNamespace  string
	MqRpcIdentifier string

	StaticOrakkiId       string
	OrakkiContainerImage string
	GipanContainerImage  string
	ProvisionMaxWait     time.Duration

	OrakkiDriverK8SConfigPath        string
	OrakkiDriverK8SNamespace         string
	OrakkiDriverK8SNodeSelectorKey   string
	OrakkiDriverK8SNodeSelectorValue string

	TurnServerUri      string
	TurnServerUsername string
	TurnServerPassword string
}
