package repository

import (
	"context"

	"github.com/DevisArya/learn-microservices/user-service/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(ctx context.Context, tx *gorm.DB, user *entity.User) (*uint, error)
	Update(ctx context.Context, tx *gorm.DB, user *entity.User) error
	Delete(ctx context.Context, tx *gorm.DB, userId uint) error
	FindById(ctx context.Context, tx *gorm.DB, userId uint) (*entity.User, error)
	FindByEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error)
	FindAll(ctx context.Context, tx *gorm.DB) (*[]entity.User, error)
}

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

// Save implements UserRepository
func (*UserRepositoryImpl) Save(ctx context.Context, tx *gorm.DB, user *entity.User) (*uint, error) {

	if err := tx.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return &user.Id, nil
}

// Update implements UserRepository
func (*UserRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, user *entity.User) error {

	if err := tx.WithContext(ctx).Where("id = ?", user.Id).Updates(&user).Error; err != nil {
		return err
	}

	return nil
}

// Delete implements UserRepository
func (*UserRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, userId uint) error {

	if err := tx.WithContext(ctx).Delete(&entity.User{}, userId).Error; err != nil {
		return err
	}

	return nil
}

// FindById implements UserRepository
func (*UserRepositoryImpl) FindById(ctx context.Context, tx *gorm.DB, userId uint) (*entity.User, error) {
	var user entity.User

	if err := tx.WithContext(ctx).First(&user, userId).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail implements UserRepository
func (*UserRepositoryImpl) FindByEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error) {
	var user entity.User

	if err := tx.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

// FindAll implements UserRepository
func (*UserRepositoryImpl) FindAll(ctx context.Context, tx *gorm.DB) (*[]entity.User, error) {

	var users []entity.User

	if err := tx.WithContext(ctx).Where("role = ?", "user").Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

// GetMyOrder implements UserRepository
func (*UserRepositoryImpl) GetMyOrder(ctx context.Context, tx *gorm.DB, userId uint, offset int, pageSize int) (*entity.User, error) {

	var user entity.User

	if err := tx.WithContext(ctx).
		Preload("Transaction", func(db *gorm.DB) *gorm.DB {
			return db.Order("transaction_time DESC").Limit(pageSize).Offset(offset)
		}).
		Preload("Transaction.TransactionDetail").
		Where("id = ?", userId).
		First(&user).
		Error; err != nil {
		return nil, err
	}

	return &user, nil
}
