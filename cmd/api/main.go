package main

import (
	"context"

	"github.com/kgundo/gundo-go/loggers"
	"github.com/kgundo/gundo-go/routers"
	"github.com/kgundo/gundo-go/validators"
	"github.com/spf13/viper"
)

var logger *loggers.StandardLogger

func main() {
	ctx := context.Background()
	logger = loggers.Logger()
	loadEnvVariables()
	validators.RegisterValidators()
	r := routers.Setup(ctx)
	_ = r.Run(":" + viper.Get("APP_PORT").(string))
}

func loadEnvVariables() {
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logger.CommonException(("Failed to read environment file"))
	}
}
