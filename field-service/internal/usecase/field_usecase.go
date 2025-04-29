package usecase

import (
	"context"

	"github.com/DevisArya/learn-microservices/field-service/internal/dto"
	"github.com/DevisArya/learn-microservices/field-service/internal/entity"
	"github.com/DevisArya/learn-microservices/field-service/internal/helper"
	"github.com/DevisArya/learn-microservices/field-service/internal/repository"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type FieldUseCase interface {
	Save(ctx context.Context, request *dto.FieldRequest) (*entity.Field, error)
	Update(ctx context.Context, request *dto.FieldRequest, id uint) error
	Delete(ctx context.Context, fieldId uint) error
	FindById(ctx context.Context, fieldId uint) (*entity.Field, error)
	FindAll(ctx context.Context, limit uint32, page uint32) (*[]entity.Field, *dto.PaginationResponse, error)
}

type FieldUseCaseImpl struct {
	FieldRepository repository.FieldRepository
	DB              *gorm.DB
	validate        *validator.Validate
}

func NewFieldUseCase(FieldRepository repository.FieldRepository, DB *gorm.DB, validate *validator.Validate) FieldUseCase {
	return &FieldUseCaseImpl{
		FieldRepository,
		DB,
		validate,
	}
}

// Save implements FieldUseCase
func (service *FieldUseCaseImpl) Save(ctx context.Context, request *dto.FieldRequest) (*entity.Field, error) {

	if err := service.validate.Struct(request); err != nil {
		return nil, err
	}

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	fieldData := entity.Field{
		Name:  request.Name,
		Type:  request.Type,
		Price: uint32(request.Price),
	}

	response, err := service.FieldRepository.Save(ctx, tx, &fieldData)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Update implements FieldUseCase
func (service *FieldUseCaseImpl) Update(ctx context.Context, request *dto.FieldRequest, id uint) error {

	if err := service.validate.Struct(request); err != nil {
		return err
	}

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if _, err := service.FieldRepository.FindById(ctx, tx, id); err != nil {
		return err
	}

	fieldData := entity.Field{
		Id:          id,
		Name:        request.Name,
		Type:        request.Type,
		Description: request.Description,
		Price:       request.Price,
	}

	if err := service.FieldRepository.Update(ctx, tx, &fieldData); err != nil {
		return err
	}

	return nil
}

// Delete implements FieldUseCase
func (service *FieldUseCaseImpl) Delete(ctx context.Context, fieldId uint) error {

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if _, err := service.FieldRepository.FindById(ctx, tx, fieldId); err != nil {
		return err
	}

	if err := service.FieldRepository.Delete(ctx, tx, fieldId); err != nil {
		return err
	}

	return nil
}

// FindById implements FieldUseCase
func (service *FieldUseCaseImpl) FindById(ctx context.Context, fieldId uint) (*entity.Field, error) {

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	field, err := service.FieldRepository.FindById(ctx, tx, fieldId)
	if err != nil {
		return nil, err
	}

	return field, nil
}

// FindAll implements FieldUseCase
func (service *FieldUseCaseImpl) FindAll(ctx context.Context, limit uint32, page uint32) (*[]entity.Field, *dto.PaginationResponse, error) {

	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	fields, totalRecord, err := service.FieldRepository.FindAll(ctx, tx, limit, offset)

	if err != nil {
		return nil, nil, err
	}

	totalPage := totalRecord / int64(limit)

	return fields, &dto.PaginationResponse{
		CurrentPage: page,
		Limit:       limit,
		TotalRecord: uint32(totalRecord),
		TotalPage:   uint32(totalPage),
	}, nil
}
