package database

import (
	"grpc_service/internal/domain"
	"grpc_service/internal/ports"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) (ports.ProductRepository, error) {
	// Automigración: GORM crea o actualiza la tabla 'products' basándose en la entidad de dominio
	err := db.AutoMigrate(&domain.Product{})
	if err != nil {
		return nil, err
	}
	return &gormRepository{db: db}, nil
}
