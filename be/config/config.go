package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type appConfig struct {
	Port                  string
	DataBase              DBConfig
	AccessTokenPrivatekey string
}

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

func Init() (appConfig, error) {
	err := loadEnvironment()
	if err != nil {
		return appConfig{}, fmt.Errorf("error loading environment: %w", err)
	}

	appConfig, err := loadConfig()
	if err != nil {
		return appConfig, fmt.Errorf("error loading config: %w", err)
	}

	return appConfig, nil
}

func loadEnvironment() error {
	if len(os.Getenv("APP_ENVIRONMENT")) == 0 || os.Getenv("APP_ENVIRONMENT") == "localhost" {
		fmt.Println("app is running in the localhost")
		err := godotenv.Load("./config/localhost.configmap.env")
		if err != nil {
			return fmt.Errorf("error loading localhost.configmap.env file: %w", err)
		}
	}

	return nil
}

func loadConfig() (appConfig, error) {

	DB := DBConfig{
		Username: os.Getenv("APP_DATABASE_USERNAME"),
		Password: os.Getenv("APP_DATABASE_PASSWORD"),
		Host:     os.Getenv("APP_DATABASE_HOST"),
		Port:     os.Getenv("APP_DATABASE_PORT"),
		DBName:   os.Getenv("APP_DATABASE_DBNAME"),
	}

	accessTokenPrivatekey := strings.ReplaceAll(os.Getenv("APP_PRIVATEKEY"), "\\n", "\n")
	return appConfig{
		Port:                  os.Getenv("APP_PORT"),
		DataBase:              DB,
		AccessTokenPrivatekey: accessTokenPrivatekey,
	}, nil
}
