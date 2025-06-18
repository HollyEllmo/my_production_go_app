# Устранение неполадок

## Проблемы с импортами proto-зависимостей

### Симптомы

```
could not import github.com/HollyEllmo/my_proto_repo/gen/go/prod_service/products/v1 (no required module provides package)
```

### Причина

Несоответствие между именем репозитория в GitHub (`my_proto_repo` с подчеркиванием) и внутренними зависимостями proto-модулей (которые ожидают `my-proto-repo` с дефисом).

### Правильная архитектура решения

В нашем проекте используется схема с replace директивами:

1. **В импортах Go-кода** используется `my-proto-repo` (с дефисом)
2. **В go.mod require секции** также используется `my-proto-repo` (с дефисом)
3. **Replace директивы** перенаправляют запросы на реальный репозиторий `my_proto_repo` (с подчеркиванием)

### Правильная структура go.mod

```go
require (
    github.com/HollyEllmo/my-proto-repo/gen/go/prod_service v0.0.0-YYYYMMDDHHMMSS-commithash
)

replace github.com/HollyEllmo/my-proto-repo/gen/go/prod_service => github.com/HollyEllmo/my_proto_repo/gen/go/prod_service v0.0.0-YYYYMMDDHHMMSS-commithash
replace github.com/HollyEllmo/my-proto-repo/gen/go/filter => github.com/HollyEllmo/my_proto_repo/gen/go/filter v0.0.0-YYYYMMDDHHMMSS-commithash
```

### Правильные импорты в Go-коде

```go
import pb_prod_products "github.com/HollyEllmo/my-proto-repo/gen/go/prod_service/products/v1"
import pb_common_filter "github.com/HollyEllmo/my-proto-repo/gen/go/filter/v1"
```

3. **Команды для исправления**:

   ```bash
   cd app
   go mod tidy
   go build -o ../bin/app ./cmd/app
   ```

4. **Автоматическое обновление**:
   ```bash
   ./update-proto-deps.sh
   ```

### Профилактика

- Всегда используйте скрипт `update-proto-deps.sh` для обновления зависимостей
- После каждого обновления proto-репозитория запускайте `go mod tidy`
- Проверяйте сборку после обновления: `go build`

## Быстрая диагностика

```bash
# Проверить состояние модулей
go list -m all

# Проверить конкретную зависимость
go list -m github.com/HollyEllmo/my_proto_repo/gen/go/prod_service

# Очистить кеш модулей (крайняя мера)
go clean -modcache
```
