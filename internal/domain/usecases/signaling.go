package usecases

import (
	"encoding/json"
	"errors"

	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"
	"github.com/oraksil/azumma/pkg/utils"

	"time"
)

type SignalingUseCase struct {
	ServiceConfig  *services.ServiceConfig
	MessageService services.MessageService
	GameRepo       models.GameRepository
	SignalingRepo  models.SignalingRepository
}

func (uc *SignalingUseCase) NewOffer(gameId int64, token, b64EncodedOffer string, sessionCtx services.SessionContext) (*models.SdpInfo, error) {
	// validation
	var offer map[string]interface{}
	err := utils.DecodeFromB64EncodedJsonStr(b64EncodedOffer, &offer)
	if err != nil {
		return nil, err
	}

	if v, ok := offer["type"]; !ok || v.(string) != "offer" {
		return nil, errors.New("invalid sdp type")
	}

	game, err := uc.GameRepo.GetById(gameId)
	if game == nil {
		return nil, errors.New("no game exists with given gameId")
	}

	if game.Orakki.State != models.OrakkiStateReady {
		return nil, errors.New("game machine is not ready for signaling")
	}

	// sdp response from orakki
	session, _ := sessionCtx.GetSession()
	resp, err := uc.MessageService.Request(
		game.Orakki.Id,
		models.MsgSetupWithNewOffer,
		models.SdpInfo{
			Peer: models.PeerInfo{
				Token:    token,
				GameId:   game.Id,
				PlayerId: session.Player.Id,
			},
			SdpBase64Encoded: b64EncodedOffer,
		},
		10*time.Second,
	)

	var answerSdpInfo models.SdpInfo
	respBytes, ok := resp.([]byte)
	if ok {
		json.Unmarshal(respBytes, &answerSdpInfo)
	}

	return &answerSdpInfo, err
}

func (uc *SignalingUseCase) GetOrakkiIceCandidates(
	token string, lastSeq int64, sessionCtx services.SessionContext) ([]*models.IceCandidate, error) {

	session, _ := sessionCtx.GetSession()

	signalings, err := uc.SignalingRepo.FindByToken(token, lastSeq)
	if err != nil {
		return nil, err
	}

	iceCandidates := make([]*models.IceCandidate, 0, 0)

	for _, s := range signalings {
		ice := &models.IceCandidate{
			Peer: models.PeerInfo{
				Token:    s.Token,
				GameId:   s.GameId,
				PlayerId: session.Player.Id,
			},
			IceBase64Encoded: s.Data,
			Seq:              s.Id,
		}

		iceCandidates = append(iceCandidates, ice)
	}

	return iceCandidates, nil
}

func (uc *SignalingUseCase) OnOrakkiIceCandidate(iceCandidate models.IceCandidate) error {
	game, err := uc.GameRepo.GetById(iceCandidate.Peer.GameId)
	if game == nil {
		return errors.New("no game matched to given gameId")
	}

	signaling := models.Signaling{
		Token:    iceCandidate.Peer.Token,
		GameId:   game.Id,
		PlayerId: iceCandidate.Peer.PlayerId,
		Data:     iceCandidate.IceBase64Encoded,
	}

	_, err = uc.SignalingRepo.Save(&signaling)
	if err != nil {
		return err
	}

	return nil
}

func (uc *SignalingUseCase) OnPlayerIceCandidate(
	gameId int64, token, b64EncodedIceCandidate string, sessionCtx services.SessionContext) error {

	game, _ := uc.GameRepo.GetById(gameId)
	if game == nil {
		return errors.New("no game exists with given gameId")
	}

	session, _ := sessionCtx.GetSession()
	_, err := uc.MessageService.Request(
		game.Orakki.Id,
		models.MsgRemoteIceCandidate,
		models.IceCandidate{
			Peer: models.PeerInfo{
				Token:    token,
				GameId:   game.Id,
				PlayerId: session.Player.Id,
			},
			IceBase64Encoded: b64EncodedIceCandidate,
		},
		5*time.Second,
	)

	if err != nil {
		return err
	}

	return nil
}

func (uc *SignalingUseCase) CreateUserAuth(userId string) (*models.TurnAuth, error) {
	username, password := utils.NewTurnAuth(
		userId,
		uc.ServiceConfig.TurnServerSecretKey,
		uc.ServiceConfig.TurnServerTTL,
	)

	return &models.TurnAuth{
		Username: username,
		Password: password,
	}, nil
}
