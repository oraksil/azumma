package di

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/golobby/container"
	"github.com/jmoiron/sqlx"
	"github.com/sangwonl/mqrpc"
	"gitlab.com/oraksil/azumma/internal/domain/models"
	"gitlab.com/oraksil/azumma/internal/domain/services"
	"gitlab.com/oraksil/azumma/internal/domain/usecases"
	"gitlab.com/oraksil/azumma/internal/presenter/data"
	"gitlab.com/oraksil/azumma/internal/presenter/mq/handlers"
	"gitlab.com/oraksil/azumma/internal/presenter/web"
	"gitlab.com/oraksil/azumma/internal/presenter/web/ctrls"
)

func newWebService() *web.WebService {
	return web.NewWebService()
}

func newMqService() *mqrpc.MqService {
	svc, _ := mqrpc.NewMqService("amqp://oraksil:oraksil@localhost:5672/", "oraksil")
	return svc
}

func newMessageService() services.MessageService {
	var mqService *mqrpc.MqService
	container.Make(&mqService)

	return &mqrpc.DefaultMessageServiceImpl{MqService: mqService}
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

	var msgService services.MessageService
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
