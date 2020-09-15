package usecases

import (
	"errors"

	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"
	"github.com/oraksil/azumma/pkg/utils"

	"github.com/pion/webrtc/v2"

	"time"
)

type SignalingUseCase struct {
	GameRepo       models.GameRepository
	SignalingRepo  models.SignalingRepository
	MessageService services.MessageService
}

func (uc *SignalingUseCase) NewOffer(gameId int64, playerId int64, sdpString string) (*models.Signaling, error) {
	offer := webrtc.SessionDescription{}

	err := json.Unmarshal([]byte(sdpString), &offer)

	if err != nil {
		return nil, err
	}

	game, err := uc.GameRepo.FindById(gameId)

	if game == nil {
		return nil, errors.New("No game exists with given ID")
	}

	b64EncodedOffer, err := utils.EncodeToB64EncodedJsonStr(offer)
	if err != nil {
		return nil, err
	}

	// sdp response from orakki
	resp, err := uc.MessageService.Request(
		game.PeerName,
		models.MSG_SETUP_WITH_NEW_OFFER,
		b64EncodedOffer,
		5*time.Second,
	)

	var setupAnswer models.SetupAnswer
	mapstructure.Decode(resp, &setupAnswer)

	Signaling := models.Signaling{
		Game: game,
		// PlayerId: playerId,
		// State:    models.POLLING_STATE_DATA_FETCHED,
		Data: setupAnswer.Answer,
	}

	return &Signaling, err
}

func (uc *SignalingUseCase) GetIceCandidate(gameId int64, sinceId int64) (*models.Signaling, error) {
	signaling, err := uc.SignalingRepo.FindByGameId(gameId, sinceId)

	if err != nil {
		return nil, err
	}

	return signaling, nil
}

func (uc *SignalingUseCase) AddServerIceCandidate(gameId int64, iceCandidate string) (*models.Signaling, error) {
	game, err := uc.GameRepo.FindById(gameId)

	if game == nil {
		return nil, errors.New("No game matched to given ID")
	}

	signaling := models.Signaling{
		Game: game,
		Data: iceCandidate,
	}

	var saved *models.Signaling
	if iceCandidate == "" {
		// get lastly added signaling info to set is_last to 1
		lastSignaling, _ := uc.SignalingRepo.FindByGameId(gameId, 1)
		lastSignaling.IsLast = true

		saved, err = uc.SignalingRepo.Save(lastSignaling)
	} else {
		saved, err = uc.SignalingRepo.Save(&signaling)
	}

	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (uc *SignalingUseCase) AddIceCandidate(gameId int64, playerId int64, iceCandidate string) (*models.Signaling, error) {
	game, _ := uc.GameRepo.FindById(gameId)

	if game == nil {
		return nil, errors.New("No game exists with given ID")
	}

	resp, err := uc.MessageService.Request(
		game.PeerName,
		models.MSG_REMOTE_ICE_CANDIDATE,
		models.Icecandidate{
			PlayerId:  playerId,
			IceString: iceCandidate,
		},
		5*time.Second,
	)

	if err != nil {
		return nil, err
	}

	var setupAnswer models.SetupAnswer
	mapstructure.Decode(resp, &setupAnswer)

	Signaling := models.Signaling{
		Game: game,
	}

	return &Signaling, nil
}
