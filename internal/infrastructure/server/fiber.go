package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/handler"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/infrastructure/middleware"
)

func ProvideFiberServer(handlerProvider *handler.HandlerProvider, middlewareProvider *middleware.MiddlewareProvider) *fiber.App {
	app := fiber.New()
	registerRouter(app, handlerProvider, middlewareProvider)
	return app
}
