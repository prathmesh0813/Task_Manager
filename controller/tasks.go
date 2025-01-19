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
