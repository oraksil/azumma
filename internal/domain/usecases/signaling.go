package usecases

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"
	"github.com/oraksil/azumma/pkg/utils"

	"time"
)

type SignalingUseCase struct {
	GameRepo       models.GameRepository
	SignalingRepo  models.SignalingRepository
	MessageService services.MessageService
}

func (uc *SignalingUseCase) NewOffer(gameId int64, playerId int64, b64EncodedOffer string) (*models.SdpInfo, error) {
	// validation
	var offer map[string]interface{}
	err := utils.DecodeFromB64EncodedJsonStr(b64EncodedOffer, &offer)
	if err != nil {
		return nil, err
	}

	if v, ok := offer["type"]; !ok || v.(string) != "offer" {
		return nil, errors.New("invalid sdp type")
	}

	game, err := uc.GameRepo.FindById(gameId)
	if game == nil {
		return nil, errors.New("no game exists with given gameId")
	}

	// sdp response from orakki
	resp, err := uc.MessageService.Request(
		game.Orakki.Id,
		models.MsgSetupWithNewOffer,
		models.SdpInfo{PeerId: playerId, SdpBase64Encoded: b64EncodedOffer},
		10*time.Second,
	)

	var answerSdpInfo models.SdpInfo
	mapstructure.Decode(resp, &answerSdpInfo)

	return &answerSdpInfo, err
}

func (uc *SignalingUseCase) GetOrakkiIceCandidates(gameId int64, lastSeq int64) ([]*models.IceCandidate, error) {
	signalings, err := uc.SignalingRepo.FindByGameId(gameId, lastSeq)
	if err != nil {
		return nil, err
	}

	iceCandidates := make([]*models.IceCandidate, 0, 0)

	for _, s := range signalings {
		ice := &models.IceCandidate{
			PeerId:           s.GameId,
			IceBase64Encoded: s.Data,
			Seq:              s.Id,
		}

		iceCandidates = append(iceCandidates, ice)
	}

	return iceCandidates, nil
}

func (uc *SignalingUseCase) OnOrakkiIceCandidate(gameId int64, iceBase64Encoded string) error {
	game, err := uc.GameRepo.FindById(gameId)
	if game == nil {
		return errors.New("no game matched to given gameId")
	}

	signaling := models.Signaling{
		GameId: game.Id,
		Data:   iceBase64Encoded,
	}

	_, err = uc.SignalingRepo.Save(&signaling)
	if err != nil {
		return err
	}

	return nil
}

func (uc *SignalingUseCase) OnPlayerIceCandidate(gameId int64, playerId int64, b64EncodedIceCandidate string) error {
	game, _ := uc.GameRepo.FindById(gameId)
	if game == nil {
		return errors.New("no game exists with given gameId")
	}

	_, err := uc.MessageService.Request(
		game.Orakki.Id,
		models.MsgRemoteIceCandidate,
		models.IceCandidate{
			PeerId:           playerId,
			IceBase64Encoded: b64EncodedIceCandidate,
		},
		5*time.Second,
	)

	if err != nil {
		return err
	}

	return nil
}
