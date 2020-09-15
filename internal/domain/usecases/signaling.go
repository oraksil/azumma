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
	RunningGameRepo models.RunningGameRepository
	SignalingRepo   models.SignalingRepository
	MessageService  services.MessageService
}

func (uc *SignalingUseCase) NewOffer(runningGameId int64, playerId int64, sdpString string) (*models.SignalingInfo, error) {
	offer := webrtc.SessionDescription{}

	err := json.Unmarshal([]byte(sdpString), &offer)

	if err != nil {
		return nil, err
	}

	game, err := uc.RunningGameRepo.FindById(runningGameId)

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

	SignalingInfo := models.SignalingInfo{
		Game: game,
		// PlayerId: playerId,
		// State:    models.POLLING_STATE_DATA_FETCHED,
		Data: setupAnswer.Answer,
	}

	return &SignalingInfo, err
}

func (uc *SignalingUseCase) GetIceCandidate(runningGameId int64, sinceId int64) (*models.SignalingInfo, error) {
	signalingInfo, err := uc.SignalingRepo.FindByRunningGameId(runningGameId, sinceId)

	if err != nil {
		return nil, err
	}

	return signalingInfo, nil
}

func (uc *SignalingUseCase) AddServerIceCandidate(runningGameId int64, iceCandidate string) (*models.SignalingInfo, error) {
	game, err := uc.RunningGameRepo.FindById(runningGameId)

	if game == nil {
		return nil, errors.New("No game matched to given ID")
	}

	SignalingInfo := models.SignalingInfo{
		Game: game,
		Data: iceCandidate,
	}

	var saved *models.SignalingInfo
	if iceCandidate == "" {
		// get lastly added signaling info to set is_last to 1
		lastSignalingInfo, _ := uc.SignalingRepo.FindByRunningGameId(runningGameId, 1)
		lastSignalingInfo.IsLast = true

		saved, err = uc.SignalingRepo.Save(lastSignalingInfo)
	} else {
		saved, err = uc.SignalingRepo.Save(&SignalingInfo)
	}

	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (uc *SignalingUseCase) AddIceCandidate(runningGameId int64, playerId int64, iceCandidate string) (*models.SignalingInfo, error) {
	game, _ := uc.RunningGameRepo.FindById(runningGameId)

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

	SignalingInfo := models.SignalingInfo{
		Game: game,
	}

	return &SignalingInfo, nil
}
