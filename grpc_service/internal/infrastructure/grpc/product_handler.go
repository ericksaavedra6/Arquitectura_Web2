package grpc

import (
	"context"
	"grpc_service/internal/app"
	"grpc_service/internal/domain"
	"grpc_service/internal/infrastructure/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProductHandler implementa la interfaz autogenerada por gRPC (ProductServiceServer)
type ProductHandler struct {
	pb.UnimplementedProductServiceServer // Embeber para asegurar compatibilidad hacia adelante
	service                              app.ProductService
}

// NewProductHandler crea una nueva instancia del adaptador de entrada
func NewProductHandler(service app.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Create(ctx context.Context, req *pb.Producto) (*pb.Producto, error) {
	// Mapeo de gRPC Message -> Entidad de Dominio
	productDom := &domain.Product{

		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Precio:      req.Precio,
	}

	res, err := h.service.CreateProduct(ctx, productDom)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error al crear el producto: %v", err)
	}

	// Mapeo de Entidad de Dominio -> gRPC Message de retorno
	return &pb.Producto{
		Id:          res.ID,
		Nombre:      res.Nombre,
		Descripcion: res.Descripcion,
		Precio:      res.Precio,
	}, nil
}

func (h *ProductHandler) Get(ctx context.Context, req *pb.ProductId) (*pb.Producto, error) {
	res, err := h.service.GetProduct(ctx, req.Id)
	if err != nil {
		if err.Error() == "product_not_found" {
			// Manejo de error específico usando códigos nativos gRPC (NOT_FOUND = 5)
			return nil, status.Errorf(codes.NotFound, "El producto con ID %s no existe", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "Error interno: %v", err)
	}

	return &pb.Producto{
		Id:          res.ID,
		Nombre:      res.Nombre,
		Descripcion: res.Descripcion,
		Precio:      res.Precio,
	}, nil
}

func (h *ProductHandler) Update(ctx context.Context, req *pb.Producto) (*pb.Producto, error) {
	productDom := &domain.Product{
		ID:          req.Id,
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Precio:      req.Precio,
	}

	res, err := h.service.UpdateProduct(ctx, productDom)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error al actualizar el producto: %v", err)
	}

	return &pb.Producto{
		Id:          res.ID,
		Nombre:      res.Nombre,
		Descripcion: res.Descripcion,
		Precio:      res.Precio,
	}, nil
}

func (h *ProductHandler) Delete(ctx context.Context, req *pb.ProductId) (*pb.Producto, error) {
	// Primero buscamos el producto para poder retornar sus datos al eliminarlo (como pide el .proto)
	p, err := h.service.GetProduct(ctx, req.Id)
	if err != nil {
		if err.Error() == "product_not_found" {
			return nil, status.Errorf(codes.NotFound, "No se puede eliminar, el producto no existe")
		}
		return nil, status.Errorf(codes.Internal, "Error interno: %v", err)
	}

	err = h.service.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error al eliminar el producto: %v", err)
	}

	return &pb.Producto{
		Id:          p.ID,
		Nombre:      p.Nombre,
		Descripcion: p.Descripcion,
		Precio:      p.Precio,
	}, nil
}

func (h *ProductHandler) List(ctx context.Context, req *pb.Empty) (*pb.ProductosList, error) {
	products, err := h.service.ListProducts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error al listar productos: %v", err)
	}

	// Mapear la lista slice de dominio al repeated de gRPC
	var gProducts []*pb.Producto
	for _, p := range products {
		gProducts = append(gProducts, &pb.Producto{
			Id:          p.ID,
			Nombre:      p.Nombre,
			Descripcion: p.Descripcion,
			Precio:      p.Precio,
		})
	}

	return &pb.ProductosList{Items: gProducts}, nil
}
