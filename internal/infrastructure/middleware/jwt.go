package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/usecase"
	"github.com/rs/zerolog/log"
)

var _ JWTAuthMiddleware = &jwtAuthMiddlewareImpl{}

type JWTAuthMiddleware interface {
	Auth() fiber.Handler
}

type jwtAuthMiddlewareImpl struct {
	tokenUsecase usecase.TokenUsecase
}

func ProvideJWTAuthMiddleware(tokenUsecase usecase.TokenUsecase) JWTAuthMiddleware {
	return &jwtAuthMiddlewareImpl{
		tokenUsecase: tokenUsecase,
	}
}

func (j *jwtAuthMiddlewareImpl) Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		const BearerSchema = "Bearer "
		authHeader := c.GetReqHeaders()["Authorization"][0]
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}
		tokenString := authHeader[len(BearerSchema):]
		claims, err := j.tokenUsecase.ParseAccessToken(tokenString)
		if err != nil {
			log.Error().Ctx(c.Context()).Err(err)
			return fiber.ErrUnauthorized
		}
		c.Locals("claims", claims)
		return c.Next()
	}
}
