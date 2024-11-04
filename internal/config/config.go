package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	Database   string
	DBUsername string
	DBPassword string
	Env        string
	HttpPort   int
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
	httpPort, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))
	log.Printf("Succses load env %s", pass)

	return &Config{
		DBHost:     host,
		DBPort:     port,
		Database:   dbname,
		DBUsername: pass,
		DBPassword: pass,
		Env:        env,
		HttpPort:   httpPort,
	}, nil
}
