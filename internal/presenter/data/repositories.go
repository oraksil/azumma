package data

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/presenter/data/dto"
)

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

func (r *PackRepositoryMySqlImpl) Find(offset, limit int) []*models.Pack {
	var packs []*models.Pack

	result := []dto.PackData{}
	err := r.DB.Select(&result, "SELECT * FROM pack LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return packs
	}

	mapstructure.Decode(result, &packs)

	return packs
}

type GameRepositoryMySqlImpl struct {
	DB *sqlx.DB
}

func (r *GameRepositoryMySqlImpl) Find(offset, limit int) []*models.Game {
	return nil
}

func (r *GameRepositoryMySqlImpl) FindById(id int64) (*models.Game, error) {
	result := dto.GameData{}
	err := r.DB.Get(&result, "SELECT * FROM game WHERE id = ? LIMIT 1", id)
	if err != nil {
		return nil, err
	}

	game := models.Game{
		Id: result.Id,
		Orakki: &models.Orakki{
			Id:    result.OrakkiId,
			State: result.OrakkiState,
		},
	}

	return &game, nil
}

func (r *GameRepositoryMySqlImpl) Save(game *models.Game) (*models.Game, error) {
	// map models to dto
	data := dto.GameData{
		OrakkiId:      game.Orakki.Id,
		OrakkiState:   game.Orakki.State,
		PackId:        game.Pack.Id,
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
		INSERT INTO game (
			orakki_id,
			orakki_state,
			pack_id,
			first_player_id,
			joined_player_ids,
			created_at)
		VALUES
			(?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			orakki_state = ?,
			first_player_id = ?,
			joined_player_ids = ?`

	result, err := r.DB.Exec(
		// insert args
		insertQuery,
		data.OrakkiId,
		data.OrakkiState,
		data.PackId,
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

func (r *SignalingRepositoryMySqlImpl) Save(signaling *models.Signaling) (*models.Signaling, error) {
	data := dto.SignalingData{
		GameId: signaling.GameId,
		Data:   signaling.Data,
	}

	var err error
	insertQuery := `INSERT INTO signaling (game_id, data) VALUES (?, ?)`
	result, err := r.DB.Exec(insertQuery, data.GameId, data.Data)
	if err != nil {
		return nil, err
	}

	LastInsertId, _ := result.LastInsertId()
	signaling.Id = LastInsertId

	return signaling, err
}

func (r *SignalingRepositoryMySqlImpl) FindOneByGameId(gameId int64, sinceId int64) (*models.Signaling, error) {
	var signalings models.Signaling

	result := dto.SignalingData{}
	query := "SELECT * FROM signaling WHERE game_id = ? AND id > ? ORDER BY id ASC LIMIT 1"
	err := r.DB.Get(&result, query, gameId, sinceId)
	if err != nil {
		return nil, err
	}

	mapstructure.Decode(result, &signalings)

	return &signalings, nil
}
