package controller

import (
	"net/http"
	"task_manager/dao"
	"task_manager/middlewares"
	"task_manager/models"
	"task_manager/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//create task for user
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
