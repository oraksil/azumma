package services

import "time"

type ServiceConfig struct {
	// for azumma
	DbUri string

	MqRpcUri        string
	MqRpcNamespace  string
	MqRpcIdentifier string

	StaticOrakkiId   string
	ProvisionMaxWait time.Duration

	// for orakki
	OrakkiMqRpcUri       string
	OrakkiMqRpcNamespace string

	OrakkiContainerImage string
	GipanContainerImage  string

	TurnServerUri      string
	TurnServerUsername string
	TurnServerPassword string

	GipanResolution       string
	GipanFps              string
	GipanKeyframeInterval string

	OrakkiDriverK8SConfigPath        string
	OrakkiDriverK8SNamespace         string
	OrakkiDriverK8SNodeSelectorKey   string
	OrakkiDriverK8SNodeSelectorValue string
}
