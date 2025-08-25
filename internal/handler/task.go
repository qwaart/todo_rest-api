package handler

import (
	"net/http"
	"strconv"
	
	"rest_api/internal/db/sqlite"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	storage *sqlite.Storage
}

func NewTaskHandler(s *sqlite.Storage) *TaskHandler {
	return &TaskHandler{storage: s}
}

type NewTask struct {
	Title 	string 	`json:"title" binding:"required"`
	Private bool	`json:"private"`
}

// POST (/task)
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req NewTask

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	id, err := h.storage.AddTask(req.Title, req.Private)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task: " + err.Error()})
	 	return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
		"title": req.Title,
		"completed": false,
		"private": req.Private,
	})
}

// GET (/task/:id)
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	idInt := c.Param("id")
	id, err := strconv.ParseInt(idInt, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	title, completed, private, err := h.storage.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// if task is private -> need token
	if private {
		token := c.GetHeader("Authorization")
		if token != "Bearer supersecrettoken" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
		"title": title,
		"completed": completed,
		"private": private,
	})
}

// DELETE (/task/:id)
func (h *TaskHandler) DeleteTaskByID(c *gin.Context) {
	idInt := c.Param("id")
	id, err := strconv.ParseInt(idInt, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	//get task
	_, _, private, err := h.storage.GetTaskFullByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if private {
		token := c.GetHeader("Authorization")
		if token != "Bearer supersecrettoken" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}

	err = h.storage.DeleteTaskByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

// UPDATE/PUT
func (h *TaskHandler) UpdateTaskCompletedByID(c *gin.Context) {
	idInt := c.Param("id")
	id, err := strconv.ParseInt(idInt, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "inavlid id"})
		return
	}

	//get task
	_, _, private, err := h.storage.GetTaskFullByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "unauthorized"})
		return
	}

	if private {
		token := c.GetHeader("Authorization")
		if token != "Bearer supersecrettoken" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}

	err = h.storage.MarkTaskTrue(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated"})
}

func (h *TaskHandler) UpdateTaskUncompletedByID(c *gin.Context) {
	idInt := c.Param("id")
	id, err := strconv.ParseInt(idInt, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// get task
	_, _, private, err := h.storage.GetTaskFullByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "unauthorized"})
		return
	}
	if private {
		token := c.GetHeader("Authorization")
		if token != "Bearer supersecrettoken" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}

	err = h.storage.MarkTaskFalse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated"})
}