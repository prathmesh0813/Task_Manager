package dao

import (
	"task_manager/models"
	"time"
)

// save task in db
func SaveTask(t *models.Task) error {
	result := DB.Create(&t)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// fetch rask by id
func GetTaskByID(id, userId int64) (*models.Task, error) {
	var task models.Task
	result := DB.Where("id = ? AND user_id = ?", id, userId).First(&task)
	if result.Error != nil {
		return &task, result.Error
	}

	return &task, nil
}

// fetch all tasks using filters
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
		return nil, 0, result.Error
	}

	return tasks, totalTasks, nil
}

// update task in db
func Update(t *models.Task) error {
	result := DB.Model(&Task{}).Where("id = ?", t.ID).Updates(Task{Completed: t.Completed, UpdatedAt: time.Now()})
	if result.Error != nil {
		return result.Error
	}

	return result.Error
}

// delete task in db
func Delete(t *models.Task) error {
	result := DB.Delete(t)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
