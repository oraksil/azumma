package dto

import (
	"testing"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestMapGameEntityToDto(t *testing.T) {
	e := models.Pack{Id: 1, Status: models.PackStatusReady, Title: "Game", Description: "Desc", MaxPlayers: 3}

	dto := PackToDto(&e)

	assert.Equal(t, e.Id, dto.Id)
	assert.Equal(t, e.Title, dto.Title)
	assert.Equal(t, e.Description, dto.Desc)
	assert.Equal(t, e.MaxPlayers, dto.MaxPlayers)
	assert.Equal(t, dto.Status, "ready")
}

func TestMapNestedFields(t *testing.T) {
	ices := make([]*models.IceCandidate, 0)
	ices = append(ices, &models.IceCandidate{
		Peer: models.PeerInfo{
			Token:    "abcd",
			GameId:   1234,
			PlayerId: 789,
		},
		IceBase64Encoded: "xyz",
	})

	iceDto := IcesToDto(ices)

	assert.Equal(t, iceDto[0].Token, ices[0].Peer.Token)
	assert.Equal(t, iceDto[0].Base64Encoded, ices[0].IceBase64Encoded)
}
