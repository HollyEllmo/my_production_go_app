# Используем официальный образ Go как базовый для сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем Git для загрузки зависимостей из Git
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /build

# Копируем go.mod и go.sum файлы из app директории
COPY app/go.mod app/go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код приложения
COPY app/ ./

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Используем минимальный образ для финального контейнера
FROM alpine:latest

# Добавляем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем собранное приложение из builder stage
COPY --from=builder /build/main .

# Копируем миграции
COPY migrations ./migrations

# Создаем директорию для конфигов
RUN mkdir -p /app/configs

# Открываем порты для gRPC и HTTP
EXPOSE 30000 30001

# Запускаем приложение
CMD ["./main"]
