package graphql

import (
	"context"
	"errors"
	"strconv"

	"graph_service/internal/app"
	"graph_service/internal/domain"

	// Importamos el paquete raíz para acceder al tipo estricto ID
	"github.com/graph-gophers/graphql-go"
)

// Resolver es el punto de entrada central para todas las operaciones de GraphQL
type Resolver struct {
	service app.ProductService
}

// NewResolver inicializa el adaptador de entrada con el caso de uso
func NewResolver(service app.ProductService) *Resolver {
	return &Resolver{service: service}
}

// --- STRUCT INTERMEDIO CORREGIDO CON TIPO DE LA LIBRERÍA ---
type productResolver struct {
	id          graphql.ID // Cambiado de string a graphql.ID para satisfacer la reflexión
	nombre      string
	descripcion string
	precio      float64
}

// Métodos de lectura con el tipo estricto que exige el esquema
func (r *productResolver) ID() graphql.ID      { return r.id }
func (r *productResolver) Nombre() string      { return r.nombre }
func (r *productResolver) Descripcion() string { return r.descripcion }
func (r *productResolver) Precio() float64     { return r.precio }

// Función auxiliar de conversión ajustada
func toGraphQL(p *domain.Product) *productResolver {
	return &productResolver{
		id:          graphql.ID(strconv.FormatInt(p.ID, 10)), // Casteo limpio a graphql.ID
		nombre:      p.Nombre,
		descripcion: p.Descripcion,
		precio:      p.Precio,
	}
}

// --- QUERIES (LECTURA) ---

// Cambiamos el argumento esperado en la estructura a graphql.ID
func (r *Resolver) Get(ctx context.Context, args struct{ ID graphql.ID }) (*productResolver, error) {
	// Convertimos el graphql.ID (que es un alias de string) a int64 para el Core Hexagonal
	idInt, err := strconv.ParseInt(string(args.ID), 10, 64)
	if err != nil {
		return nil, errors.New("el ID provisto no es valido")
	}

	res, err := r.service.GetProduct(ctx, idInt)
	if err != nil {
		if err.Error() == "product_not_found" {
			return nil, errors.New("el producto solicitado no existe")
		}
		return nil, err
	}

	return toGraphQL(res), nil
}

func (r *Resolver) List(ctx context.Context) ([]*productResolver, error) {
	products, err := r.service.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	var lResolvers []*productResolver
	for _, p := range products {
		prodCopy := p
		lResolvers = append(lResolvers, toGraphQL(&prodCopy))
	}

	return lResolvers, nil
}

// --- MUTATIONS (ESCRITURA) ---

func (r *Resolver) Create(ctx context.Context, args struct {
	Nombre      string
	Descripcion string
	Precio      float64
}) (*productResolver, error) {
	productDom := &domain.Product{
		Nombre:      args.Nombre,
		Descripcion: args.Descripcion,
		Precio:      args.Precio,
	}

	res, err := r.service.CreateProduct(ctx, productDom)
	if err != nil {
		return nil, err
	}

	return toGraphQL(res), nil
}

// Cambiamos el argumento ID de la mutación a graphql.ID
func (r *Resolver) Update(ctx context.Context, args struct {
	ID          graphql.ID
	Nombre      string
	Descripcion string
	Precio      float64
}) (*productResolver, error) {
	idInt, err := strconv.ParseInt(string(args.ID), 10, 64)
	if err != nil {
		return nil, errors.New("ID invalido para actualizacion")
	}

	productDom := &domain.Product{
		ID:          idInt,
		Nombre:      args.Nombre,
		Descripcion: args.Descripcion,
		Precio:      args.Precio,
	}

	res, err := r.service.UpdateProduct(ctx, productDom)
	if err != nil {
		return nil, err
	}

	return toGraphQL(res), nil
}

// Cambiamos el argumento ID de la mutación a graphql.ID
func (r *Resolver) Delete(ctx context.Context, args struct{ ID graphql.ID }) (*productResolver, error) {
	idInt, err := strconv.ParseInt(string(args.ID), 10, 64)
	if err != nil {
		return nil, errors.New("ID invalido para eliminacion")
	}

	p, err := r.service.GetProduct(ctx, idInt)
	if err != nil {
		return nil, errors.New("no se puede eliminar un producto inexistente")
	}

	err = r.service.DeleteProduct(ctx, idInt)
	if err != nil {
		return nil, err
	}

	return toGraphQL(p), nil
}
