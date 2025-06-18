#!/bin/bash

# Скрипт для обновления proto зависимостей до последней версии из main ветки

set -e

echo "🔄 Обновление proto зависимостей..."

# Получаем последний коммит из main ветки
LATEST_COMMIT=$(git ls-remote https://github.com/HollyEllmo/my_proto_repo.git HEAD | cut -f1)
echo "📍 Последний коммит: $LATEST_COMMIT"

# Создаем псевдоверсию из коммита
TIMESTAMP=$(date -u +"%Y%m%d%H%M%S")
PSEUDO_VERSION="v0.0.0-${TIMESTAMP}-${LATEST_COMMIT:0:12}"

echo "🏷️  Новая версия: $PSEUDO_VERSION"

# Путь к go.mod файлу
GO_MOD_FILE="./app/go.mod"

# Обновляем версии в go.mod
sed -i '' "s/github.com\/HollyEllmo\/my-proto-repo\/gen\/go\/prod_service v[0-9].*$/github.com\/HollyEllmo\/my-proto-repo\/gen\/go\/prod_service $PSEUDO_VERSION/" "$GO_MOD_FILE"
sed -i '' "s/github.com\/HollyEllmo\/my-proto-repo\/gen\/go\/filter v[0-9].*$/github.com\/HollyEllmo\/my-proto-repo\/gen\/go\/filter $PSEUDO_VERSION/" "$GO_MOD_FILE"

# Обновляем replace директивы
sed -i '' "s/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/prod_service v[0-9].*$/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/prod_service $PSEUDO_VERSION/" "$GO_MOD_FILE"
sed -i '' "s/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/filter v[0-9].*$/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/filter $PSEUDO_VERSION/" "$GO_MOD_FILE"

echo "📝 go.mod обновлен"

# Переходим в директорию app и запускаем go mod tidy
cd app
echo "🧹 Выполняется go mod tidy..."
go mod tidy

echo "✅ Зависимости успешно обновлены!"
echo "🐳 Теперь можно запускать Docker: docker-compose up --build"
