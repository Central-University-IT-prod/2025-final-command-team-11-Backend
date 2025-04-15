package app

import (
	"github.com/nikitaSstepanov/coffee-id/internal/controller"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase"
	config "github.com/nikitaSstepanov/tools/configurator"
)

type Config struct {
	Controller controller.Config `yaml:"controller"`
	UseCase    usecase.Config    `yaml:"usecase"`
}

func getConfig() (*Config, error) {
	var cfg Config

	if err := config.Get(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
