# Dockerfile

# --- Etapa de Construcción (Build Stage) ---
# Usamos una imagen base que incluya el compilador de Go.
# Esto mantiene la imagen final más limpia y pequeña.
FROM golang:1.24.4-alpine AS builder

# Establece el directorio de trabajo dentro del contenedor.
WORKDIR /app

# Copia los archivos go.mod y go.sum (si usas módulos Go)
# para que Docker pueda descargar las dependencias primero.
COPY go.mod ./
COPY go.sum ./

# Descarga las dependencias del módulo.
# Si tus dependencias no cambian a menudo, esta capa se cacheará.
RUN go mod download

# Copia todo el código fuente de tu aplicación al directorio de trabajo.
COPY . .

# Construye tu aplicación. El resultado es un ejecutable estático.
# CGO_ENABLED=0 deshabilita la vinculación de CGO, haciendo que el binario sea completamente estático
# y más fácil de mover.
# -o /app/streaming-app nombra tu ejecutable.
RUN CGO_ENABLED=0 go build -o /app/streaming-app .

# --- Etapa Final (Run Stage) ---
# Usamos una imagen base más pequeña para la imagen final.
# 'scratch' es la imagen más pequeña posible, no contiene nada.
FROM scratch

# Establece metadatos sobre el puerto que la aplicación escuchará (si aplica).
# Aunque nuestra app no escucha puertos, es buena práctica.
EXPOSE 8080

# Copia el ejecutable desde la etapa de construcción a la imagen final.
COPY --from=builder /app/streaming-app /streaming-app

# Define el comando por defecto para ejecutar tu aplicación cuando el contenedor se inicie.
CMD ["/streaming-app"]