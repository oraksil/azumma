package usecases

import (
	"gitlab.com/oraksil/azumma/internal/domain/models"

	"encoding/json"

	"github.com/pion/webrtc/v2"
)

type SignalingUseCase struct {
	GameRepository models.GameRepository
}

// NewUserSdp : Accept user's sdp
func (uc *SignalingUseCase) NewOffer(orakkiId string, playerId int64, sdpString string) (*models.ConnectionInfo, error) {
	// validation needed????
	offer := webrtc.SessionDescription{}

	// b, err := base64.StdEncoding.DecodeString(sdpString)

	// if err != nil {
	// return nil, err
	// }

	err := json.Unmarshal([]byte(sdpString), &offer)

	if err != nil {
		return nil, err
	}

	// Player, _ := uc.GameRepository.GetPlayerById(playerId)

	connectionInfo := models.ConnectionInfo{
		OrakkiId: orakkiId,
		PlayerId: playerId,
		State:    models.ConnectionStateInit,
	}

	saved, err := uc.GameRepository.SaveConnectionInfo(&connectionInfo)
	if err != nil {
		return nil, err
	}

	return saved, err
}
