package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"graph_service/internal/app"
	"graph_service/internal/infrastructure/database"
	gqlHandler "graph_service/internal/infrastructure/graphql"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/graphql-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Inicializar la Base de Datos con el ORM (GORM y PostgreSQL)
	dsn := "host=postgres_db user=admin password=root dbname=inventarios port=5432 sslmode=disable TimeZone=America/Bogota"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fallo al conectar a PostgreSQL: %v", err)
	}

	// 2. Inyección de Dependencias (Mismo núcleo Hexagonal)
	productRepo, err := database.NewGormRepository(db)
	if err != nil {
		log.Fatalf("Fallo en la automigración del ORM: %v", err)
	}
	productService := app.NewProductService(productRepo)
	resolver := gqlHandler.NewResolver(productService)

	// 3. Leer y Compilar el Esquema de GraphQL
	schemaFile, err := ioutil.ReadFile("internal/infrastructure/graphql/schema.graphql")
	if err != nil {
		log.Fatalf("No se pudo leer el archivo de esquema graphql: %v", err)
	}

	// Compila el esquema asociándolo a las funciones del resolver maestro
	parsedSchema, err := graphql.ParseSchema(string(schemaFile), resolver)
	if err != nil {
		log.Fatalf("Error al parsear el esquema GraphQL: %v", err)
	}

	// 4. Configurar Gin y exponer el endpoint unificado de GraphQL
	r := gin.Default()

	// Interceptamos la petición POST y procesamos el cuerpo del request manualmente con Exec
	r.POST("/graphql", func(c *gin.Context) {
		// Estructura estándar que espera recibir una petición de GraphQL por HTTP JSON
		var requestBody struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}

		// Parsear el JSON entrante
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
			return
		}

		// Ejecutar la consulta contra el esquema de GraphQL usando el método Exec nativo
		response := parsedSchema.Exec(c.Request.Context(), requestBody.Query, requestBody.OperationName, requestBody.Variables)

		// Retornar el resultado estructurado de GraphQL (contiene data y/o errors)
		c.JSON(http.StatusOK, response)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"service": "GraphQL Engine con ejecutor nativo sobre Gin",
		})
	})

	log.Println("🚀 Servidor GraphQL con Gin listo y escuchando en el puerto :8080")
	log.Println("🌐 Consume tu API declarativa a través del endpoint POST http://localhost:8080/graphql")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Fallo al arrancar el servidor Gin: %v", err)
	}
}
