package repository

import (
	"context"

	"github.com/DevisArya/learn-microservices/field-service/internal/entity"
	"gorm.io/gorm"
)

type FieldRepository interface {
	Save(ctx context.Context, tx *gorm.DB, field *entity.Field) (*entity.Field, error)
	Update(ctx context.Context, tx *gorm.DB, field *entity.Field) error
	Delete(ctx context.Context, tx *gorm.DB, fieldId uint) error
	FindById(ctx context.Context, tx *gorm.DB, fieldId uint) (*entity.Field, error)
	FindAll(ctx context.Context, tx *gorm.DB, limit uint32, offset uint32) (*[]entity.Field, int64, error)
}

type FieldRepositoryImpl struct{}

func NewFieldRepository() FieldRepository {
	return &FieldRepositoryImpl{}
}

// Save implements FieldRepository
func (repository *FieldRepositoryImpl) Save(ctx context.Context, tx *gorm.DB, field *entity.Field) (*entity.Field, error) {

	if err := tx.WithContext(ctx).Create(field).Error; err != nil {
		return nil, err
	}
	return field, nil
}

// Update implements FieldRepository
func (repository *FieldRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, field *entity.Field) error {

	if err := tx.WithContext(ctx).Updates(&field).Error; err != nil {
		return err
	}

	return nil
}

// Delete implements FieldRepository
func (repository *FieldRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, fieldId uint) error {

	if err := tx.WithContext(ctx).Delete(&entity.Field{}, fieldId).Error; err != nil {
		return err
	}

	return nil
}

// FindById implements FieldRepository
func (repository *FieldRepositoryImpl) FindById(ctx context.Context, tx *gorm.DB, fieldId uint) (*entity.Field, error) {
	var field entity.Field

	if err := tx.WithContext(ctx).First(&field, fieldId).Error; err != nil {
		return nil, err
	}
	return &field, nil
}

// FindAll implements FieldRepository
func (repository *FieldRepositoryImpl) FindAll(ctx context.Context, tx *gorm.DB, limit uint32, offset uint32) (*[]entity.Field, int64, error) {
	var fields []entity.Field
	var count int64

	if err := tx.WithContext(ctx).Model(&entity.Field{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := tx.WithContext(ctx).Limit(int(limit)).Offset(int(offset)).Order("id ASC").Find(&fields).Error; err != nil {
		return nil, 0, err
	}

	return &fields, count, nil
}
