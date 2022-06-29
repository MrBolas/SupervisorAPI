package main

import (
	"log"
	"os"

	"github.com/MrBolas/SupervisorAPI/api"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const ENV_MYSQL_USERNAME = "MYSQL_USERNAME"
const ENV_MYSQL_PASSWORD = "MYSQL_PASSWORD"
const ENV_MYSQL_HOST = "MYSQL_HOSTNAME"
const ENV_MYSQL_PORT = "MYSQL_PORT"
const ENV_MYSQL_DB = "MYSQL_DATABASE"

func buildDSN() string {
	mysqlUsername := os.Getenv(ENV_MYSQL_USERNAME)
	if mysqlUsername == "" {
		panic("missing env var: " + ENV_MYSQL_USERNAME)
	}
	mysqlPassword := os.Getenv(ENV_MYSQL_PASSWORD)
	if mysqlPassword == "" {
		panic("missing env var: " + ENV_MYSQL_PASSWORD)
	}
	mysqlHost := os.Getenv(ENV_MYSQL_HOST)
	if mysqlHost == "" {
		panic("missing env var: " + ENV_MYSQL_HOST)
	}
	mysqlPort := os.Getenv(ENV_MYSQL_PORT)
	if mysqlPort == "" {
		panic("missing env var: " + ENV_MYSQL_PORT)
	}
	mysqlDB := os.Getenv(ENV_MYSQL_DB)
	if mysqlDB == "" {
		panic("missing env var: " + ENV_MYSQL_DB)
	}

	return mysqlUsername + ":" + mysqlPassword + "@(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDB + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	db, err := gorm.Open(mysql.Open(buildDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("unale to connect to database: %v", err)
	}

	// Start API
	a := api.New(db)
	err = a.Start()
	if err != nil {
		log.Fatalf("unable to start echo: %v", err)
	}

}
