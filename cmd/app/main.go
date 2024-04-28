package main

import (
	"github.com/prakaypetch-yuw/go-clean-arch/config"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/infrastructure/registry/di"
)

func main() {
	app, cleanUpFn, err := di.InitializeApplication(config.GetConfigPath())
	defer cleanUpFn()
	if err != nil {
		panic(err)
	}
	err = app.Server.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
