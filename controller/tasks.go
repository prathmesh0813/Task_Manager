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
		utils.SetResponse(c, nil, "could not parse request", true, http.StatusBadRequest)
		return
	}

	userId := c.GetInt64("userId")
	task.UserID = userId
	utils.Logger.Info("Recieved task creation request", zap.Int64("userId", userId))

	//save task in db
	err = dao.SaveTask(&task)
	if err != nil {
		utils.Logger.Error("Failed to save task", zap.Error(err))
		utils.SetResponse(c, nil, "failed to create the task", true, http.StatusBadRequest)
		return
	}

	utils.Logger.Info("Task created successfully", zap.Int64("taskId", task.ID), zap.Int64("userId", userId))
	utils.SetResponse(c, gin.H{"taskId": task.ID}, "task created successfully", false, http.StatusCreated)
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
		utils.SetResponse(c, nil, "user id not found", true, http.StatusUnauthorized)
		return
	}

	taskId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Error("Failed to parse task ID", zap.String("param", c.Param("id")), zap.Error(err))
		utils.SetResponse(c, nil, "could not parse data", true, http.StatusBadRequest)
		return
	}

	utils.Logger.Info("Fetching task", zap.Int64("taskId", taskId), zap.Int64("userId", userId.(int64)))

	task, err := dao.GetTaskByID(taskId, userId.(int64))
	if err != nil {
		utils.Logger.Error("Failed to get task or access denied", zap.Int64("taskId", taskId), zap.Int64("userId", userId.(int64)))
		utils.SetResponse(c, nil, "could not fetch task", true, http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("Task fetched successfully", zap.Int64("taskId", taskId), zap.Int64("userId", userId.(int64)))
	utils.SetResponse(c, task, "task fetched successfully", false, http.StatusOK)
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
		utils.SetResponse(c, nil, "not authorized", true, http.StatusUnauthorized)
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
		utils.SetResponse(c, nil, "could not fetch tasks", true, http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := (totalTasks + int64(limit) - 1) / int64(limit)

	// Respond with tasks and pagination metadata
	utils.Logger.Info("Tasks fetched successfully", zap.Int64("userId", userId.(int64)), zap.String("sortOrder", sortOrder), zap.String("completed", completed), zap.Int("page", page), zap.Int("limit", limit), zap.Int("totalPages", int(totalPages)))
	utils.SetResponse(c, gin.H{"tasks": tasks, "totalPages": totalPages, "currentPage": page}, "task fetched successfully", false, http.StatusOK)
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
		utils.SetResponse(c, nil, "could not parse data", true, http.StatusBadRequest)
		return
	}

	userID := c.GetInt64("userId")

	task, err := dao.GetTaskByID(taskId, userID)
	if err != nil {
		utils.Logger.Error("Failed to fetch task", zap.Int64("taskId", taskId), zap.Int64("userId", userID), zap.Error(err))
		utils.SetResponse(c, nil, "could not fetch task", true, http.StatusInternalServerError)
		return
	}

	if task.UserID != userID {
		utils.Logger.Warn("User not authorized to update task", zap.Int64("taskId", taskId), zap.Int64("userId", userID))
		utils.SetResponse(c, nil, "not authorized to update", true, http.StatusUnauthorized)
		return
	}

	var updatedTask models.Task

	err = c.ShouldBindJSON(&updatedTask)
	if err != nil {
		utils.Logger.Error("Failed to bind json", zap.Error(err))
		utils.SetResponse(c, nil, "could not parse request", true, http.StatusBadRequest)
		return
	}

	updatedTask.ID = taskId

	err = dao.Update(&updatedTask)
	if err != nil {
		utils.Logger.Error("Failed to update task", zap.Int64("taskId", taskId), zap.Error(err))
		utils.SetResponse(c, nil, "could not update task", true, http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("Task updated successfully", zap.Int64("taskId", taskId), zap.Int64("userId", userID))
	utils.SetResponse(c, nil, "task updated successfully", false, http.StatusOK)
}

// delete task
func DeleteTask(c *gin.Context) {
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	taskId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Error("Failed to parse task ID", zap.String("param", c.Param("id")), zap.Error(err))
		utils.SetResponse(c, nil, "could not parse data", true, http.StatusBadRequest)
		return
	}

	userID := c.GetInt64("userId")

	task, err := dao.GetTaskByID(taskId, userID)
	if err != nil {
		utils.Logger.Error("Failed to fetch task for deletion", zap.Int64("taskId", taskId), zap.Int64("userId", userID), zap.Error(err))
		utils.SetResponse(c, nil, "could not fetch task", true, http.StatusInternalServerError)
		return
	}

	if task.UserID != userID {
		utils.Logger.Warn("User not authorized to delete task", zap.Int64("taskId", taskId), zap.Int64("userId", userID))
		utils.SetResponse(c, nil, "not authorized to delete", true, http.StatusUnauthorized)
		return
	}

	err = dao.Delete(task)
	if err != nil {
		utils.Logger.Error("Failed to delete task", zap.Int64("taskId", taskId), zap.Error(err))
		utils.SetResponse(c, nil, "could not delete task", true, http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("Task deleted successfully", zap.Int64("taskId", taskId), zap.Int64("userId", userID))
	utils.SetResponse(c, nil, "task deleted successfully", false, http.StatusOK)
}
