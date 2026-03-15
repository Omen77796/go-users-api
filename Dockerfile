# Imagen base oficial de Go
FROM golang:1.25

# Carpeta de trabajo dentro del contenedor
WORKDIR /app

# Copiar archivos go.mod y go.sum primero
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar todo el proyecto
COPY . .

# Compilar la aplicación
RUN go build -o main ./cmd/api 

# Exponer el puerto
EXPOSE 8080

# Comando para ejecutar
CMD ["./main"]