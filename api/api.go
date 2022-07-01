package api

import (
	"os"

	"github.com/MrBolas/SupervisorAPI/auth"
	"github.com/MrBolas/SupervisorAPI/handlers"
	"github.com/MrBolas/SupervisorAPI/models"
	"github.com/MrBolas/SupervisorAPI/repositories"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Api struct {
	echo *echo.Echo
}

const ENV_PUBLIC_KEY_URL = "AUTH0_PUBLIC_KEY_URL"
const ENV_CRYPTO_KEY = "CRYPTO_KEY"

func New(db *gorm.DB) *Api {

	err := db.AutoMigrate(models.Task{})
	if err != nil {
		panic(err)
	}

	e := echo.New()

	// encryption
	cryptKey := os.Getenv(ENV_CRYPTO_KEY)
	if cryptKey == "" {
		panic("missing env var: " + ENV_CRYPTO_KEY)
	}
	ce := handlers.NewCryptoEngine(cryptKey)

	// repositories
	tasksRepo := repositories.NewTasksRepository(db)

	// handlers
	tasksHandler := handlers.NewTasksHandler(tasksRepo, ce)

	// auth
	publicKeyUrl := os.Getenv(ENV_PUBLIC_KEY_URL)
	if publicKeyUrl == "" {
		panic("missing env var: " + ENV_PUBLIC_KEY_URL)
	}

	jwtConfig, err := auth.JWTConfig(publicKeyUrl)
	if err != nil {
		panic(err)
	}

	g := e.Group("/v1")

	// middleware
	g.Use(middleware.JWTWithConfig(jwtConfig))

	g.POST("/tasks", tasksHandler.CreateTask)
	g.GET("/tasks", tasksHandler.GetTaskList)
	g.GET("/tasks/:id", tasksHandler.GetTaskById)
	g.PUT("/tasks/:id", tasksHandler.UpdateTask)
	g.DELETE("/tasks/:id", tasksHandler.DeleteTask)

	return &Api{
		echo: e,
	}
}

func (api *Api) Start() error {
	return api.echo.Start(":8080")
}
