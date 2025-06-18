#!/bin/bash

# Простой скрипт для обновления proto зависимостей

set -e

echo "🔄 Обновление proto зависимостей..."

# Переходим в директорию app
cd "$(dirname "$0")/app" || exit 1

echo "🧹 Очистка кэша..."
go clean -modcache

echo "📥 Обновление зависимостей..."
# Используем go get для обновления до latest
go get github.com/HollyEllmo/my_proto_repo/gen/go/prod_service@latest
go get github.com/HollyEllmo/my_proto_repo/gen/go/filter@latest

echo "🧹 Выполняется go mod tidy..."
go mod tidy

echo "✅ Зависимости успешно обновлены!"
echo "🐳 Теперь можно запускать Docker: docker-compose up --build"
