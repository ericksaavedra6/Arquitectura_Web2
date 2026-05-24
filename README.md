# Arquitectura_Web2
# Inicialización del módulo de Go
```
go mod init grpc_service <br>
go mod init graph_service
```
# Dependencias Core del Servidor gRPC y Protocol Buffers
```
go get google.golang.org/grpc<br>
go get google.golang.org/protobuf
```
# Framework Web Gin (Utilizado para hosting del servidor GraphQL y HTTP status)
```
go get [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
```

# Módulos del ORM GORM y Driver nativo para bases de datos PostgreSQL
```
go get gorm.io/gorm<br>
go get gorm.io/driver/postgres
```
# 📐 2. Comandos de Generación de Código (Compilación de Contratos)
## A. Compilación del Contrato gRPC (.proto)
<br>La definición de servicios y estructuras de mensajes estructurados se realiza en el archivo product.proto. 
<br>Para compilar los stubs y la interfaz del servidor en Go, ejecute desde el directorio donde se ubica el archivo:<br>
```
Bash

protoc --go_out=. --go_opt=paths=source_relative \\
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \\
       internal/infrastructure/grpc/proto/product.proto
```
### Herramientas requeridas instaladas en el sistema: 
protoc, protoc-gen-go y protoc-gen-go-grpc.

## B. Generación del Esquema y Resolvers de GraphQL
Si utiliza herramientas de generación automática de código basadas en esquemas estructurados SDL (como gqlgen), el comando ejecutado para mapear el archivo schema.graphql a structs nativos de Go es:<br>
```
Bash
go run [github.com/99designs/gqlgen](https://github.com/99designs/gqlgen) generate
```
## 🐳 3. Comando para Levantar todo el Ecosistema
El despliegue está completamente automatizado y orquestado mediante contenedores aislados. 
<br>Gracias a la directiva de healthcheck implementada en la receta, las aplicaciones backend esperarán de forma segura a que la instancia relacional esté lista para aceptar conexiones antes de inicializar sus puertos.
Para construir las imágenes e iniciar la base de datos de PostgreSQL junto con los servidores en segundo plano, ejecute:
```
Bash
docker-compose up --build
```
### Puertos Expuestos en el Host:

:50051 -> Endpoints nativos del servidor gRPC.
<br>:8080 -> Endpoint de GraphQL (/graphql) y consultas HTTP sobre Gin.
<br>:5432 -> Motor transaccional de PostgreSQL (inventarios).