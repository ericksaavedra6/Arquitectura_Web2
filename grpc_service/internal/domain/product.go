package domain

// Product representa la entidad de negocio limpia en nuestro Core
type Product struct {
	ID          string `gorm:"primaryKey"` // Anotación para el ORM que usaremos después
	Nombre      string
	Descripcion string
	Precio      float64
}
