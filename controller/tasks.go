package controller

import (
	"net/http"
	"strconv"
	"task_manager/dao"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// create task for user
func CreateTask(c *gin.Context) {

	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	var task models.Task
	err = c.ShouldBindJSON(&task)
	if err != nil {
		utils.Logger.Error("Failed to parse request", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not parse request")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	userId := c.GetInt64("userId")
	task.UserID = userId
	utils.Logger.Info("Recieved task creation request", zap.Int64("userId", userId))

	//save task in db
	err = dao.SaveTask(&task)
	if err != nil {
		utils.Logger.Error("Failed to save task", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "failed to create the task")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	utils.Logger.Info("Task created successfully", zap.Int64("taskId", task.ID), zap.Int64("userId", userId))

	c.Set("response", gin.H{"taskId": task.ID})
	c.Set("message", "task created successfully")
	c.Set("error", false)
	c.Status(http.StatusCreated)
}

// fetch task by id
func GetTask(c *gin.Context) {

	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		utils.Logger.Error("User ID not found in context", zap.String("context", "userId"))

		c.Set("response", nil)
		c.Set("message", "user id not found")
		c.Set("error", true)
		c.Status(http.StatusUnauthorized)
		return
	}

	taskId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Error("Failed to parse task ID", zap.String("param", c.Param("id")), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not parse data")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	utils.Logger.Info("Fetching task", zap.Int64("taskId", taskId), zap.Int64("userId", userId.(int64)))

	task, err := dao.GetTaskByID(taskId, userId.(int64))
	if err != nil {
		utils.Logger.Error("Failed to get task or access denied", zap.Int64("taskId", taskId), zap.Int64("userId", userId.(int64)))

		c.Set("response", nil)
		c.Set("message", "could not fetch task")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("Task fetched successfully", zap.Int64("taskId", taskId), zap.Int64("userId", userId.(int64)))

	c.Set("response", task)
	c.Set("message", "task fetched successfully")
	c.Set("error", false)
	c.Status(http.StatusOK)
}

// Fetch all tasks using query params also
func GetTasksByQuery(c *gin.Context) {

	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		utils.Logger.Warn("Unauthorized access attempt in getTasksByQuery")

		c.Set("response", nil)
		c.Set("message", "not authorized")
		c.Set("error", true)
		c.Status(http.StatusUnauthorized)
		return
	}

	// Retrieve query parameters
	sortOrder := c.DefaultQuery("sort", "asc") // Default sort order: ascending
	completed := c.Query("completed")          // Optional filter

	// Pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1")) // Default page: 1
	if err != nil || page < 1 {
		utils.Logger.Warn("Invalid page parameter", zap.String("page", c.DefaultQuery("page", "1")), zap.Error(err))
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5")) // Default limit: 5
	if err != nil || limit < 1 {
		utils.Logger.Warn("Invalid limit parameter", zap.String("limit", c.DefaultQuery("limit", "5")), zap.Error(err))
		limit = 5
	}

	offset := (page - 1) * limit

	// Fetch tasks with filters, sorting, and pagination
	tasks, totalTasks, err := dao.GetTasksWithFilters(userId.(int64), sortOrder, completed, limit, offset)
	if err != nil {
		utils.Logger.Error("Failed to get tasks", zap.Int64("userId", userId.(int64)), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not fetch tasks")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := (totalTasks + int64(limit) - 1) / int64(limit)

	// Respond with tasks and pagination metadata
	utils.Logger.Info("Tasks fetched successfully", zap.Int64("userId", userId.(int64)), zap.String("sortOrder", sortOrder), zap.String("completed", completed), zap.Int("page", page), zap.Int("limit", limit), zap.Int("totalPages", int(totalPages)))

	c.Set("response", gin.H{
		"tasks":       tasks,
		"totalPages":  totalPages,
		"currentPage": page})
	c.Set("message", "task fetched successfully")
	c.Set("error", false)
	c.Status(http.StatusOK)
}

// Update Task
func UpdateTask(c *gin.Context) {
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	taskId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Error("Failed to parse task ID", zap.String("param", c.Param("id")), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not parse data")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	userID := c.GetInt64("userId")

	task, err := dao.GetTaskByID(taskId, userID)
	if err != nil {
		utils.Logger.Error("Failed to fetch task", zap.Int64("taskId", taskId), zap.Int64("userId", userID), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not fetch task")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	if task.UserID != userID {
		utils.Logger.Warn("User not authorized to update task", zap.Int64("taskId", taskId), zap.Int64("userId", userID))

		c.Set("response", nil)
		c.Set("message", "not authorized to update")
		c.Set("error", true)
		c.Status(http.StatusUnauthorized)
		return
	}

	var updatedTask models.Task

	err = c.ShouldBindJSON(&updatedTask)
	if err != nil {
		utils.Logger.Error("Failed to bind json", zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not parse request")
		c.Set("error", true)
		c.Status(http.StatusBadRequest)
		return
	}

	updatedTask.ID = taskId

	err = dao.Update(&updatedTask)
	if err != nil {
		utils.Logger.Error("Failed to update task", zap.Int64("taskId", taskId), zap.Error(err))

		c.Set("response", nil)
		c.Set("message", "could not update task")
		c.Set("error", true)
		c.Status(http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("Task updated successfully", zap.Int64("taskId", taskId), zap.Int64("userId", userID))

	c.Set("response", nil)
	c.Set("message", "task updated successfully")
	c.Set("error", false)
	c.Status(http.StatusOK)
}
