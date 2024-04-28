package usecase

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prakaypetch-yuw/go-clean-arch/config"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/model"
)

var _ TokenUsecase = tokenUsecaseImpl{}

var (
	accessTokenExpire  = time.Minute * 15
	refreshTokenExpire = time.Hour * 24 * 7
)

type TokenUsecase interface {
	NewToken(id uint, email string) (token *model.Token, err error)
	NewAccessToken(claims model.UserClaims) (string, error)
	NewRefreshToken(claims jwt.RegisteredClaims) (string, error)
	ParseAccessToken(accessToken string) (*model.UserClaims, error)
	ParseRefreshToken(refreshToken string) *jwt.RegisteredClaims
}

type tokenUsecaseImpl struct {
	jwtKey []byte
}

func ProvideTokenUsecase(config config.Config) TokenUsecase {
	return &tokenUsecaseImpl{
		[]byte(config.JWT.Secret),
	}
}

func (t tokenUsecaseImpl) NewAccessToken(claims model.UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString(t.jwtKey)
}

func (t tokenUsecaseImpl) NewRefreshToken(claims jwt.RegisteredClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString(t.jwtKey)
}

func (t tokenUsecaseImpl) ParseAccessToken(accessToken string) (*model.UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return t.jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	return parsedAccessToken.Claims.(*model.UserClaims), nil
}

func (t tokenUsecaseImpl) ParseRefreshToken(refreshToken string) *jwt.RegisteredClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return t.jwtKey, nil
	})

	return parsedRefreshToken.Claims.(*jwt.RegisteredClaims)
}

func (t tokenUsecaseImpl) NewToken(id uint, email string) (token *model.Token, err error) {
	idStr := strconv.Itoa(int(id))
	accessToken, err := t.NewAccessToken(model.UserClaims{
		Id:    strconv.Itoa(int(id)),
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    idStr,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpire)),
		},
	})
	if err != nil {
		return
	}
	refreshToken, err := t.NewRefreshToken(jwt.RegisteredClaims{
		Issuer:    idStr,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpire)),
	})
	if err != nil {
		return
	}

	token = &model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return
}
