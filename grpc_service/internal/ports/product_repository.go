package ports

import (
	"context"
	"grpc_service/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) (*domain.Product, error)
	Get(ctx context.Context, id int64) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]domain.Product, error)
}
