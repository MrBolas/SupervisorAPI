package repositories

import (
	"github.com/MrBolas/SupervisorAPI/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetTaskById(id uuid.UUID) (models.Task, error)
	CreateTask(t models.Task, workerId string) (models.Task, error)
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
	return models.Task{}, nil
}

func (r TaskRepository) CreateTask(t models.Task, workerId string) (models.Task, error) {
	return models.Task{}, nil
}

func (r TaskRepository) ListTasks(filters map[string]interface{}) ([]models.Task, error) {
	return []models.Task{}, nil
}

func (r TaskRepository) DeleteTask(id uuid.UUID) error {
	return nil
}
