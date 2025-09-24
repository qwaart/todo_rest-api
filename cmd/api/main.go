package main

import (
	"os"
	"log/slog"
	_ "net/http"

	"github.com/gin-gonic/gin"

	sl "rest_api/internal/lib/logger/slog"
	"rest_api/internal/config"
	"rest_api/internal/db/sqlite"
	"rest_api/internal/handler"
	"rest_api/internal/auth"
)

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
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
	    log.Error("failed to init storage", sl.Err(err))
	    os.Exit(1)
	}

	// init auth storage(work with api_keys, permission)
	authStorage := auth.NewStorage(storage.DB())
	if err := authStorage.Init(); err != nil {
		log.Error("failed to init auth tables", sl.Err(err))
		os.Exit(1)
	}
	if err := authStorage.EnsureAdminSetup(log); err != nil {
		log.Error("failed to ensure admin key", sl.Err(err))
		os.Exit(1)
	}

	// create permission
	permissionsToCreate := []string{"task.create", "task.delete", "task.update"}
	for _, perm := range permissionsToCreate {
    if err := authStorage.CreatePermission(perm); err != nil {
        log.Warn("failed to create permission", slog.String("permission", perm), sl.Err(err))
    } else {
        log.Info("created permission", slog.String("permission", perm))
    }
}
	// init services & handlers
	taskHandler := handler.NewTaskHandler(storage)
	authService := auth.NewService(authStorage)
	authHandler := handler.NewAuthorization(authService, log)


	// init router: gin
	r := gin.Default()

	//Public routes
	r.GET("/task/:id", taskHandler.GetTaskByID)
	r.POST("/register", authHandler.Register)

	// Admin routes
	admin := r.Group("/admin",
		auth.AuthMiddleware(authStorage, log),
		auth.RequirePermission(authStorage, "admin", log),
	)
	{
		admin.GET("/keys", authHandler.ListKeys)
		admin.GET("/permissions", authHandler.ListPermissions)
		admin.POST("/permissions", authHandler.CreatePermission)
		admin.POST("/keys/:id/permissions", authHandler.GrantPermission)
	}

	// Protected routes
	authorized := r.Group("/", auth.AuthMiddleware(authStorage, log))
{
    authorized.POST("/task", auth.RequirePermission(authStorage, "task.create", log), taskHandler.CreateTask)
    authorized.DELETE("/task/:id", auth.RequirePermission(authStorage, "task.delete", log), taskHandler.DeleteTaskByID)
    authorized.PATCH("/task/:id/completed", auth.RequirePermission(authStorage, "task.update", log), taskHandler.CompletedTask)
    authorized.PATCH("/task/:id/uncompleted", auth.RequirePermission(authStorage, "task.update", log), taskHandler.UncompletedTask)
}
	if err := r.Run(":8080"); err != nil {
		log.Error("Failed to run server", sl.Err(err))
		os.Exit(1)
	}
}

//  setupLog configures the logger depending on the environment (local/dev/prod).
func setupLog(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envlocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
 	case envProd:
 		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
 	default:
 		log = slog.New (slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
 	}
 	return log
}