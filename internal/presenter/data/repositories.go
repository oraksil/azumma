package data

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/presenter/data/dto"
)

type PlayerRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *PlayerRepositoryMySqlImpl) GetById(id int64) (*models.Player, error) {
	var player *models.Player

	playerData := dto.PlayerData{}
	err := r.DB.Get(&playerData, "SELECT * FROM player WHERE id = ? LIMIT 1", id)
	if err != nil {
		return player, err
	}

	mapstructure.Decode(playerData, &player)

	return player, nil
}

func (r *PlayerRepositoryMySqlImpl) Save(player *models.Player) (*models.Player, error) {
	data := dto.PlayerData{
		Name:       player.Name,
		TotalCoins: player.TotalCoins,
	}

	if player.Id > 0 {
		updateQuery := `UPDATE player SET name = ?, total_coins = ? WHERE id = ?`
		_, err := r.DB.Exec(updateQuery, data.Name, data.TotalCoins, player.Id)
		if err != nil {
			return nil, err
		}
	} else {
		insertQuery := `INSERT INTO player (name, total_coins) VALUES (?, ?)`
		result, err := r.DB.Exec(insertQuery, data.Name, data.TotalCoins)
		if err != nil {
			return nil, err
		}

		lastInsertId, _ := result.LastInsertId()
		player.Id = lastInsertId
	}

	return player, nil
}

type PackRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *PackRepositoryMySqlImpl) GetById(id int) (*models.Pack, error) {
	var pack *models.Pack

	result := dto.PackData{}
	err := r.DB.Get(&result, "SELECT * FROM pack WHERE id = ? LIMIT 1", id)
	if err != nil {
		return pack, err
	}

	mapstructure.Decode(result, &pack)

	return pack, nil
}

func (r *PackRepositoryMySqlImpl) FindByStatus(status, offset, limit int) []*models.Pack {
	var packs []*models.Pack

	result := []dto.PackData{}

	query := "SELECT * FROM pack WHERE status = ? LIMIT ? OFFSET ?"
	err := r.DB.Select(&result, query, status, limit, offset)
	if err != nil {
		return packs
	}

	mapstructure.Decode(result, &packs)

	return packs

}

type GameRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *GameRepositoryMySqlImpl) GetById(id int64) (*models.Game, error) {
	gameData := dto.GameData{}
	err := r.DB.Get(&gameData, "SELECT * FROM game WHERE id = ? LIMIT 1", id)
	if err != nil {
		return nil, err
	}

	packData := dto.PackData{}
	err = r.DB.Get(&packData, "SELECT * FROM pack WHERE id = ? LIMIT 1", gameData.PackId)
	if err != nil {
		return nil, err
	}

	game := models.Game{
		Id: gameData.Id,
		Orakki: &models.Orakki{
			Id:    gameData.OrakkiId,
			State: gameData.OrakkiState,
		},
		Pack: &models.Pack{
			Id:         packData.Id,
			Title:      packData.Title,
			MaxPlayers: packData.MaxPlayers,
		},
		Players:   gameData.GetJoinedPlayers(),
		CreatedAt: gameData.CreatedAt,
	}

	return &game, nil
}

func (r *GameRepositoryMySqlImpl) Save(game *models.Game) (*models.Game, error) {
	// map models to dto
	data := dto.GameData{
		PackId:      game.Pack.Id,
		OrakkiId:    game.Orakki.Id,
		OrakkiState: game.Orakki.State,
		CreatedAt:   time.Now(),
	}
	data.SetJoinedPlayers(game.Players)

	if game.Id > 0 {
		updateQuery := `
			UPDATE game
				SET orakki_state = ?,
				joined_player_ids = ?
			WHERE id = ?`

		_, err := r.DB.Exec(
			updateQuery,
			data.OrakkiState,
			data.JoinedPlayerIds,
			game.Id,
		)

		if err != nil {
			return nil, err
		}
	} else {
		data.FirstPlayerId = game.Players[0].Id

		insertQuery := `
			INSERT INTO game (
				pack_id,
				orakki_id,
				orakki_state,
				first_player_id,
				joined_player_ids,
				created_at)
			VALUES
				(?, ?, ?, ?, ?, ?)`

		result, err := r.DB.Exec(
			insertQuery,
			data.PackId,
			data.OrakkiId,
			data.OrakkiState,
			data.FirstPlayerId,
			data.JoinedPlayerIds,
			data.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		lastInsertId, _ := result.LastInsertId()
		game.Id = lastInsertId
		game.CreatedAt = data.CreatedAt
	}

	return game, nil
}

type SignalingRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *SignalingRepositoryMySqlImpl) FindByToken(token string, sinceId int64) ([]*models.Signaling, error) {
	var signalings []*models.Signaling

	result := []*dto.SignalingData{}
	query := "SELECT * FROM signaling WHERE token = ? AND id > ? ORDER BY id ASC"
	err := r.DB.Select(&result, query, token, sinceId)
	if err != nil {
		return nil, err
	}

	mapstructure.Decode(result, &signalings)

	return signalings, nil
}

func (r *SignalingRepositoryMySqlImpl) Save(signaling *models.Signaling) (*models.Signaling, error) {
	data := dto.SignalingData{
		Token:    signaling.Token,
		GameId:   signaling.GameId,
		PlayerId: signaling.PlayerId,
		Data:     signaling.Data,
	}

	var err error
	insertQuery := `INSERT INTO signaling (token, game_id, player_id, data) VALUES (?, ?, ?, ?)`
	result, err := r.DB.Exec(insertQuery, data.Token, data.GameId, data.PlayerId, data.Data)
	if err != nil {
		return nil, err
	}

	lastInsertId, _ := result.LastInsertId()
	signaling.Id = lastInsertId

	return signaling, err
}
