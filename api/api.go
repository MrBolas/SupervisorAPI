package api

import (
	"github.com/MrBolas/SupervisorAPI/handlers"
	"github.com/MrBolas/SupervisorAPI/models"
	"github.com/MrBolas/SupervisorAPI/repositories"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Api struct {
	echo *echo.Echo
}

func New(db *gorm.DB) *Api {

	err := db.AutoMigrate(models.Task{})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		panic(err)
	}

	e := echo.New()

	// repositories
	tasksRepo := repositories.NewTasksRepository(db)

	// handlers
	tasksHandler := handlers.NewTasksHandler(tasksRepo)

	// auth

	g := e.Group("/v1")

	g.POST("/tasks", tasksHandler.CreateTask)
	g.GET("/tasks", tasksHandler.GetTaskList)
	g.GET("/tasks/:id", tasksHandler.GetTaskById)
	g.DELETE("/tasks/:id", tasksHandler.DeleteTask)

	return &Api{
		echo: e,
	}
}

func (api *Api) Start() error {
	return api.echo.Start(":8080")
}
