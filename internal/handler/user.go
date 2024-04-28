package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/model"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

var _ UserHandler = &userHandlerImpl{}

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	UserInfo(c *fiber.Ctx) error
}

type userHandlerImpl struct {
	userUseCase  usecase.UserUsecase
	tokenUsecase usecase.TokenUsecase
}

func (u userHandlerImpl) UserInfo(c *fiber.Ctx) error {
	claims := c.Locals("claims")
	userClaims, ok := claims.(*model.UserClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"message": "server error",
			})
	}
	user, err := u.userUseCase.GetUserByEmail(c.Context(), userClaims.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "user not found",
			})
	}

	return c.JSON(user)
}

func (u userHandlerImpl) Login(c *fiber.Ctx) error {
	var credentials model.Credentials

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	user, err := u.userUseCase.GetUserByEmail(c.Context(), credentials.Username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "user not found",
			})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(credentials.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": "incorrect password",
			})
	}

	token, err := u.tokenUsecase.NewToken(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"message": "could not login",
			})
	}
	return c.JSON(token)
}

func (u userHandlerImpl) Register(c *fiber.Ctx) error {
	var req *model.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}

	_, err := u.userUseCase.Register(c.Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{})
}

func ProvideUserHandler(userUsecase usecase.UserUsecase, tokenUsecase usecase.TokenUsecase) UserHandler {
	return &userHandlerImpl{
		userUseCase:  userUsecase,
		tokenUsecase: tokenUsecase,
	}
}
