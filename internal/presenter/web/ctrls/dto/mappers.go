package dto

import (
	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
)

func PlayerToDto(src *models.Player) *PlayerDto {
	var playerDto PlayerDto
	mapstructure.Decode(src, &playerDto)

	return &playerDto
}

func PackToDto(src []*models.Pack) []*PackDto {
	var packsDto []*PackDto
	mapstructure.Decode(src, &packsDto)

	return packsDto
}

func GameToDto(src *models.Game) *GameDto {
	gameDto := GameDto{
		Id:        src.Id,
		CreatedAt: src.CreatedAt.Unix(),
	}

	return &gameDto
}

func SdpToDto(src *models.SdpInfo) *SdpDto {
	return &SdpDto{
		Token:         src.Peer.Token,
		Base64Encoded: src.SdpBase64Encoded,
	}
}

func IcesToDto(src []*models.IceCandidate) []*IceCandidateDto {
	icesDto := make([]*IceCandidateDto, 0)

	for _, ice := range src {
		icesDto = append(icesDto, &IceCandidateDto{
			Token:         ice.Peer.Token,
			Base64Encoded: ice.IceBase64Encoded,
		})
	}

	return icesDto
}
