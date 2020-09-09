package dto

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestMapGameEntityToDto(t *testing.T) {
	e := models.Game{Id: 1, Title: "Game", Description: "Desc", MaxPlayers: 3}

	var dto AvailableGameDto
	mapstructure.Decode(e, &dto)

	assert.Equal(t, e.Id, dto.Id)
	assert.Equal(t, e.Title, dto.Title)
	assert.Equal(t, e.Description, dto.Desc)
	assert.Equal(t, e.MaxPlayers, dto.MaxPlayers)
}
