# STAGE 1: Compilación
FROM golang:1.26.2-alpine AS builder

# Instalamos git por si alguna dependencia lo requiere
RUN apk add --no-cache git

# Directorio de trabajo interno
WORKDIR /app

# Copiamos archivos de dependencias primero para aprovechar el cache de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiamos todo el código fuente (incluyendo la carpeta internal/ con los archivos .pb.go)
COPY . .

# CORRECCIÓN DE RUTA: Apuntamos la compilación a la carpeta ./cmd donde reside tu main.go de gRPC
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# STAGE 2: Ejecución (Imagen ultraligera de producción)
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copiamos el binario estático desde la etapa de compilación
COPY --from=builder /app/main .

# EXPOSE: Cambiamos el puerto al :50051 que es el puerto estándar que configuramos para gRPC
EXPOSE 50051
# Si además de gRPC dejas corriendo Gin en paralelo en ese proyecto, puedes exponer ambos:
# EXPOSE 50051 8080

# Comando por defecto para arrancar el contenedor
CMD ["./main"]