package controller

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"task_manager/dao"
	"task_manager/logger"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// create task for user
func CreateTask(c *gin.Context) {
	requestID := requestid.Get(c)

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	requestBody := string(bodyBytes)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	var task models.Task
	err = c.ShouldBindJSON(&task)
	if err != nil {
		logger.Error(requestID, "failed to parse request", err.Error(), requestBody)
		utils.SetResponse(c, requestID, nil, "cannot parsed the requested data", true, http.StatusBadRequest)
		return
	}

	userId := c.GetInt64("userId")
	task.UserID = userId

	// If Completed is not set, set it to false
	// if task.Completed == nil {
	// 	falseVal := false
	// 	task.Completed = &falseVal
	// }

	logger.Info(requestID, "Recieved task creation request", "userID: "+strconv.Itoa(int(userId)), requestBody)

	//save task in db
	err = dao.SaveTask(&task)
	if err != nil {
		logger.Error(requestID, "failed to save task", err.Error(), "userID: "+strconv.Itoa(int(userId)), requestBody)
		utils.SetResponse(c, requestID, nil, "failed to create the task", true, http.StatusBadRequest)
		return
	}

	logger.Info(requestID, "task created successfully", "taskID: "+strconv.Itoa(int(task.ID)), "userID: "+strconv.Itoa(int(userId)), requestBody)
	utils.SetResponse(c, requestID, gin.H{"taskId": task.ID}, "task created successfully", false, http.StatusCreated)
}

// fetch task by id
func GetTask(c *gin.Context) {
	requestID := requestid.Get(c)

	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		logger.Error(requestID, "User Id not found in context", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "user id not found", true, http.StatusUnauthorized)
		return
	}

	taskId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(requestID, "failed to parse task Id", c.Param("id"), err.Error())
		utils.SetResponse(c, requestID, nil, "could not parse request data", true, http.StatusBadRequest)
		return
	}

	logger.Info(requestID, "Fetching Task", "userID: "+strconv.Itoa(int(userId.(int64))), "taskID: "+strconv.Itoa(int(taskId)))

	task, err := dao.GetTaskByID(taskId, userId.(int64))
	if err != nil {
		logger.Error(requestID, "Failed to fetch task or access denied", "userID: "+strconv.Itoa(int(userId.(int64))), "taskID: "+strconv.Itoa(int(taskId)), err.Error())
		utils.SetResponse(c, requestID, nil, "could not fetch task or access denied", true, http.StatusBadRequest)
		return
	}

	logger.Info(requestID, "Task fetched successfully", "userID: "+strconv.Itoa(int(userId.(int64))), "taskID: "+strconv.Itoa(int(taskId)))
	utils.SetResponse(c, requestID, task, "task fetched successfully", false, http.StatusOK)
}

// Fetch all tasks using query params also
func GetTasksByQuery(c *gin.Context) {
	requestID := requestid.Get(c)

	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		logger.Warn(requestID, "Unauthorized", "userID: "+strconv.Itoa(int(userId.(int64))))
		utils.SetResponse(c, requestID, nil, "not authorized", true, http.StatusUnauthorized)
		return
	}

	// Retrieve query parameters
	sortOrder := c.DefaultQuery("sort", "asc")
	completed := c.Query("completed")

	// Pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		logger.Warn(requestID, "Invalid page parameter", c.DefaultQuery("page", "1"), err.Error())
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit < 1 {
		logger.Warn(requestID, "Invalid limit parameter", c.DefaultQuery("limit", "5"), err.Error())
		limit = 5
	}

	offset := (page - 1) * limit

	// Fetch tasks with filters, sorting, and pagination
	tasks, totalTasks, err := dao.GetTasksWithFilters(userId.(int64), sortOrder, completed, limit, offset)
	if err != nil {
		logger.Error(requestID, "failed to fetch tasks", "userID: "+strconv.Itoa(int(userId.(int64))), err.Error())
		utils.SetResponse(c, requestID, nil, "could not fetch tasks", true, http.StatusBadRequest)
		return
	}

	// Calculate total pages
	totalPages := (totalTasks + int64(limit) - 1) / int64(limit)

	// Respond with tasks and pagination metadata
	logger.Info(requestID, "task fetched successfully", "userID: "+strconv.Itoa(int(userId.(int64))), "sortOrder: "+sortOrder, "completed: "+completed, "page: "+strconv.Itoa(int(page)), "limit: "+strconv.Itoa(int(limit)), "totalPages: "+strconv.Itoa(int(totalPages)))
	utils.SetResponse(c, requestID, gin.H{"tasks": tasks, "totalPages": totalPages, "currentPage": page}, "task fetched successfully", true, http.StatusOK)
}

// Update Task
func UpdateTask(c *gin.Context) {
	requestID := requestid.Get(c)

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	requestBody := string(bodyBytes)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	taskId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(requestID, "failed to parse task id", c.Param("id"), err.Error(), requestBody)
		utils.SetResponse(c, requestID, nil, "could not parse task id", true, http.StatusBadRequest)
		return
	}

	userID := c.GetInt64("userId")

	task, err := dao.GetTaskByID(taskId, userID)
	if err != nil {
		logger.Error(requestID, "failed to fetch task", "userID: "+strconv.Itoa(int(userID)), "taskID: "+strconv.Itoa(int(taskId)), err.Error(), requestBody)
		utils.SetResponse(c, requestID, nil, "could not fetch task", true, http.StatusBadRequest)
		return
	}

	if task.UserID != userID {
		logger.Warn(requestID, "User not authorized to update task", "userID: "+strconv.Itoa(int(userID)), "taskID: "+strconv.Itoa(int(taskId)), requestBody)
		utils.SetResponse(c, requestID, nil, "not authorized to update", true, http.StatusUnauthorized)
		return
	}

	if task.Completed == "false" {
		task.Completed = "true"
	}

	err = dao.Update(task)
	if err != nil {
		logger.Error(requestID, "failed to update task", "taskID: "+strconv.Itoa(int(taskId)), err.Error(), requestBody)
		utils.SetResponse(c, requestID, nil, "could not update task", true, http.StatusBadRequest)
		return
	}

	logger.Info(requestID, "task updated successfully", "userID: "+strconv.Itoa(int(userID)), "taskID: "+strconv.Itoa(int(taskId)), requestBody)
	utils.SetResponse(c, requestID, nil, "task updated successfully", false, http.StatusOK)
}

// delete task
func DeleteTask(c *gin.Context) {
	requestID := requestid.Get(c)
	err := middlewares.CheckTokenPresent(c)
	if err != nil {
		return
	}

	taskId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(requestID, "failed to parse task id", c.Param("id"), err.Error())
		utils.SetResponse(c, requestID, nil, "could not parse data", true, http.StatusBadRequest)
		return
	}

	userID := c.GetInt64("userId")

	task, err := dao.GetTaskByID(taskId, userID)
	if err != nil {
		logger.Error(requestID, "failed to fetch task for deletion", "userID: "+strconv.Itoa(int(userID)), "taskID: "+strconv.Itoa(int(taskId)), err.Error())
		utils.SetResponse(c, requestID, nil, "could not fetch task", true, http.StatusBadRequest)
		return
	}

	if task.UserID != userID {
		logger.Warn(requestID, "user not authorized to delete task", "userID: "+strconv.Itoa(int(userID)), "taskID: "+strconv.Itoa(int(taskId)))
		utils.SetResponse(c, requestID, nil, "user not authorized to delete task", true, http.StatusUnauthorized)
		return
	}

	err = dao.Delete(task)
	if err != nil {
		logger.Error(requestID, "failed to delete task", "taskID: "+strconv.Itoa(int(taskId)), err.Error())
		utils.SetResponse(c, requestID, nil, "could not delete task", true, http.StatusBadRequest)
		return
	}

	logger.Info(requestID, "Task deleted successfully", "userID: "+strconv.Itoa(int(userID)), "taskID: "+strconv.Itoa(int(taskId)))
	utils.SetResponse(c, requestID, nil, "task deleted successfully", false, http.StatusOK)
}
