package app

import (
	"context"
	"graph_service/internal/domain"
	"graph_service/internal/ports"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetProduct(ctx context.Context, id int64) (*domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	ListProducts(ctx context.Context) ([]domain.Product, error)
}

type productService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	// Aquí podrías agregar validaciones de negocio (ej. precio > 0)
	return s.repo.Create(ctx, product)
}

func (s *productService) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	return s.repo.Get(ctx, id)
}

func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	return s.repo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context) ([]domain.Product, error) {
	return s.repo.GetAll(ctx)
}
