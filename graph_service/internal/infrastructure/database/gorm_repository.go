package database

import (
	"context"
	"errors"
	"graph_service/internal/domain"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) (*gormRepository, error) {
	// Automigración: GORM crea o actualiza la tabla 'products' basándose en la entidad de dominio
	err := db.AutoMigrate(&domain.Product{})
	if err != nil {
		return nil, err
	}
	return &gormRepository{db: db}, nil
}
func (r *gormRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	result := r.db.WithContext(ctx).Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *gormRepository) Get(ctx context.Context, id int64) (*domain.Product, error) {
	var product domain.Product
	result := r.db.WithContext(ctx).First(&product, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("product_not_found")
		}
		return nil, result.Error
	}
	return &product, nil
}

func (r *gormRepository) Update(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	// Actualiza todos los campos del producto en la base de datos
	result := r.db.WithContext(ctx).Save(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *gormRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&domain.Product{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("product_not_found")
	}
	return nil
}

func (r *gormRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	result := r.db.WithContext(ctx).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
