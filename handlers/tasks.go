package handlers

import (
	"github.com/MrBolas/SupervisorAPI/repositories"
	"github.com/labstack/echo/v4"
)

type TasksHandler struct {
	repo repositories.Repository
}

func NewTasksHandler(repo repositories.Repository) *TasksHandler {
	return nil
}

func (th *TasksHandler) GetTaskById(c echo.Context) error {
	return nil
}

func (th *TasksHandler) GetTaskList(c echo.Context) error {
	return nil
}

func (th *TasksHandler) CreateTask(c echo.Context) error {
	return nil
}

func (th *TasksHandler) DeleteTask(c echo.Context) error {
	return nil
}
