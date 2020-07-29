package dto

import (
	"github.com/mitchellh/mapstructure"
	"gitlab.com/oraksil/azumma/internal/domain/models"
)

func GamesToDto(src []*models.Game) []*AvailableGameDto {
	var gamesDto []*AvailableGameDto
	mapstructure.Decode(src, &gamesDto)

	return gamesDto
}

func RunningGameToDto(src *models.RunningGame) *RunningGameDto {
	gameDto := RunningGameDto{
		Id:        src.Id,
		CreatedAt: src.CreatedAt.Unix(),
	}

	return &gameDto
}
