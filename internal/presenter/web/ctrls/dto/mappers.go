package dto

import (
	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
)

func PlayerToDto(src *models.Player) *PlayerDto {
	var playerDto PlayerDto
	mapstructure.Decode(src, &playerDto)

	playerDto.LastCoinsUsedAt = src.LastCoinsUsedAt.Unix()

	return &playerDto
}

func PlayerToCoinDto(src *models.Player) *CoinDto {
	return &CoinDto{
		LastCoins:       src.LastCoins,
		LastCoinsUsedAt: src.LastCoinsUsedAt.Unix(),
	}
}

func PackToDto(src *models.Pack) *PackDto {
	var packDto PackDto
	mapstructure.Decode(src, &packDto)

	packDto.Status = src.GetStatusAsString()

	return &packDto
}

func PacksToDto(src []*models.Pack) []*PackDto {
	packsDto := make([]*PackDto, 0)

	for _, pack := range src {
		packDto := PackToDto(pack)
		packsDto = append(packsDto, packDto)
	}

	return packsDto
}

func GameToDto(src *models.Game) *GameDto {
	return &GameDto{Id: src.Id}
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
