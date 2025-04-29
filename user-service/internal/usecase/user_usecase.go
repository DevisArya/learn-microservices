package usecase

import (
	"context"
	"errors"

	"github.com/DevisArya/learn-microservices/user-service/internal/dto"
	"github.com/DevisArya/learn-microservices/user-service/internal/entity"
	"github.com/DevisArya/learn-microservices/user-service/internal/helper"
	"github.com/DevisArya/learn-microservices/user-service/internal/repository"
	"github.com/DevisArya/learn-microservices/user-service/internal/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserService interface {
	Create(ctx context.Context, request *dto.UserCreateRequest, role entity.Role) (*uint, error)
	UpdatePassword(ctx context.Context, request *dto.UserupdatePasswordRequest, id uint) error
	UpdateEmail(ctx context.Context, request *dto.UserupdateEmailRequest, id uint) error
	UpdateProfile(ctx context.Context, request *dto.UserUpdateProfileRequest, id uint) error
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*entity.User, error)
	FindAll(ctx context.Context) (*[]entity.User, error)
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *gorm.DB
	validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, DB *gorm.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		validate:       validate,
	}
}

// Create implements UserService
func (service *UserServiceImpl) Create(ctx context.Context, request *dto.UserCreateRequest, role entity.Role) (*uint, error) {
	if err := service.validate.Struct(request); err != nil {
		return nil, err
	}

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	//check used email
	cekEmail, err := service.UserRepository.FindByEmail(ctx, tx, request.Email)
	if err != nil {
		return nil, err
	}

	if !cekEmail {
		return nil, errors.New("email already use")
	}

	//hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	userData := entity.User{
		Email:       request.Email,
		Name:        request.Name,
		Password:    hashedPassword,
		PhoneNumber: request.PhoneNumbner,
		Role:        role,
	}

	id, err := service.UserRepository.Save(ctx, tx, &userData)

	if err != nil {
		return nil, err
	}
	return id, nil
}

// UpdatePassword implements UserService
func (service *UserServiceImpl) UpdatePassword(ctx context.Context, request *dto.UserupdatePasswordRequest, id uint) error {
	if err := service.validate.Struct(request); err != nil {
		return err
	}

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return err
	}

	// hash new password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}

	//validate new password same or not
	if user.Password == hashedPassword {
		return errors.New("new password must be different from the current password")
	}

	userData := entity.User{
		Id:       id,
		Password: hashedPassword,
	}

	if err := service.UserRepository.Update(ctx, tx, &userData); err != nil {
		return err
	}

	return nil
}

// UpdateEmail implements UserService
func (service *UserServiceImpl) UpdateEmail(ctx context.Context, request *dto.UserupdateEmailRequest, id uint) error {
	if err := service.validate.Struct(request); err != nil {
		return err
	}

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	//check used email
	_, err := service.UserRepository.FindByEmail(ctx, tx, request.Email)
	if err != nil {
		return err
	}

	userData := entity.User{
		Id:    id,
		Email: request.Email,
	}

	if err := service.UserRepository.Update(ctx, tx, &userData); err != nil {
		return err
	}

	return nil
}

// UpdateProfile implements UserService
func (service *UserServiceImpl) UpdateProfile(ctx context.Context, request *dto.UserUpdateProfileRequest, id uint) error {
	if err := service.validate.Struct(request); err != nil {
		return err
	}

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	userData := entity.User{
		Id:          id,
		Name:        request.Name,
		PhoneNumber: request.PhoneNumbner,
	}

	if err := service.UserRepository.Update(ctx, tx, &userData); err != nil {
		return err
	}

	return nil
}

// Delete implements UserService
func (service *UserServiceImpl) Delete(ctx context.Context, id uint) error {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if _, err := service.UserRepository.FindById(ctx, tx, id); err != nil {
		return err
	}

	if err := service.UserRepository.Delete(ctx, tx, id); err != nil {
		return err
	}

	return nil
}

// FindById implements UserService
func (service *UserServiceImpl) FindById(ctx context.Context, id uint) (*entity.User, error) {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindAll implements UserService
func (service *UserServiceImpl) FindAll(ctx context.Context) (*[]entity.User, error) {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	users, err := service.UserRepository.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
