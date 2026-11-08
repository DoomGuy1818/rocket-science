package env

import "github.com/caarlos0/env/v11"

type migrationEnvConfig struct {
	migrationDir string `env:"MIGRATION_DIRECTORY, required"`
}

type migrationConfig struct {
	raw *migrationEnvConfig
}

func NewMigrationConfig() (*migrationConfig, error) {
	var raw migrationEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &migrationConfig{
		raw: &raw,
	}, nil
}

func (m *migrationConfig) Path() string {
	return m.raw.migrationDir
}
