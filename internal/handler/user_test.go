package handler

import (
	"errors"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/entity"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/model"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/usecase/mock"
	"github.com/prakaypetch-yuw/go-clean-arch/tool/util"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandlerTestSuite struct {
	suite.Suite
	ctrl             *gomock.Controller
	underTest        UserHandler
	mockUserUseCase  *mock.MockUserUseCase
	mockTokenUseCase *mock.MockTokenUseCase
	fiber            *fiber.App
}

func (s *UserHandlerTestSuite) SetupSuite() {
	s.fiber = fiber.New()
	s.ctrl = gomock.NewController(s.T())
	s.mockUserUseCase = mock.NewMockUserUseCase(s.ctrl)
	s.mockTokenUseCase = mock.NewMockTokenUseCase(s.ctrl)
	s.underTest = ProvideUserHandler(
		s.mockUserUseCase,
		s.mockTokenUseCase,
	)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (s *UserHandlerTestSuite) TestLogin() {
	id := gofakeit.Uint()
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, 10)
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	accessToken := gofakeit.UUID()
	refreshToken := gofakeit.UUID()

	s.Run("Login: success", func() {
		body := &model.Credentials{
			Username: email,
			Password: password,
		}
		jsonBody, _ := json.Marshal(body)
		fctx, cleanup, err := util.GetFiberCtx(s.fiber, "localhost", "/auth/login", fasthttp.MethodPost, jsonBody)
		defer cleanup()

		user := &entity.User{
			Model: gorm.Model{
				ID: id,
			},
			Name:     name,
			Email:    email,
			Password: encryptedPassword,
		}

		token := &model.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		s.mockUserUseCase.EXPECT().GetUserByEmail(fctx.Context(), email).Return(user, nil)
		s.mockTokenUseCase.EXPECT().NewToken(user.ID, user.Email).Return(token, nil)
		err = s.underTest.Login(fctx)
		s.Assert().NoError(err)
		var res model.Token
		err = json.Unmarshal(fctx.Response().Body(), res)
		if err != nil {
			return
		}
		s.Assert().Equal(accessToken, res.AccessToken)
		s.Assert().Equal(refreshToken, res.RefreshToken)
	})

	s.Run("Login: user not found", func() {
		invalidEmail := gofakeit.Email()
		body := &model.Credentials{
			Username: invalidEmail,
			Password: password,
		}
		jsonBody, _ := json.Marshal(body)
		fctx, cleanup, err := util.GetFiberCtx(s.fiber, "localhost", "/auth/login", fasthttp.MethodPost, jsonBody)
		defer cleanup()

		s.mockUserUseCase.EXPECT().GetUserByEmail(fctx.Context(), invalidEmail).Return(nil, gorm.ErrRecordNotFound)
		err = s.underTest.Login(fctx)
		s.Assert().NoError(err)
		var res map[string]string
		err = json.Unmarshal(fctx.Response().Body(), res)
		if err != nil {
			return
		}
		s.Assert().Equal("user not found", res["message"])
	})

	s.Run("Login: incorrect password", func() {
		invalidPassword := gofakeit.Password(true, true, true, true, true, 10)
		body := &model.Credentials{
			Username: email,
			Password: invalidPassword,
		}
		jsonBody, _ := json.Marshal(body)
		fctx, cleanup, err := util.GetFiberCtx(s.fiber, "localhost", "/auth/login", fasthttp.MethodPost, jsonBody)
		defer cleanup()

		user := &entity.User{
			Model: gorm.Model{
				ID: id,
			},
			Name:     name,
			Email:    email,
			Password: encryptedPassword,
		}

		s.mockUserUseCase.EXPECT().GetUserByEmail(fctx.Context(), email).Return(user, nil)
		err = s.underTest.Login(fctx)
		s.Assert().NoError(err)
		var res map[string]string
		err = json.Unmarshal(fctx.Response().Body(), res)
		if err != nil {
			return
		}
		s.Assert().Equal("incorrect password", res["message"])
	})

	s.Run("Login: new token failed", func() {
		body := &model.Credentials{
			Username: email,
			Password: password,
		}
		jsonBody, _ := json.Marshal(body)
		fctx, cleanup, err := util.GetFiberCtx(s.fiber, "localhost", "/auth/login", fasthttp.MethodPost, jsonBody)
		defer cleanup()

		user := &entity.User{
			Model: gorm.Model{
				ID: id,
			},
			Name:     name,
			Email:    email,
			Password: encryptedPassword,
		}

		newTokenErr := errors.New("new token failed")

		s.mockUserUseCase.EXPECT().GetUserByEmail(fctx.Context(), email).Return(user, nil)
		s.mockTokenUseCase.EXPECT().NewToken(user.ID, user.Email).Return(nil, newTokenErr)
		err = s.underTest.Login(fctx)
		s.Assert().NoError(err)
		var res map[string]string
		err = json.Unmarshal(fctx.Response().Body(), res)
		if err != nil {
			return
		}
		s.Assert().Equal("could not login", res["message"])
	})
}

func (s *UserHandlerTestSuite) TestUserInfo() {
	id := gofakeit.Uint()
	email := gofakeit.Email()
	name := gofakeit.Name()

	s.Run("GetUserInfo", func() {
		fctx, cleanup, err := util.GetFiberCtx(s.fiber, "localhost", "/api/user", fasthttp.MethodGet, nil)
		defer cleanup()
		userClaims := &model.UserClaims{
			Id:               strconv.Itoa(int(id)),
			Email:            email,
			RegisteredClaims: jwt.RegisteredClaims{},
		}
		fctx.Locals("claims", userClaims)
		user := &entity.User{
			Model: gorm.Model{
				ID: id,
			},
			Name:  name,
			Email: email,
		}
		s.mockUserUseCase.EXPECT().GetUserByEmail(fctx.Context(), email).Return(user, nil)

		err = s.underTest.UserInfo(fctx)
		s.Assert().NoError(err)
		var res entity.User
		err = json.Unmarshal(fctx.Response().Body(), res)
		if err != nil {
			return
		}
		s.Assert().Equal(name, res.Name)
		s.Assert().Equal(email, res.Email)
		s.Assert().Equal(id, res.ID)
	})
	s.Run("GetUserInfo: fail to get claims", func() {
		fctx, cleanup, err := util.GetFiberCtx(s.fiber, "localhost", "/api/user", fasthttp.MethodGet, nil)
		defer cleanup()

		err = s.underTest.UserInfo(fctx)
		s.Assert().NoError(err)
		var res map[string]string
		err = json.Unmarshal(fctx.Response().Body(), res)
		if err != nil {
			return
		}
		s.Assert().Equal("server error", res["message"])
	})
	s.Run("GetUserInfo: user not found", func() {
		fctx, cleanup, err := util.GetFiberCtx(s.fiber, "localhost", "/api/user", fasthttp.MethodGet, nil)
		defer cleanup()
		userClaims := &model.UserClaims{
			Id:               strconv.Itoa(int(id)),
			Email:            email,
			RegisteredClaims: jwt.RegisteredClaims{},
		}
		fctx.Locals("claims", userClaims)
		s.mockUserUseCase.EXPECT().GetUserByEmail(fctx.Context(), email).Return(nil, gorm.ErrRecordNotFound)

		err = s.underTest.UserInfo(fctx)
		s.Assert().NoError(err)
		var res map[string]string
		err = json.Unmarshal(fctx.Response().Body(), res)
		if err != nil {
			return
		}
		s.Assert().Equal("server error", res["message"])
	})
}

func (s *UserHandlerTestSuite) TestRegister() {
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, 10)
	s.Run("Register: success", func() {
		registerRequest := &model.RegisterRequest{
			Name:     name,
			Email:    email,
			Password: password,
		}
		jsonBody, _ := json.Marshal(registerRequest)
		fctx, cleanup, err := util.GetFiberCtx(s.fiber, "localhost", "/api/user", fasthttp.MethodGet, jsonBody)
		defer cleanup()

		s.mockUserUseCase.EXPECT().Register(fctx.Context(), registerRequest).Return(&entity.User{}, nil)
		err = s.underTest.Register(fctx)
		s.Assert().NoError(err)
	})
	s.Run("Register: fail", func() {
		registerRequest := &model.RegisterRequest{
			Name:     name,
			Email:    email,
			Password: password,
		}
		jsonBody, _ := json.Marshal(registerRequest)
		fctx, cleanup, err := util.GetFiberCtx(s.fiber, "localhost", "/api/user", fasthttp.MethodGet, jsonBody)
		defer cleanup()

		s.mockUserUseCase.EXPECT().Register(fctx.Context(), registerRequest).Return(nil, errors.New("register failed"))
		err = s.underTest.Register(fctx)
		s.Assert().NoError(err)
		var res map[string]string
		err = json.Unmarshal(fctx.Response().Body(), res)
		if err != nil {
			return
		}
		s.Assert().Equal("could not register", res["message"])
	})
}
