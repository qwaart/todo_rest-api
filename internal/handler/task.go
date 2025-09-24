package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"fmt"
	
	"rest_api/internal/db/sqlite"

	"github.com/gin-gonic/gin"
)

// The TaskHandler is responsible for processing HTTP requests related to tasks (creating, receiving, updating, deleting).
type TaskHandler struct {
	storage *sqlite.Storage
}

func NewTaskHandler(s *sqlite.Storage) *TaskHandler {
	return &TaskHandler{storage: s}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	id, err := h.storage.AddTask(req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task: " + err.Error()})
	 	return
	}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	title, completed, err := h.storage.GetTaskByID(id)
	if err != nil {
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get task"})
	return
	}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err = h.storage.DeleteTaskByID(id)
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.storage.MarkTaskTrue(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated"})
}

func (h *TaskHandler) UncompletedTask(c *gin.Context) {
	id, err := GetID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.storage.MarkTaskFalse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated"})
}