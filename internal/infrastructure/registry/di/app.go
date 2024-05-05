package di

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/prakaypetch-yuw/go-clean-arch/config"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/repository"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/handler"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/infrastructure/db"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/infrastructure/middleware"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/infrastructure/server"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/usecase"
)

type Application struct {
	Cfg    config.Config
	Server *fiber.App
}

var MainBindingSet = wire.NewSet(config.ProvideConfig,
	ServerSet,
	DatabaseSet,
	UseCaseSet,
	MiddlewareSet,
	ProviderSet,
	RepositorySet,
	HandlerSet,
)

var ServerSet = wire.NewSet(
	server.ProvideFiberServer,
)

var DatabaseSet = wire.NewSet(
	db.ProvideDB,
)

var UseCaseSet = wire.NewSet(
	usecase.ProvideUserUseCase,
	usecase.ProvideTokenUseCase,
)

var MiddlewareSet = wire.NewSet(
	middleware.ProvideJWTAuthMiddleware,
)

var ProviderSet = wire.NewSet(
	middleware.ProvideMiddlewareProvider,
	handler.ProvideHandlerProvider,
)

var RepositorySet = wire.NewSet(
	repository.ProvideUserRepository,
)

var HandlerSet = wire.NewSet(
	handler.ProvideUserHandler,
)

var ApplicationSet = wire.NewSet(MainBindingSet, wire.Struct(new(Application), "*"))
