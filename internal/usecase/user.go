package usecase

import (
	"context"

	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/entity"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/model"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

var _ UserUsecase = &userUsecaseImpl{}

type UserUsecase interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userUsecaseImpl struct {
	userRepository repository.UserRepository
}

func (u userUsecaseImpl) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.userRepository.FindByEmail(ctx, email)
}

func (u userUsecaseImpl) Register(ctx context.Context, req *model.RegisterRequest) (*entity.User, error) {
	password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	}

	err := u.userRepository.Store(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func ProvideUserUsecase(userRepository repository.UserRepository) UserUsecase {
	return &userUsecaseImpl{
		userRepository,
	}
}
