package main

import (
	"context"
	"daos_core/internal/boot"
	log "daos_core/internal/utils/logger"
	"daos_core/internal/utils/worker"
	"net/http"
	"os"
)

func main() {
	log.InitLogger()
	defer log.StopAsyncLogger(context.Background())

	// семафор
	worker.InitGlobalPool()
	defer worker.Pool.Stop()

	app, err := boot.InitApp()
	if err != nil {
		log.Logg.Error(err.Error())
		os.Exit(-1)
	}

	if app == nil {
		log.Logg.Error(err.Error())
		os.Exit(-1)
	}

	defer app.DisposeApp()
	log.Logg.Info("app was init")

	srv := http.Server{
		Addr:    "localhost:3000",
		Handler: app.Engine,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Logg.Error(err.Error())
		os.Exit(-1)
	}
}
