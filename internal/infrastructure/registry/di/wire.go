//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/prakaypetch-yuw/go-clean-arch/config"
)

func InitializeApplication(configPath config.FilePath) (*Application, func(), error) {
	wire.Build(ApplicationSet)
	return &Application{}, func() {}, nil
}
