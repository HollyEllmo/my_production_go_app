#!/bin/bash

# Скрипт для обновления proto зависимостей до последней версии из main ветки

set -e

echo "🔄 Обновление proto зависимостей..."

# Получаем последний коммит из main ветки
LATEST_COMMIT=$(git ls-remote https://github.com/HollyEllmo/my_proto_repo.git HEAD | cut -f1)
echo "📍 Последний коммит: $LATEST_COMMIT"

# Получаем время коммита в правильном формате
COMMIT_TIME=$(git ls-remote --heads https://github.com/HollyEllmo/my_proto_repo.git main | head -n1 | cut -f1 | xargs -I {} curl -s "https://api.github.com/repos/HollyEllmo/my_proto_repo/commits/{}" | grep '"date"' | head -n1 | sed 's/.*"\([0-9]\{4\}-[0-9]\{2\}-[0-9]\{2\}T[0-9]\{2\}:[0-9]\{2\}:[0-9]\{2\}Z\)".*/\1/')

# Конвертируем в формат для псевдоверсии (YYYYMMDDHHMMSS)
FORMATTED_TIME=$(date -u -j -f "%Y-%m-%dT%H:%M:%SZ" "$COMMIT_TIME" "+%Y%m%d%H%M%S" 2>/dev/null || echo "20250618070144")

# Создаем псевдоверсию из коммита
PSEUDO_VERSION="v0.0.0-${FORMATTED_TIME}-${LATEST_COMMIT:0:12}"

echo "🏷️  Новая версия: $PSEUDO_VERSION"

# Путь к go.mod файлу
GO_MOD_FILE="./app/go.mod"

# Обновляем версии в go.mod
sed -i '' "s/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/prod_service v[0-9].*$/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/prod_service $PSEUDO_VERSION/" "$GO_MOD_FILE"
sed -i '' "s/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/filter v[0-9].*$/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/filter $PSEUDO_VERSION/" "$GO_MOD_FILE"

# Обновляем replace директивы для proto-repo (с дефисом) -> my_proto_repo (с подчеркиванием)
sed -i '' "s/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/prod_service v[0-9].*$/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/prod_service $PSEUDO_VERSION/" "$GO_MOD_FILE"
sed -i '' "s/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/filter v[0-9].*$/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/filter $PSEUDO_VERSION/" "$GO_MOD_FILE"

# Также обновляем require секцию (используется с дефисом как псевдоним)
sed -i '' "s/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/prod_service v[0-9].*$/github.com\/HollyEllmo\/my_proto_repo\/gen\/go\/prod_service $PSEUDO_VERSION/" "$GO_MOD_FILE"

echo "📝 go.mod обновлен"

# Переходим в директорию app и запускаем go mod tidy
cd app
echo "🧹 Выполняется go mod tidy..."
go mod tidy

echo "✅ Зависимости успешно обновлены!"
echo "🐳 Теперь можно запускать Docker: docker-compose up --build"
