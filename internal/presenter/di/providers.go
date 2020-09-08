package di

import (
	"os"
	"time"

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
	"gitlab.com/oraksil/azumma/pkg/drivers"
	"gitlab.com/oraksil/azumma/pkg/utils"
)

func newServiceConfig() *services.ServiceConfig {
	useStaticOrakki := utils.GetBoolEnv("USE_STATIC_ORAKKI", false)
	hostname, _ := os.Hostname()
	return &services.ServiceConfig{
		UseStaticOrakki:      useStaticOrakki,
		StaticOrakkiId:       utils.GetStrEnv("STATIC_ORAKKI_ID", "orakki-static"),
		StaticOrakkiPeerName: utils.GetStrEnv("STATIC_ORAKKI_PEER_NAME", "orakki-local"),
		PeerName:             utils.GetStrEnv("PEER_NAME", hostname),
		ProvisionMaxWait:     time.Duration(utils.GetIntEnv("PROVISION_MAX_WAIT", 30)),
	}
}

func newOrakkiDriver() services.OrakkiDriver {
	drv, err := drivers.NewK8SOrakkiDriver("", "orakki:latest")
	if err != nil {
		panic(err)
	}
	return drv
}

func newWebService() *web.WebService {
	return web.NewWebService()
}

func newMqService() *mqrpc.MqService {
	svc, err := mqrpc.NewMqService("amqp://oraksil:oraksil@localhost:5672/", "oraksil")
	if err != nil {
		panic(err)
	}
	return svc
}

func newMessageService() services.MessageService {
	var mqService *mqrpc.MqService
	container.Make(&mqService)

	return &mqrpc.DefaultMessageServiceImpl{MqService: mqService}
}

func newMySqlDb() *sqlx.DB {
	db, err := sqlx.Open("mysql", "oraksil:qlqjswha!@(localhost:3306)/oraksil?parseTime=true")
	if err != nil {
		panic(err)
	}

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

	var orakkiDrv services.OrakkiDriver
	container.Make(&orakkiDrv)

	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	return &usecases.GameCtrlUseCase{
		GameRepository: repo,
		MessageService: msgService,
		OrakkiDriver:   orakkiDrv,
		ServiceConfig:  serviceConf,
	}
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

func newSignalingUseCases() *usecases.SignalingUseCase {
	var repo models.GameRepository
	container.Make(&repo)

	var msgService services.MessageService
	container.Make(&msgService)

	return &usecases.SignalingUseCase{
		GameRepository: repo,
		MessageService: msgService,
	}
}

func newSignalingController() *ctrls.SignalingController {
	var signalingUseCase *usecases.SignalingUseCase
	container.Make(&signalingUseCase)

	return &ctrls.SignalingController{
		SignalingUseCase: signalingUseCase,
	}
}

func newHelloHandler() *handlers.HelloHandler {
	var gameCtrlUseCase *usecases.GameCtrlUseCase
	container.Make(&gameCtrlUseCase)

	return &handlers.HelloHandler{
		GameCtrlUseCase: gameCtrlUseCase,
	}
}
