package config

import "os"

type Config struct {
	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	Host string
	Port string
}

var Cfg Config

func init() {
	docker := os.Getenv("DOCKER_ENV")
	if docker == "1" {
		Cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
		Cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
		Cfg.Postgres.User = os.Getenv("POSTGRES_USER")
		Cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
		Cfg.Postgres.DBName = os.Getenv("POSTGRES_DB")
		Cfg.Port = os.Getenv("PORT")
		Cfg.Host = os.Getenv("HOST")
	} else {
		Cfg.Postgres.Host = "localhost"
		Cfg.Postgres.Port = "5432"
		Cfg.Postgres.User = "postgres"
		Cfg.Postgres.Password = "secret"
		Cfg.Postgres.DBName = "byteme"
		Cfg.Port = "8080"
		Cfg.Host = "localhost"
	}
}
