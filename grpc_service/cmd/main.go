package main

import (
	"grpc_service/internal/app"
	"grpc_service/internal/infrastructure/database"
	gHandler "grpc_service/internal/infrastructure/grpc"
	"grpc_service/internal/infrastructure/grpc/pb"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Definir la cadena de conexión (DSN) para PostgreSQL
	// Reemplaza estos valores con tus credenciales locales o las de Render / Supabase
	dsn := "host=postgres_db user=admin password=root dbname=inventarios port=5432 sslmode=disable TimeZone=America/Bogota"

	// Inicializar la Base de Datos con el ORM (GORM)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fallo al conectar a PostgreSQL: %v", err)
	}

	// 2. Inyección de Dependencias (Arquitectura Hexagonal)
	productRepo, err := database.NewGormRepository(db)
	if err != nil {
		log.Fatalf("Fallo en la automigración del ORM: %v", err)
	}
	productService := app.NewProductService(productRepo)
	grpcHandler := gHandler.NewProductHandler(productService)

	// 3. Levantar el Servidor gRPC en un puerto de fondo (Goroutine)
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Fallo al escuchar en puerto 50051: %v", err)
		}

		baseGrpcServer := grpc.NewServer()
		pb.RegisterProductServiceServer(baseGrpcServer, grpcHandler)

		log.Println("🚀 Servidor gRPC corriendo en el puerto :50051")
		if err := baseGrpcServer.Serve(lis); err != nil {
			log.Fatalf("Fallo al servir gRPC: %v", err)
		}
	}()

	// 4. Levantar el Servidor HTTP de Gin en el hilo principal (puerto 8080)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"service": "gRPC Backend Multi-server con Postgres",
		})
	})

	log.Println("🌐 Servidor Gin HTTP corriendo en el puerto :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Fallo al arrancar Gin: %v", err)
	}
}
