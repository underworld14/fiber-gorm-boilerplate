package repository

import (
	"gorm.io/gorm"

	"fiber-gorm/internal/models"
)

type TaskRepository struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) CreateTask(task *models.Task) error {
	return r.DB.Create(task).Error
}

func (r *TaskRepository) FindAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	return tasks, r.DB.Find(&tasks).Error
}
