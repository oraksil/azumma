package usecases

import (
	"errors"

	"gitlab.com/oraksil/azumma/internal/domain/models"

	"encoding/json"

	"github.com/pion/webrtc/v2"
)

type SignalingUseCase struct {
	GameRepository models.GameRepository
}

// NewUserSdp : Accept user's sdp
func (uc *SignalingUseCase) NewOffer(gameId int64, playerId int64, sdpString string) (*models.ConnectionInfo, error) {
	// validation needed????
	offer := webrtc.SessionDescription{}

	err := json.Unmarshal([]byte(sdpString), &offer)

	if err != nil {
		return nil, err
	}

	// Player, _ := uc.GameRepository.GetPlayerById(playerId)
	game, err := uc.GameRepository.FindRunningGameById(gameId)

	if game == nil {
		return nil, errors.New("No game exists with given ID")
	}

	connectionInfo := models.ConnectionInfo{
		Game:     game,
		PlayerId: playerId,
		State:    models.CONNECTION_STATE_OFFER_REQUESTED,
	}

	saved, err := uc.GameRepository.SaveConnectionInfo(&connectionInfo)
	if err != nil {
		return nil, err
	}

	return saved, err
}

func (uc *SignalingUseCase) TryFetchAnswer(gameId int64, playerId int64) (bool, string, error) {
	game, err := uc.GameRepository.FindRunningGameById(gameId)

	if err != nil {
		return false, "", errors.New("No game exists with given ID")
	}

	connectionInfo, _ := uc.GameRepository.GetConnectionInfo(game.Orakki.Id, playerId)

	if connectionInfo.State == models.CONNECTION_STATE_ANSWER_SET {
		if connectionInfo.ServerData != "" {
			// change state ??
			// connectionInfo.State = models.CONNECTION_STATE_ICE_EXCHANGING
			return true, connectionInfo.ServerData, nil
		} else {
			return true, "", errors.New("Empty Answer is set")
		}
	} else {
		return false, "", nil
	}
}
