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
