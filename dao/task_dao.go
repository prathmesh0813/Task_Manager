package dao

import (
	"task_manager/models"
	"task_manager/utils"

	"go.uber.org/zap"
)

// save task in db
func SaveTask(t *models.Task) error {
	result := DB.Create(&t)
	if result.Error != nil {
		utils.Logger.Error("Failed to save task", zap.Error(result.Error), zap.Int64("userId", t.ID))
		return result.Error
	}

	utils.Logger.Info("Task saved successfully", zap.Int64("taskId", t.ID))
	return nil
}

//fetch rask by id
func GetTaskByID(id, userId int64) (*Task, error) {
	var task Task
	result := DB.Where("id = ? AND user_id = ?", id, userId).First(&task)
	if result.Error != nil {
		utils.Logger.Error("Failed to fetch task by id", zap.Error(result.Error), zap.Int64("taskId", id))
		return &Task{}, result.Error
	}

	utils.Logger.Info("Task fetched by id successfully", zap.Int64("taskId", id), zap.Int64("userId", userId))
	return &task, nil
}
