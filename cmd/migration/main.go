package main

import (
	"github.com/kgundo/gundo-go/cmd/migration/tables"
	"github.com/kgundo/gundo-go/db"
	"github.com/kgundo/gundo-go/loggers"
	"github.com/spf13/viper"
)

var logger *loggers.StandardLogger

func main() {
	logger = loggers.Logger()
	loadEnvVariables()

	dbInstance := db.InitDB()
	tables.MigrateUser(dbInstance)
}

func loadEnvVariables() {
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logger.CommonException(("Failed to read environment file"))
	}
}
