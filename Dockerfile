# Etapa de construcción
FROM golang:1.21.0 AS builder

# Configura el directorio de trabajo
WORKDIR /build

# Copia los archivos de la aplicación
COPY . .

# Compila el binario con la arquitectura correcta
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o main ./

# Crear la imagen final con Alpine
FROM alpine:latest

# Actualiza los paquetes necesarios
RUN apk upgrade libssl3 libcrypto3 busybox

# Copia el binario desde la etapa de construcción
COPY --from=builder /build/main /app/main

# Exponer el puerto
EXPOSE 8080

# Define el ENTRYPOINT para ejecutar el binario
ENTRYPOINT ["/app/main"]
