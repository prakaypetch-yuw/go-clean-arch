package repository

import (
	"context"
	"errors"

	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/customerror"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/entity"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/infrastructure/db"
	"gorm.io/gorm"
)

var _ UserRepository = &userRepositoryImpl{}

type UserRepository interface {
	Store(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userRepositoryImpl struct {
	defaultDB *gorm.DB
}

func (u *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	err = u.tx(ctx).Where("email = ?", email).First(&user).Error
	return
}

func (u *userRepositoryImpl) tx(ctx context.Context) *gorm.DB {
	tx := db.GormFromContext(ctx)
	if tx == nil {
		return u.defaultDB.WithContext(ctx)
	}
	return tx.WithContext(ctx)
}

func (u userRepositoryImpl) Store(ctx context.Context, ent *entity.User) error {
	var result *gorm.DB

	if ent == nil {
		return customerror.ErrEntityIsNull
	}

	if ent.ID == 0 {
		result = u.tx(ctx).WithContext(ctx).Create(ent)
		if result.Error == nil && result.RowsAffected == 0 {
			return customerror.ErrNoRowsAffected
		}
		return result.Error
	}

	result = u.tx(ctx).WithContext(ctx).Where("id = ?", ent.ID).Take(&entity.User{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			resultIfNotFound := u.tx(ctx).WithContext(ctx).Create(&ent)
			return resultIfNotFound.Error
		}
		return result.Error
	}

	result = u.tx(ctx).WithContext(ctx).Select("*").Updates(ent)
	if result.Error == nil && result.RowsAffected == 0 {
		return customerror.ErrNoRowsAffected
	}
	return result.Error
}

func ProvideUserRepository(defaultDB *gorm.DB) UserRepository {
	return &userRepositoryImpl{defaultDB: defaultDB}
}
