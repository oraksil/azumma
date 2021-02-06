package di

import (
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

		TurnServerUri:       utils.GetStrEnv("TURN_URI", ""),
		TurnServerSecretKey: utils.GetStrEnv("TURN_SECRET_KEY", ""),
		TurnServerTTL:       utils.GetIntEnv("TURN_TTL", 3600),

		PlayerHealthCheckTimeout: utils.GetIntEnv("PLAYER_HEALTHCHECK_TIMEOUT", 20),
		PlayerIdleCheckTimeout:   utils.GetIntEnv("PLAYER_IDLECHECK_TIMEOUT", 600),

		GipanResourceCpu:      utils.GetIntEnv("GIPAN_RESOURCE_CPU", 400),
		GipanResourceMemory:   utils.GetIntEnv("GIPAN_RESOURCE_MEMORY", 896),
		GipanResolution:       utils.GetStrEnv("GIPAN_RESOLUTION", "640x480"),
		GipanFps:              utils.GetStrEnv("GIPAN_FPS", "25"),
		GipanKeyframeInterval: utils.GetStrEnv("GIPAN_KEYFRAME_INTERVAL", "150"),

		OrakkiDriverK8SConfigPath:        utils.GetStrEnv("ORAKKI_DRIVER_K8S_CONFIG_PATH", ""),
		OrakkiDriverK8SNamespace:         utils.GetStrEnv("ORAKKI_DRIVER_K8S_NAMESPACE", ""),
		OrakkiDriverK8SNodeSelectorKey:   utils.GetStrEnv("ORAKKI_DRIVER_K8S_NODE_SELECTOR_KEY", ""),
		OrakkiDriverK8SNodeSelectorValue: utils.GetStrEnv("ORAKKI_DRIVER_K8S_NODE_SELECTOR_VALUE", ""),
	}
}

func newOrakkiDriver(serviceConf *services.ServiceConfig) services.OrakkiDriver {
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
		serviceConf.TurnServerSecretKey,
		serviceConf.TurnServerTTL,
		serviceConf.PlayerHealthCheckTimeout,
		serviceConf.PlayerIdleCheckTimeout,
		serviceConf.GipanResourceCpu,
		serviceConf.GipanResourceMemory,
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

func newMqService(serviceConf *services.ServiceConfig) *mqrpc.MqService {
	svc, err := mqrpc.NewMqService(serviceConf.MqRpcUri, serviceConf.MqRpcNamespace, serviceConf.MqRpcIdentifier)
	if err != nil {
		panic(err)
	}
	return svc
}

func newMessageService(mqService *mqrpc.MqService) services.MessageService {
	return &mqrpc.DefaultMessageService{MqService: mqService}
}

func newMySqlDb(serviceConf *services.ServiceConfig) *sqlx.DB {
	db, err := sqlx.Open("mysql", serviceConf.DbUri)
	if err != nil {
		panic(err)
	}

	db.DB.SetMaxOpenConns(10)
	_ = db.Ping()
	return db
}

func newPlayerRepository(db *sqlx.DB) models.PlayerRepository {
	return &data.PlayerRepositoryMySqlImpl{DB: db}
}

func newPackRepository(db *sqlx.DB) models.PackRepository {
	return &data.PackRepositoryMySqlImpl{DB: db}
}

func newGameRepository(db *sqlx.DB) models.GameRepository {
	return &data.GameRepositoryMySqlImpl{DB: db}
}

func newSignalingRepository(db *sqlx.DB) models.SignalingRepository {
	return &data.SignalingRepositoryMySqlImpl{DB: db}
}

func newUserFeedbackRepository(db *sqlx.DB) models.UserFeedbackRepository {
	return &data.UserFeedbackRepositoryMySqlImpl{DB: db}
}

func newGameFetchUseCase(packRepo models.PackRepository, gameRepo models.GameRepository) *usecases.GameFetchUseCase {
	return &usecases.GameFetchUseCase{
		PackRepo: packRepo,
		GameRepo: gameRepo,
	}
}

func newGameCtrlUseCase(
	serviceConf *services.ServiceConfig,
	packRepo models.PackRepository,
	gameRepo models.GameRepository,
	playerRepo models.PlayerRepository,
	msgService services.MessageService,
	orakkiDrv services.OrakkiDriver) *usecases.GameCtrlUseCase {
	return &usecases.GameCtrlUseCase{
		ServiceConfig:  serviceConf,
		PackRepo:       packRepo,
		GameRepo:       gameRepo,
		PlayerRepo:     playerRepo,
		MessageService: msgService,
		OrakkiDriver:   orakkiDrv,
	}
}

func newGameController(
	gameFetchUseCase *usecases.GameFetchUseCase,
	gameCtrlUseCase *usecases.GameCtrlUseCase) *ctrls.GameController {
	return &ctrls.GameController{
		GameFetchUseCase: gameFetchUseCase,
		GameCtrlUseCase:  gameCtrlUseCase,
	}
}

func newGameHandler(gameCtrlUseCase *usecases.GameCtrlUseCase) *handlers.GameHandler {
	return &handlers.GameHandler{GameCtrlUseCase: gameCtrlUseCase}
}

func newSignalingUseCases(
	serviceConf *services.ServiceConfig,
	msgService services.MessageService,
	gameRepo models.GameRepository,
	signalingRepo models.SignalingRepository) *usecases.SignalingUseCase {
	return &usecases.SignalingUseCase{
		ServiceConfig:  serviceConf,
		MessageService: msgService,
		GameRepo:       gameRepo,
		SignalingRepo:  signalingRepo,
	}
}

func newSignalingController(
	gameCtrlUseCase *usecases.GameCtrlUseCase,
	signalingUseCase *usecases.SignalingUseCase) *ctrls.SignalingController {
	return &ctrls.SignalingController{
		GameCtrlUseCase:  gameCtrlUseCase,
		SignalingUseCase: signalingUseCase,
	}
}

func newSignalingHandler(signalingUseCase *usecases.SignalingUseCase) *handlers.SignalingHandler {
	return &handlers.SignalingHandler{
		SignalingUseCase: signalingUseCase,
	}
}

func newPlayerUseCase(playerRepo models.PlayerRepository) *usecases.PlayerUseCase {
	return &usecases.PlayerUseCase{
		PlayerRepo: playerRepo,
	}
}

func newPlayerController(playerUseCase *usecases.PlayerUseCase) *ctrls.PlayerController {
	return &ctrls.PlayerController{
		PlayerUseCase: playerUseCase,
	}
}

func newUserFeedbackUseCase(feedbackRepo models.UserFeedbackRepository) *usecases.UserFeedbackUseCase {
	return &usecases.UserFeedbackUseCase{
		FeedbackRepo: feedbackRepo,
	}
}

func newUserFeedbackController(feedbackUseCase *usecases.UserFeedbackUseCase) *ctrls.UserFeedbackController {
	return &ctrls.UserFeedbackController{
		UserFeedbackUseCase: feedbackUseCase,
	}
}
