package boot

import (
	"daos_core/internal/boot/di/adapter"
	"daos_core/internal/boot/di/config"
	"daos_core/internal/boot/di/controller"
	"daos_core/internal/boot/di/database"
	"daos_core/internal/boot/di/repository"
	"daos_core/internal/boot/di/route"
	"daos_core/internal/boot/di/service"
	"daos_core/internal/boot/di/utils"
	"daos_core/internal/transport/middleware"
	l "daos_core/internal/utils/logger"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
)

type App struct {
	Engine       *gin.Engine
	DB           *database.Container
	Routes       *route.Container
	Controllers  *controller.Container
	Config       *config.Container
	Services     *service.Container
	Adapters     *adapter.Container
	Repositories *repository.Container
	Utils        *utils.Container
	// middleware
}

func InitApp() (*App, error) {
	// cfg
	path, exist := os.LookupEnv("DEV_ENV")
	if !exist {
		return nil, errors.New("DEV_ENV return empty path")
	}

	cfg, err := config.OpenCfg(path)
	if err != nil {
		l.Logg.Error(err.Error())
		return nil, err
	}

	l.Logg.Info("App: cfg was init")

	// db
	storages, err := database.NewDatabases(
		cfg.Common.Postgres,
		cfg.Common.Redis,
	)
	if err != nil {
		l.Logg.Error(err.Error())
		return nil, err
	}

	l.Logg.Info("App: database was init")

	// router
	g := gin.Default()

	//utils layer
	utils, err := utils.RegisterAll(&cfg.Common.Service.Auth)
	if err != nil {
		return nil, err
	}

	l.Logg.Info("App: utils was init")

	// external API
	adapters, err := adapter.RegisterAll(cfg.Amo, cfg.Telegram)
	if err != nil {
		return nil, err
	}

	l.Logg.Info("App: adapters was init")

	// database/cache
	repositories, err := repository.RegisterAll(storages.PostgresDB, storages.RedisDB)
	if err != nil {
		return nil, err
	}

	l.Logg.Info("App: repositories was init")

	// bussines logic
	services, err := service.RegisterAll(repositories, adapters, cfg.Common.Service, utils)
	if err != nil {
		return nil, err
	}

	l.Logg.Info("App: services was init")

	// api handlers
	controllers, err := controller.RegisterAll(services)
	if err != nil {
		return nil, err
	}

	l.Logg.Info("App: controllers was init")

	// middleware
	mid := middleware.NewMiddleware(services.Auth)
	l.Logg.Info("App: middleware was init")

	// routes
	routes, err := route.RegisterAll(g, controllers, &mid)
	if err != nil {
		return nil, err
	}

	l.Logg.Info("App: routes was init")

	return &App{
		DB:           storages,
		Engine:       g,
		Routes:       routes,
		Controllers:  controllers,
		Repositories: repositories,
		Services:     services,
		Adapters:     adapters,
		Config:       cfg,
		Utils:        utils,
	}, nil
}

// graceful shutdown
func (app *App) DisposeApp() {
	app.DB.PostgresDB.Disconnect()
	app.DB.RedisDB.Disconnect()
}
