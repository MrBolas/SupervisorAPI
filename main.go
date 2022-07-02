package main

import (
	"log"
	"os"
	"strconv"

	"github.com/MrBolas/SupervisorAPI/api"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const ENV_MYSQL_USERNAME = "MYSQL_USERNAME"
const ENV_MYSQL_PASSWORD = "MYSQL_PASSWORD"
const ENV_MYSQL_HOST = "MYSQL_HOSTNAME"
const ENV_MYSQL_PORT = "MYSQL_PORT"
const ENV_MYSQL_DB = "MYSQL_DATABASE"

const ENV_REDIS_HOST = "REDIS_HOST"
const ENV_REDIS_PORT = "REDIS_PORT"
const ENV_REDIS_DB = "REDIS_DB"
const ENV_REDIS_PASSWORD = "REDIS_PASSWORD"

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

	// database
	db, err := gorm.Open(mysql.Open(buildDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("unale to connect to database: %v", err)
	}

	redisHost := os.Getenv(ENV_REDIS_HOST)
	if redisHost == "" {
		panic("missing env var: " + ENV_REDIS_HOST)
	}
	redisPort := os.Getenv(ENV_REDIS_PORT)
	if redisPort == "" {
		panic("missing env var: " + ENV_REDIS_PORT)
	}
	redisDB := os.Getenv(ENV_REDIS_DB)
	if redisDB == "" {
		panic("missing env var: " + ENV_REDIS_DB)
	}

	dBNumber, err := strconv.Atoi(redisDB)
	if err != nil {
		panic("invalid Redis DB number: " + redisDB)
	}

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",
		DB:       dBNumber,
	})

	// Start API
	a := api.New(db, rdb)
	err = a.Start()
	if err != nil {
		log.Fatalf("unable to start echo: %v", err)
	}

}
