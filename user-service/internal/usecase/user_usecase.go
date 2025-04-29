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

type UserUseCase interface {
	Create(ctx context.Context, request *dto.UserCreateRequest, role entity.Role) (*uint, error)
	UpdatePassword(ctx context.Context, request *dto.UserupdatePasswordRequest, id uint) error
	UpdateEmail(ctx context.Context, request *dto.UserupdateEmailRequest, id uint) error
	UpdateProfile(ctx context.Context, request *dto.UserUpdateProfileRequest, id uint) error
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*entity.User, error)
	FindAll(ctx context.Context, limit, page uint32) (*[]entity.User, *dto.PaginationResponse, error)
}

type UserUseCaseImpl struct {
	UserRepository repository.UserRepository
	DB             *gorm.DB
	validate       *validator.Validate
}

func NewUserUseCase(userRepository repository.UserRepository, DB *gorm.DB, validate *validator.Validate) UserUseCase {
	return &UserUseCaseImpl{
		UserRepository: userRepository,
		DB:             DB,
		validate:       validate,
	}
}

// Create implements UserUseCase
func (service *UserUseCaseImpl) Create(ctx context.Context, request *dto.UserCreateRequest, role entity.Role) (*uint, error) {
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

// UpdatePassword implements UserUseCase
func (service *UserUseCaseImpl) UpdatePassword(ctx context.Context, request *dto.UserupdatePasswordRequest, id uint) error {
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

// UpdateEmail implements UserUseCase
func (service *UserUseCaseImpl) UpdateEmail(ctx context.Context, request *dto.UserupdateEmailRequest, id uint) error {
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

// UpdateProfile implements UserUseCase
func (service *UserUseCaseImpl) UpdateProfile(ctx context.Context, request *dto.UserUpdateProfileRequest, id uint) error {
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

// Delete implements UserUseCase
func (service *UserUseCaseImpl) Delete(ctx context.Context, id uint) error {
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

// FindById implements UserUseCase
func (service *UserUseCaseImpl) FindById(ctx context.Context, id uint) (*entity.User, error) {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindAll implements UserUseCase
func (service *UserUseCaseImpl) FindAll(ctx context.Context, limit, page uint32) (*[]entity.User, *dto.PaginationResponse, error) {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	offset := (page - 1) * limit

	users, totalRecord, err := service.UserRepository.FindAll(ctx, tx, int(limit), int(offset))
	if err != nil {
		return nil, nil, err
	}

	totalPage := uint32(*totalRecord) / limit

	return users, &dto.PaginationResponse{
		CurrentPage: page,
		Limit:       limit,
		TotalRecord: uint32(*totalRecord),
		TotalPage:   totalPage,
	}, nil
}
