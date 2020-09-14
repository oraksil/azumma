package usecases

import (
	"errors"

	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"

	"github.com/pion/webrtc/v2"

	"time"
)

type SignalingUseCase struct {
	GameRepository      models.GameRepository
	SignalingRepository models.SignalingRepository
	MessageService      services.MessageService
}

func (uc *SignalingUseCase) NewOffer(orakkiId string, playerId int64, sdpString string) (*models.SignalingInfo, error) {
	offer := webrtc.SessionDescription{}

	err := json.Unmarshal([]byte(sdpString), &offer)

	if err != nil {
		return nil, err
	}

	game, err := uc.GameRepository.FindRunningGameByOrakkiId(orakkiId)

	if game == nil {
		return nil, errors.New("No game exists with given ID")
	}

	// sdp response from orakki
	resp, err := uc.MessageService.Request(
		game.PeerName,
		models.MSG_SETUP_WITH_NEW_OFFER,
		offer,
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

func (uc *SignalingUseCase) GetIceCandidate(orakkiId string, seqAfter int, num int) (*models.SignalingInfo, error) {
	signalingInfo, err := uc.SignalingRepository.FindIceCandidate(orakkiId, seqAfter, num)

	if err != nil {
		return nil, err
	}

	return signalingInfo, nil
}

func (uc *SignalingUseCase) AddServerIceCandidate(orakkiId string, iceCandidate string) (*models.SignalingInfo, error) {
	game, err := uc.GameRepository.FindRunningGameByOrakkiId(orakkiId)

	if game == nil {
		return nil, errors.New("No game matched to given ID")
	}

	SignalingInfo := models.SignalingInfo{
		Game:     game,
		Data:     iceCandidate,
		OrakkiId: orakkiId,
	}

	var saved *models.SignalingInfo
	if iceCandidate == "" {
		// get lastly added signaling info to set is_last to 1
		lastSignalingInfo, _ := uc.SignalingRepository.FindSignalingInfo(orakkiId, "desc", 1)
		lastSignalingInfo.IsLast = true

		saved, err = uc.SignalingRepository.UpdateSignalingInfo(lastSignalingInfo)
	} else {
		saved, err = uc.SignalingRepository.SaveSignalingInfo(&SignalingInfo)
	}

	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (uc *SignalingUseCase) AddIceCandidate(orakkiId string, playerId int64, iceCandidate string) (*models.SignalingInfo, error) {
	game, _ := uc.GameRepository.FindRunningGameByOrakkiId(orakkiId)

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
