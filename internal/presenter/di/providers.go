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
		// Azumma
		DbUri: utils.GetStrEnv("DB_URI", "oraksil:oraksil@(localhost:3306)/oraksil?parseTime=true"),

		MqRpcUri:        utils.GetStrEnv("MQRPC_URI", "amqp://oraksil:oraksil@localhost:5672/"),
		MqRpcNamespace:  utils.GetStrEnv("MQRPC_NAMESPACE", "oraksil"),
		MqRpcIdentifier: utils.GetStrEnv("MQRPC_IDENTIFIER", hostname),

		StaticOrakkiId:   utils.GetStrEnv("STATIC_ORAKKI_ID", ""),
		ProvisionMaxWait: time.Duration(utils.GetIntEnv("PROVISION_MAX_WAIT", 30)),

		// For Orakki
		OrakkiMqRpcUri:       utils.GetStrEnv("ORAKKI_MQRPC_URI", "amqp://oraksil:oraksil@localhost:5672/"),
		OrakkiMqRpcNamespace: utils.GetStrEnv("ORAKKI_MQRPC_NAMESPACE", "oraksil"),

		OrakkiContainerImage: utils.GetStrEnv("ORAKKI_CONTAINER_IMAGE", "oraksil/orakki:latest"),
		GipanContainerImage:  utils.GetStrEnv("GIPAN_CONTAINER_IMAGE", "oraksil/gipan:latest"),

		TurnServerUri:      utils.GetStrEnv("TURN_URI", ""),
		TurnServerUsername: utils.GetStrEnv("TURN_USERNAME", ""),
		TurnServerPassword: utils.GetStrEnv("TURN_PASSWORD", ""),

		GipanResolution:       utils.GetStrEnv("GIPAN_RESOLUTION", "640x480"),
		GipanFps:              utils.GetStrEnv("GIPAN_FPS", "25"),
		GipanKeyframeInterval: utils.GetStrEnv("GIPAN_KEYFRAME_INTERVAL", "150"),

		OrakkiDriverK8SConfigPath:        utils.GetStrEnv("ORAKKI_DRIVER_K8S_CONFIG_PATH", ""),
		OrakkiDriverK8SNamespace:         utils.GetStrEnv("ORAKKI_DRIVER_K8S_NAMESPACE", ""),
		OrakkiDriverK8SNodeSelectorKey:   utils.GetStrEnv("ORAKKI_DRIVER_K8S_NODE_SELECTOR_KEY", ""),
		OrakkiDriverK8SNodeSelectorValue: utils.GetStrEnv("ORAKKI_DRIVER_K8S_NODE_SELECTOR_VALUE", ""),
	}
}

func newOrakkiDriver() services.OrakkiDriver {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	drv, err := drivers.NewK8SOrakkiDriver(
		serviceConf.OrakkiDriverK8SConfigPath,
		serviceConf.OrakkiDriverK8SNamespace,
		serviceConf.OrakkiDriverK8SNodeSelectorKey,
		serviceConf.OrakkiDriverK8SNodeSelectorValue,
		serviceConf.OrakkiContainerImage,
		serviceConf.GipanContainerImage,
		serviceConf.OrakkiMqRpcUri,
		serviceConf.OrakkiMqRpcNamespace,
		serviceConf.TurnServerUri,
		serviceConf.TurnServerUsername,
		serviceConf.TurnServerPassword,
		serviceConf.GipanResolution,
		serviceConf.GipanFps,
		serviceConf.GipanKeyframeInterval,
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
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	db, err := sqlx.Open("mysql", serviceConf.DbUri)
	if err != nil {
		panic(err)
	}

	db.DB.SetMaxOpenConns(10)
	_ = db.Ping()
	return db
}

func newPlayerRepository() models.PlayerRepository {
	var db *sqlx.DB
	container.Make(&db)

	return &data.PlayerRepositoryMySqlImpl{DB: db}
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

func newUserFeedbackRepository() models.UserFeedbackRepository {
	var db *sqlx.DB
	container.Make(&db)

	return &data.UserFeedbackRepositoryMySqlImpl{DB: db}
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
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	var packRepo models.PackRepository
	container.Make(&packRepo)

	var gameRepo models.GameRepository
	container.Make(&gameRepo)

	var playerRepo models.PlayerRepository
	container.Make(&playerRepo)

	var msgService services.MessageService
	container.Make(&msgService)

	var orakkiDrv services.OrakkiDriver
	container.Make(&orakkiDrv)

	return &usecases.GameCtrlUseCase{
		ServiceConfig:  serviceConf,
		PackRepo:       packRepo,
		GameRepo:       gameRepo,
		PlayerRepo:     playerRepo,
		MessageService: msgService,
		OrakkiDriver:   orakkiDrv,
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

func newGameHandler() *handlers.GameHandler {
	var gameCtrlUseCase *usecases.GameCtrlUseCase
	container.Make(&gameCtrlUseCase)

	return &handlers.GameHandler{GameCtrlUseCase: gameCtrlUseCase}
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

func newPlayerUseCase() *usecases.PlayerUseCase {
	var playerRepo models.PlayerRepository
	container.Make(&playerRepo)

	return &usecases.PlayerUseCase{
		PlayerRepo: playerRepo,
	}
}

func newPlayerController() *ctrls.PlayerController {
	var playerUseCase *usecases.PlayerUseCase
	container.Make(&playerUseCase)

	return &ctrls.PlayerController{
		PlayerUseCase: playerUseCase,
	}
}

func newUserFeedbackUseCase() *usecases.UserFeedbackUseCase {
	var feedbackRepo models.UserFeedbackRepository
	container.Make(&feedbackRepo)

	return &usecases.UserFeedbackUseCase{
		FeedbackRepo: feedbackRepo,
	}
}

func newUserFeedbackController() *ctrls.UserFeedbackController {
	var feedbackUseCase *usecases.UserFeedbackUseCase
	container.Make(&feedbackUseCase)

	return &ctrls.UserFeedbackController{
		UserFeedbackUseCase: feedbackUseCase,
	}
}
