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

func (r *GameRepositoryMySqlImpl) GetById(id int) (*models.Game, error) {
	var game *models.Game

	result := dto.GameData{}
	err := r.DB.Get(&result, "select * from game where id = ? limit 1", id)
	if err != nil {
		return game, err
	}

	mapstructure.Decode(result, &game)

	return game, nil
}

func (r *GameRepositoryMySqlImpl) Find(offset, limit int) []*models.Game {
	var games []*models.Game

	result := []dto.GameData{}
	err := r.DB.Select(&result, "select * from game limit ? offset ?", limit, offset)
	if err != nil {
		return games
	}

	mapstructure.Decode(result, &games)

	return games
}

type RunningGameRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *RunningGameRepositoryMySqlImpl) Find(offset, limit int) []*models.RunningGame {
	return nil
}

func (r *RunningGameRepositoryMySqlImpl) FindById(id int64) (*models.RunningGame, error) {
	result := dto.RunningGameData{}
	err := r.DB.Get(&result, "select * from running_game where id = ? limit 1", id)
	if err != nil {
		return nil, err
	}

	game := models.RunningGame{Id: result.Id, PeerName: result.PeerName, Orakki: &models.Orakki{Id: result.OrakkiId}}

	return &game, nil
}

func (r *RunningGameRepositoryMySqlImpl) Save(game *models.RunningGame) (*models.RunningGame, error) {
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

type SignalingRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *SignalingRepositoryMySqlImpl) Save(signalingInfo *models.SignalingInfo) (*models.SignalingInfo, error) {
	data := dto.SignalingInfoData{
		OrakkiID: signalingInfo.Game.Orakki.Id,
		Data:     signalingInfo.Data,
		IsLast:   signalingInfo.IsLast,
	}

	var err error
	if signalingInfo.Id > 0 {
		updateQuery := `update signaling_info set is_last = ? where id = ? `
		_, err := r.DB.Exec(updateQuery, data.IsLast, signalingInfo.Id)
		if err != nil {
			return nil, err
		}
	} else {
		insertQuery := `insert into signaling_info (orakki_id, data, is_last) values (?, ?, ?)`
		result, err := r.DB.Exec(insertQuery, data.OrakkiID, data.Data, data.IsLast)
		if err != nil {
			return nil, err
		}

		LastInsertId, _ := result.LastInsertId()
		signalingInfo.Id = LastInsertId
	}

	return signalingInfo, err
}

func (r *SignalingRepositoryMySqlImpl) FindByRunningGameId(runningGameId int64, sinceId int64) (*models.SignalingInfo, error) {
	var signalingInfo *models.SignalingInfo
	result := dto.SignalingInfoData{}
	err := r.DB.Get(&result, "select * from signaling_info where running_game_id = ? and id > ? order by id asc", runningGameId, sinceId)
	if err != nil {
		return nil, err
	}

	mapstructure.Decode(result, &signalingInfo)

	return signalingInfo, nil
}
