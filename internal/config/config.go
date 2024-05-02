package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type configApp struct {
	REDIS_URL 		string 	`env:"REDIS_URL" env-default:"localhost:1234"`
	REDIS_PASSWORD 	string 	`env:"REDIS_PASSWORD" env-default:""`
	REDIS_DATABASE 	int 	`env:"REDIS_DATABASE" env-default:"0"`
	REINDEXER_URL	string 	`env:"REINDEXER_URL" env-default:"cproto://127.0.0.1:6534/db"`
}

func GetEnvironment () (configApp) {
	var cfg configApp
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		fmt.Println("Error reading environment: ", err.Error())
		os.Exit(-1)
	}
	return cfg
}