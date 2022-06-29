package api

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Api struct {
	echo *echo.Echo
}

func New(db *gorm.DB) *Api {

	err := db.AutoMigrate()
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

	// handlers

	// auth

	g := e.Group("/v1")

	g.POST("/tasks", echo.NotFoundHandler)
	g.GET("/tasks", echo.NotFoundHandler)
	g.GET("/tasks/:id", echo.NotFoundHandler)
	g.DELETE("/tasks/:id", echo.NotFoundHandler)

	return &Api{
		echo: e,
	}
}

func (api *Api) Start() error {
	return api.echo.Start(":8080")
}
