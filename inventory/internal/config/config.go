package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	Mongo         MongoConfig
	InventoryGRPC GRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		Mongo:         mongoCfg,
		InventoryGRPC: inventoryGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
