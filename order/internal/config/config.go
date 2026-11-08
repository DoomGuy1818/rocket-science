package config

import (
	"os"

	"github.com/DoomGuy1818/rocket-science/order/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	Postgres      PostgresConfig
	OrderHTTP     HTTPConfig
	InventoryGRPC InventoryGRPCConfig
	Migrator      MigratorConfig
	PaymentGRPC   PaymentGRPCConfig
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

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	OrderHTTP, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	InventoryGRPC, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	migrator, err := env.NewMigrationConfig()
	if err != nil {
		return err
	}

	PaymentGRPC, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		Postgres:      postgresCfg,
		OrderHTTP:     OrderHTTP,
		InventoryGRPC: InventoryGRPC,
		Migrator:      migrator,
		PaymentGRPC:   PaymentGRPC,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
