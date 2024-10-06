package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvsConfig struct {
	DBProvider string
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     int
}

var envs EnvsConfig

func Setup() {
	envs = initEnvsConfigs()
}

func Envs() *EnvsConfig {
	return &envs
}

// TODO: Make default values
func initEnvsConfigs() EnvsConfig {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Failed to load envs: ", err)
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))

	return EnvsConfig{
		DBProvider: os.Getenv("DB_PROVIDER"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     port,
	}
}

func (e *EnvsConfig) DBUrl() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s", e.DBProvider, e.DBUser, e.DBPassword, e.DBHost, e.DBPort, e.DBName)
}

func (e *EnvsConfig) DBAddress() string {
	return fmt.Sprintf("%s:%d", e.DBHost, e.DBPort)
}
