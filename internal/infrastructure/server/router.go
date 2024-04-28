package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/handler"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/infrastructure/middleware"
)

func registerRouter(app *fiber.App, handlerProvider *handler.HandlerProvider, middlewareProvider *middleware.MiddlewareProvider) {
	// Auth
	auth := app.Group("/auth")
	auth.Post("/register", func(c *fiber.Ctx) error {
		return handlerProvider.UserHandler.Register(c)
	})
	auth.Post("/login", func(c *fiber.Ctx) error {
		return handlerProvider.UserHandler.Login(c)
	})

	// API
	apis := app.Group("/api").Use(middlewareProvider.JWTAuthMiddleware.Auth())
	apis.Get("/user", func(c *fiber.Ctx) error {
		return handlerProvider.UserHandler.UserInfo(c)
	})
}
