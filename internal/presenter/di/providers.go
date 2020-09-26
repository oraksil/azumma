package di

import (
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golobby/container"
	"github.com/jmoiron/sqlx"
	"github.com/oraksil/azumma/internal/domain/models"
	"github.com/oraksil/azumma/internal/domain/services"
	"github.com/oraksil/azumma/internal/domain/usecases"
	"github.com/oraksil/azumma/internal/presenter/data"
	"github.com/oraksil/azumma/internal/presenter/mq/handlers"
	"github.com/oraksil/azumma/internal/presenter/web"
	"github.com/oraksil/azumma/internal/presenter/web/ctrls"
	"github.com/oraksil/azumma/pkg/drivers"
	"github.com/oraksil/azumma/pkg/utils"
	"github.com/sangwonl/mqrpc"
)

func newServiceConfig() *services.ServiceConfig {
	hostname, _ := os.Hostname()

	return &services.ServiceConfig{
		MqRpcUri:        utils.GetStrEnv("MQRPC_URI", "amqp://oraksil:oraksil@oraksil-mq-svc:5672/"),
		MqRpcNamespace:  utils.GetStrEnv("MQRPC_NAMESPACE", "oraksil"),
		MqRpcIdentifier: utils.GetStrEnv("MQRPC_IDENTIFIER", hostname),

		StaticOrakkiId: utils.GetStrEnv("STATIC_ORAKKI_ID", ""),

		OrakkiContainerImage: utils.GetStrEnv("ORAKKI_CONTAINER_IMAGE", "oraksil/orakki:latest"),
		GipanContainerImage:  utils.GetStrEnv("GIPAN_CONTAINER_IMAGE", "oraksil/gipan:latest"),
		ProvisionMaxWait:     time.Duration(utils.GetIntEnv("PROVISION_MAX_WAIT", 30)),
	}
}

func newOrakkiDriver() services.OrakkiDriver {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	drv, err := drivers.NewK8SOrakkiDriver(
		"",
		serviceConf.OrakkiContainerImage,
		serviceConf.GipanContainerImage,
		serviceConf.MqRpcUri,
		serviceConf.MqRpcNamespace,
	)
	if err != nil {
		panic(err)
	}
	return drv
}

func newWebService() *web.WebService {
	return web.NewWebService()
}

func newMqService() *mqrpc.MqService {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	svc, err := mqrpc.NewMqService(serviceConf.MqRpcUri, serviceConf.MqRpcNamespace)
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
	db, err := sqlx.Open("mysql", "oraksil:oraksil@(localhost:3306)/oraksil?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.DB.SetMaxOpenConns(10)
	_ = db.Ping()
	return db
}

func newPackRepository() models.PackRepository {
	var db *sqlx.DB
	container.Make(&db)

	return &data.PackRepositoryMySqlImpl{DB: db}
}

func newGameRepository() models.GameRepository {
	var db *sqlx.DB
	container.Make(&db)

	return &data.GameRepositoryMySqlImpl{DB: db}
}

func newSignalingRepository() models.SignalingRepository {
	var db *sqlx.DB
	container.Make(&db)

	return &data.SignalingRepositoryMySqlImpl{DB: db}
}

func newGameFetchUseCase() *usecases.GameFetchUseCase {
	var packRepo models.PackRepository
	container.Make(&packRepo)

	var gameRepo models.GameRepository
	container.Make(&gameRepo)

	return &usecases.GameFetchUseCase{
		PackRepo: packRepo,
		GameRepo: gameRepo,
	}
}

func newGameCtrlUseCase() *usecases.GameCtrlUseCase {
	var packRepo models.PackRepository
	container.Make(&packRepo)

	var gameRepo models.GameRepository
	container.Make(&gameRepo)

	var msgService services.MessageService
	container.Make(&msgService)

	var orakkiDrv services.OrakkiDriver
	container.Make(&orakkiDrv)

	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	return &usecases.GameCtrlUseCase{
		PackRepo:       packRepo,
		GameRepo:       gameRepo,
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
	var gameRepo models.GameRepository
	container.Make(&gameRepo)

	var signalingRepo models.SignalingRepository
	container.Make(&signalingRepo)

	var msgService services.MessageService
	container.Make(&msgService)

	return &usecases.SignalingUseCase{
		GameRepo:       gameRepo,
		SignalingRepo:  signalingRepo,
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

func newSignalingHandler() *handlers.SignalingHandler {
	var signalingUseCase *usecases.SignalingUseCase
	container.Make(&signalingUseCase)

	return &handlers.SignalingHandler{
		SignalingUseCase: signalingUseCase,
	}
}
