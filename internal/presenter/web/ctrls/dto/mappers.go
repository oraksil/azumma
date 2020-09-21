package dto

import (
	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
)

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
	var sdpDto SdpDto
	mapstructure.Decode(src, &sdpDto)

	return &sdpDto
}

func IceToDto(src *models.IceCandidate) *IceCandidateDto {
	var iceDto IceCandidateDto
	mapstructure.Decode(src, &iceDto)

	return &iceDto
}
