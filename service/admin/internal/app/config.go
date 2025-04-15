package app

import (
	config "github.com/nikitaSstepanov/tools/configurator"
	"REDACTED/team-11/backend/admin/internal/controller"
	"REDACTED/team-11/backend/admin/internal/usecase"
	"REDACTED/team-11/backend/admin/internal/usecase/storage"
)

type Config struct {
	Controller controller.Config `yaml:"controller"`
	Storage    storage.Config    `yaml:"storage"`
	UseCase    usecase.Config    `yaml:"usecase"`
}

func getConfig() (*Config, error) {
	var cfg Config

	if err := config.Get(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
