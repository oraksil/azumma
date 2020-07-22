package dto

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"gitlab.com/oraksil/sil/backend/internal/domain/models"
)

func TestMapGameEntityToDto(t *testing.T) {
	e := models.Game{Id: 1, Title: "Game", Description: "Desc", MaxPlayers: 3}

	var dto AvailableGame
	mapstructure.Decode(e, &dto)

	assert.Equal(t, e.Id, dto.Id)
	assert.Equal(t, e.Title, dto.Title)
	assert.Equal(t, e.Description, dto.Desc)
	assert.Equal(t, e.MaxPlayers, dto.MaxPlayers)
}
