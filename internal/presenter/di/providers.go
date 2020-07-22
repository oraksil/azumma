package di

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golobby/container"
	"github.com/jmoiron/sqlx"
	"gitlab.com/oraksil/sil/backend/internal/domain/models"
	"gitlab.com/oraksil/sil/backend/internal/domain/usecases"
	"gitlab.com/oraksil/sil/backend/internal/presenter/data"
	"gitlab.com/oraksil/sil/backend/internal/presenter/mq/handlers"
	"gitlab.com/oraksil/sil/backend/internal/presenter/web"
	"gitlab.com/oraksil/sil/backend/internal/presenter/web/ctrls"
	"gitlab.com/oraksil/sil/backend/pkg/mq"
)

func newWebService() *web.WebService {
	return web.NewWebService()
}

func newMqService() *mq.MqService {
	return mq.NewMqService("amqp://oraksil:oraksil@localhost:5672/", "oraksil.mq.p2p", "oraksil.mq.broadcast")
}

func newMessageService() mq.MessageService {
	var mqService *mq.MqService
	container.Make(&mqService)

	return &mq.DefaultMessageServiceImpl{MqService: mqService}
}

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

func newGameCtrlUseCase() *usecases.GameCtrlUseCase {
	var repo models.GameRepository
	container.Make(&repo)

	var msgService mq.MessageService
	container.Make(&msgService)

	return &usecases.GameCtrlUseCase{GameRepository: repo, MessageService: msgService}
}

func newGameController() *ctrls.GameController {
	var gameFetchUseCase *usecases.GameFetchUseCase
	container.Make(&gameFetchUseCase)

	var gameCtrlUseCase *usecases.GameCtrlUseCase
	container.Make(&gameCtrlUseCase)

	return &ctrls.GameController{
		GameFetchUseCase: gameFetchUseCase,
		GameCtrlUseCase:  gameCtrlUseCase,
	}
}

func newHelloHandler() *handlers.HelloHandler {
	var gameCtrlUseCase *usecases.GameCtrlUseCase
	container.Make(&gameCtrlUseCase)

	return &handlers.HelloHandler{
		GameCtrlUseCase: gameCtrlUseCase,
	}
}
