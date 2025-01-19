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

//fetch all tasks using filters
func GetTasksWithFilters(userId int64, sortOrder, completed string, limit, offset int) ([]Task, int64, error) {
	var tasks []Task
	var totalTasks int64

	// Start building the query
	query := DB.Order("created_at "+sortOrder).Where("user_id = ?", userId)

	// Apply the Completed filter if provided
	if completed != "" {
		query = query.Where("completed = ?", completed)
	}

	// Count the total number of tasks (without limit/offset)
	query.Model(&Task{}).Count(&totalTasks)

	// Apply pagination
	query = query.Limit(limit).Offset(offset)

	// Execute the query
	result := query.Find(&tasks)
	if result.Error != nil {
		utils.Logger.Error("Failed to fetch tasks", zap.Error(result.Error), zap.Int64("userId", userId), zap.String("sortOrder", sortOrder), zap.String("completed", completed), zap.Int("limit", limit), zap.Int("offset", offset))
		return nil, 0, result.Error
	}

	utils.Logger.Info("Tasks fetched successfully", zap.Int64("userId", userId), zap.String("sortOrder", sortOrder), zap.String("completed", completed), zap.Int("tasksCount", len(tasks)), zap.Int("totalTasks", int(totalTasks)))
	return tasks, totalTasks, nil
}

