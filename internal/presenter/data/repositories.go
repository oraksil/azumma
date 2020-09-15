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
	err := r.DB.Get(&result, "select * from pack where id = ? limit 1", id)
	if err != nil {
		return pack, err
	}

	mapstructure.Decode(result, &pack)

	return pack, nil
}

func (r *PackRepositoryMySqlImpl) Find(offset, limit int) []*models.Pack {
	var packs []*models.Pack

	result := []dto.PackData{}
	err := r.DB.Select(&result, "select * from pack limit ? offset ?", limit, offset)
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
	err := r.DB.Get(&result, "select * from game where id = ? limit 1", id)
	if err != nil {
		return nil, err
	}

	game := models.Game{
		Id:       result.Id,
		PeerName: result.PeerName,
		Orakki:   &models.Orakki{Id: result.OrakkiId},
	}

	return &game, nil
}

func (r *GameRepositoryMySqlImpl) Save(game *models.Game) (*models.Game, error) {
	// map models to dto
	data := dto.GameData{
		OrakkiId:      game.Orakki.Id,
		OrakkiState:   game.Orakki.State,
		PeerName:      game.Orakki.PeerName,
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
		insert into game (
			peer_name,
			orakki_id,
			orakki_state,
			pack_id,
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
		OrakkiId: signaling.Game.Orakki.Id,
		Data:     signaling.Data,
		IsLast:   signaling.IsLast,
	}

	var err error
	if signaling.Id > 0 {
		updateQuery := `update signaling set is_last = ? where id = ? `
		_, err := r.DB.Exec(updateQuery, data.IsLast, signaling.Id)
		if err != nil {
			return nil, err
		}
	} else {
		insertQuery := `insert into signaling (orakki_id, data, is_last) values (?, ?, ?)`
		result, err := r.DB.Exec(insertQuery, data.OrakkiId, data.Data, data.IsLast)
		if err != nil {
			return nil, err
		}

		LastInsertId, _ := result.LastInsertId()
		signaling.Id = LastInsertId
	}

	return signaling, err
}

func (r *SignalingRepositoryMySqlImpl) FindByGameId(gameId int64, sinceId int64) (*models.Signaling, error) {
	var signaling *models.Signaling
	result := dto.SignalingData{}
	err := r.DB.Get(&result, "select * from signaling where game_id = ? and id > ? order by id asc", gameId, sinceId)
	if err != nil {
		return nil, err
	}

	mapstructure.Decode(result, &signaling)

	return signaling, nil
}
