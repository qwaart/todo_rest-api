package main

import (
	"os"
	"log/slog"
	_ "net/http"

	"github.com/gin-gonic/gin"

	"rest_api/internal/config"
	"rest_api/internal/lib/logger/slog"
	"rest_api/internal/db/sqlite"
	"rest_api/internal/handler"
)

var storage *sqlite.Storage

const (
	envlocal 	= "local"
	envDev 		= "dev"
	envProd		= "prod"
)

func main() {
	//initialization config
	cfg := config.MustLoad()

	// initialization logger
	log := setupLog(cfg.Env)

	log.Info("Start api", slog.String("env", cfg.Env))
	log.Debug("Debug message are enabled")

	//initialization storage(sqlite)
	var _ = sqlite.New

	var err error
	storage, err = sqlite.New(cfg.StoragePath)
	if err != nil {
	    log.Error("failed to init storage", sl.Err(err))
	    os.Exit(1)
	}


	taskHandler := handler.NewTaskHandler(storage)

	// init router: gin
	r := gin.Default()
	r.POST("/task", taskHandler.CreateTask)
	r.GET("/task/:id", taskHandler.GetTaskByID)
	r.DELETE("/task/:id", taskHandler.DeleteTaskByID)
	r.PATCH("/task/:id/completed", taskHandler.UpdateTaskCompletedByID)
	r.PATCH("/task/:id/uncompleted", taskHandler.UpdateTaskUncompletedByID)
	r.Run(":8080")
}

// For logs
func setupLog(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envlocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
 	case envProd:
 		log = slog.New(
 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
 		)
 	default:
 		log = slog.New (
 			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
 		)
 	}
 	return log
}