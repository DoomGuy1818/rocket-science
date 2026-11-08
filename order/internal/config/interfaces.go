package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type HTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type MigratorConfig interface {
	Path() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
}
