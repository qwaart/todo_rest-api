package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"fmt"
	"log/slog"
	
	"rest_api/internal/db/sqlite"

	"github.com/gin-gonic/gin"
)

// The TaskHandler is responsible for processing HTTP requests related to tasks (creating, receiving, updating, deleting).
type TaskHandler struct {
	storage *sqlite.Storage
	log 	*slog.Logger
}

func NewTaskHandler(s *sqlite.Storage, log *slog.Logger) *TaskHandler {
	return &TaskHandler{storage: s, log: log}
}

type NewTask struct {
	Title 	string 	`json:"title" binding:"required"`
}

//get id
func GetID(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	} 
	return id, nil
}
// gets the title from json and creates a task in the database. Sends back the task(json)
// POST (/task)
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req NewTask

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid request", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	id, err := h.storage.AddTask(req.Title)
	if err != nil {
		h.log.Error("Failed to create task", slog.String("title", req.Title), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task: " + err.Error()})
	 	return
	}

	h.log.Info("Task created successfully", slog.Int64("id", id), slog.String("title", req.Title))
	c.JSON(http.StatusCreated, gin.H{
		"id": id,
		"title": req.Title,
		"completed": false,
	})
}

// returns the corresponding task
// GET (/task/:id)
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id, err := GetID(c)
	if err != nil {
		h.log.Warn("Invalid task ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	h.log.Debug("Fetching task", slog.Int64("id", id))
	title, completed, err := h.storage.GetTaskByID(id)
	if err != nil {
	if errors.Is(err, sql.ErrNoRows) {
		h.log.Warn("Task not found", slog.Int64("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	h.log.Error("Failed to get task", slog.Int64("id", id), slog.Any("error", err))
	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get task"})
	return
	}

	h.log.Info("Task retrieved successfully", slog.Int64("id", id), slog.String("title", title))
	c.JSON(http.StatusOK, gin.H{
		"id": id,
		"title": title,
		"completed": completed,
	})
}

// delete corresponding task in the db
// DELETE (/task/:id)
func (h *TaskHandler) DeleteTaskByID(c *gin.Context) {
	id, err := GetID(c)
	if err != nil {
		h.log.Warn("Invalid task ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	h.log.Debug("Deleting task", slog.Int64("id", id))
	err = h.storage.DeleteTaskByID(id)
	if err != nil {
		h.log.Error("Failed to delete task", slog.Int64("id", id), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

// changes the “completed” field to the required ending (completed/uncompleted)
// UPDATE/PUT
func (h *TaskHandler) CompletedTask(c *gin.Context) {
	id, err := GetID(c)
	if err != nil {
		h.log.Warn("Invalid task ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.log.Debug("Marking task as completed", slog.Int64("id", id))
	err = h.storage.MarkTaskTrue(id)
	if err != nil {
		h.log.Error("Failed to mark task completed", slog.Int64("id", id), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	h.log.Info("Task marked as completed", slog.Int64("id", id))
	c.JSON(http.StatusOK, gin.H{"message": "task updated"})
}

func (h *TaskHandler) UncompletedTask(c *gin.Context) {
	id, err := GetID(c)
	if err != nil {
		h.log.Warn("Invalid task ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.log.Debug("Marking task as uncompleted", slog.Int64("id", id))
	err = h.storage.MarkTaskFalse(id)
	if err != nil {
		h.log.Error("Failed to mark task uncompleted", slog.Int64("id", id), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	h.log.Info("Task marked as uncompleted", slog.Int64("id", id))
	c.JSON(http.StatusOK, gin.H{"message": "task updated"})
}