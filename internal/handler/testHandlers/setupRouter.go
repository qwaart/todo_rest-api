package test_Handler

import (
	"testing"

	"rest_api/internal/handler"
	"rest_api/internal/db/sqlite"

	"github.com/gin-gonic/gin"
)

func SetupRouter(t *testing.T) *gin.Engine {
	storage, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	r := gin.Default()
	taskHandler := handler.NewTaskHandler(storage)
	r.POST("/task", taskHandler.CreateTask)
	r.GET("/task/:id", taskHandler.GetTaskByID)
	r.DELETE("/task/:id", taskHandler.DeleteTaskByID)
	r.PATCH("/task/:id/completed", taskHandler.UpdateTaskCompletedByID)
	r.PATCH("/task/:id/uncompleted", taskHandler.UpdateTaskUncompletedByID)
	return r
}