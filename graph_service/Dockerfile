# STAGE 1: Compilación
FROM golang:1.26.2-alpine AS builder

# Instalamos git por si alguna dependencia lo requiere
RUN apk add --no-cache git

# Directorio de trabajo interno en la etapa de construcción
WORKDIR /app

# Copiamos archivos de dependencias primero para aprovechar el cache de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiamos todo el código fuente del proyecto
COPY . .

# CORRECCIÓN DE RUTA: Apuntamos la compilación a la carpeta ./cmd donde reside tu main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# STAGE 2: Ejecución (Imagen ultraligera de producción)
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

# Usamos /app como directorio de trabajo para mantener las rutas relativas coherentes
WORKDIR /app

# Copiamos el binario estático desde la etapa anterior
COPY --from=builder /app/main .

# SOLUCIÓN DE ARCHIVO: Copiamos el archivo de esquema respetando su estructura de carpetas
COPY --from=builder /app/internal/infrastructure/graphql/schema.graphql ./internal/infrastructure/graphql/schema.graphql

# Exponemos el puerto de Gin que definimos para GraphQL
EXPOSE 8080

# Comando para arrancar el contenedor desde el directorio /app
CMD ["./main"]