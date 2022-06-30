package repositories

import (
	"github.com/MrBolas/SupervisorAPI/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetTaskById(id uuid.UUID) (models.Task, error)
	CreateTask(t models.Task) (models.Task, error)
	ListTasks(filters map[string]interface{}) ([]models.Task, error)
	DeleteTask(id uuid.UUID) error
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTasksRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r TaskRepository) GetTaskById(id uuid.UUID) (models.Task, error) {
	var task models.Task

	if err := r.db.First(&task, id).Error; err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (r TaskRepository) CreateTask(t models.Task) (models.Task, error) {

	if err := r.db.Create(&t).Error; err != nil {
		return models.Task{}, err
	}

	return t, nil
}

func (r TaskRepository) ListTasks(filters map[string]interface{}) ([]models.Task, error) {
	return []models.Task{}, nil
}

func (r TaskRepository) DeleteTask(id uuid.UUID) error {

	tx := r.db.Delete(&models.Task{}, id)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected <= 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
