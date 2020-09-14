package data

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/presenter/data/dto"
)

type GameRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *GameRepositoryMySqlImpl) GetGameById(id int) (*models.Game, error) {
	var game *models.Game

	result := dto.GameData{}
	err := r.DB.Get(&result, "select * from game where id = ? limit 1", id)
	if err != nil {
		return game, err
	}

	mapstructure.Decode(result, &game)

	return game, nil
}

func (r *GameRepositoryMySqlImpl) FindAvailableGames(offset, limit int) []*models.Game {
	var games []*models.Game

	result := []dto.GameData{}
	err := r.DB.Select(&result, "select * from game limit ? offset ?", limit, offset)
	if err != nil {
		return games
	}

	mapstructure.Decode(result, &games)

	return games
}

func (r *GameRepositoryMySqlImpl) FindRunningGames(offset, limit int) []*models.RunningGame {
	return nil
}

func (r *GameRepositoryMySqlImpl) SaveRunningGame(game *models.RunningGame) (*models.RunningGame, error) {
	// map models to dto
	data := dto.RunningGameData{
		OrakkiId:      game.Orakki.Id,
		OrakkiState:   game.Orakki.State,
		PeerName:      game.Orakki.PeerName,
		GameId:        game.Game.Id,
		FirstPlayerId: game.Players[0].Id,
		CreatedAt:     time.Now(),
	}

	playerNames := make([]string, len(game.Players))
	for i, p := range game.Players {
		playerNames[i] = p.Name
	}
	data.JoinedPlayerIds = strings.Join(playerNames, ",")

	// insert and return id aware model
	insertQuery := `
		insert into running_game (
			peer_name,
			orakki_id,
			orakki_state,
			game_id,
			first_player_id,
			joined_player_ids,
			created_at)
		values
			(?, ?, ?, ?, ?, ?, ?)
		on duplicate key update
			orakki_state = ?,
			first_player_id = ?,
			joined_player_ids = ?`

	result, err := r.DB.Exec(
		// insert args
		insertQuery,
		data.PeerName,
		data.OrakkiId,
		data.OrakkiState,
		data.GameId,
		data.FirstPlayerId,
		data.JoinedPlayerIds,
		data.CreatedAt,
		// update args on duplicate key
		data.OrakkiState,
		data.FirstPlayerId,
		data.JoinedPlayerIds,
	)

	if err != nil {
		return nil, err
	}

	lastInsertedId, _ := result.LastInsertId()
	game.Id = lastInsertedId
	game.CreatedAt = data.CreatedAt

	return game, nil
}

func (r *GameRepositoryMySqlImpl) FindRunningGameByOrakkiId(orakkiId string) (*models.RunningGame, error) {

	result := dto.RunningGameData{}
	err := r.DB.Get(&result, "select * from running_game where orakki_id = ? limit 1", orakkiId)
	if err != nil {
		return nil, err
	}

	game := models.RunningGame{Id: result.Id, PeerName: result.PeerName, Orakki: &models.Orakki{Id: result.OrakkiId}}

	return &game, nil
}

type SignalingRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *SignalingRepositoryMySqlImpl) SaveSignalingInfo(signalingInfo *models.SignalingInfo) (*models.SignalingInfo, error) {
	data := dto.SignalingInfoData{
		OrakkiID: signalingInfo.Game.Orakki.Id,
		Data:     signalingInfo.Data,
		IsLast:   signalingInfo.IsLast,
	}

	insertQuery := `insert into signaling_info (
		orakki_id,
		data,
		is_last) values (?, ?, ?) `

	result, err := r.DB.Exec(insertQuery, data.OrakkiID, data.Data, data.IsLast)

	if err != nil {
		return nil, err
	}

	LastInsertId, _ := result.LastInsertId()
	signalingInfo.Id = LastInsertId

	return signalingInfo, err
}
func (r *SignalingRepositoryMySqlImpl) UpdateSignalingInfo(signalingInfo *models.SignalingInfo) (*models.SignalingInfo, error) {
	data := dto.SignalingInfoData{
		OrakkiID: signalingInfo.OrakkiId,
		Data:     signalingInfo.Data,
		IsLast:   signalingInfo.IsLast,
	}

	insertQuery := `update signaling_info set 
		is_last = ? where id = ? `

	result, err := r.DB.Exec(insertQuery, data.IsLast, signalingInfo.Id)

	if err != nil {
		return nil, err
	}

	LastInsertId, _ := result.LastInsertId()
	signalingInfo.Id = LastInsertId

	return signalingInfo, err
}
func (r *SignalingRepositoryMySqlImpl) FindSignalingInfo(orakkiId string, order string, num int) (*models.SignalingInfo, error) {
	var signalingInfo *models.SignalingInfo
	result := dto.SignalingInfoData{}
	err := r.DB.Get(&result, "select * from signaling_info where orakki_id = ? order by id desc limit ?", orakkiId, num)
	if err != nil {
		return nil, err
	}

	mapstructure.Decode(result, &signalingInfo)

	return signalingInfo, nil
}

func (r *SignalingRepositoryMySqlImpl) FindIceCandidate(orakkiId string, seqAfter int) (*models.SignalingInfo, error) {
	var signalingInfo *models.SignalingInfo
	result := dto.SignalingInfoData{}
	err := r.DB.Get(&result, "select * from signaling_info where orakki_id = ? and id > ? order by id asc", orakkiId, seqAfter)
	if err != nil {
		return nil, err
	}

	mapstructure.Decode(result, &signalingInfo)

	return signalingInfo, nil

}
