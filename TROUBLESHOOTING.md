# Устранение неполадок

## Проблемы с импортами proto-зависимостей

### Симптомы

```
could not import github.com/HollyEllmo/my_proto_repo/gen/go/prod_service/products/v1 (no required module provides package)
```

### Причина

Несоответствие между именем репозитория в GitHub (`my_proto_repo` с подчеркиванием) и путями в go.mod (иногда используется `my-proto-repo` с дефисом).

### Решение

1. **Проверьте go.mod файл**:

   - В секции `require` должны быть модули с правильными именами и версиями
   - В секции `replace` должны быть корректные переадресации

2. **Правильная структура go.mod**:

   ```go
   require (
       github.com/HollyEllmo/my_proto_repo/gen/go/prod_service v0.0.0-YYYYMMDDHHMMSS-commithash
   )

   replace github.com/HollyEllmo/my-proto-repo/gen/go/filter => github.com/HollyEllmo/my_proto_repo/gen/go/filter v0.0.0-YYYYMMDDHHMMSS-commithash
   replace github.com/HollyEllmo/my-proto-repo/gen/go/prod_service => github.com/HollyEllmo/my_proto_repo/gen/go/prod_service v0.0.0-YYYYMMDDHHMMSS-commithash
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
