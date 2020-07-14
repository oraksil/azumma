package di

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golobby/container"
	"github.com/jmoiron/sqlx"
	"oraksil.com/sil/internal/domain/models"
	"oraksil.com/sil/internal/domain/usecases"
	"oraksil.com/sil/internal/presenter/data"
	"oraksil.com/sil/internal/presenter/web/ctrls"
)

func newMySqlDb() *sqlx.DB {
	db, _ := sqlx.Open("mysql", "oraksil:qlqjswha!@(localhost:3306)/oraksil")
	db.DB.SetMaxOpenConns(10)
	_ = db.Ping()
	return db
}

func newGameRepository() models.GameRepository {
	var db *sqlx.DB
	container.Make(&db)

	return &data.GameRepositoryMySqlImpl{DB: db}
}

func newGameFetchUseCase() *usecases.GameFetchUseCase {
	var repo models.GameRepository
	container.Make(&repo)

	return &usecases.GameFetchUseCase{GameRepository: repo}
}

func newGameController() *ctrls.GameController {
	var gameFetchUseCase *usecases.GameFetchUseCase
	container.Make(&gameFetchUseCase)

	return &ctrls.GameController{GameFetchUseCase: gameFetchUseCase}
}
