package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Env      string
}

func MustLoad() (db *Config, err error) {
	err = godotenv.Load()

	if err != nil {
		fmt.Println("Error is occurred  on .env file please check", err)
		return nil, err
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")
	env := os.Getenv("ENV")
	log.Printf("Succses load env %s", pass)

	return &Config{
		Host:     host,
		Port:     port,
		Database: dbname,
		Username: pass,
		Password: pass,
		Env:      env,
	}, nil
}
