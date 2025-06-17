# Интеграция gRPC сервиса prod_service

Этот документ описывает, как в проекте настроена интеграция с gRPC сервисом `prod_service`.

## Что было сделано

### 1. Настройка зависимостей

В `go.mod` добавлены зависимости на proto-репозиторий:

```go
require (
    github.com/HollyEllmo/my-proto-repo/gen/go/prod_service v0.0.0-00010101000000-000000000000
    // ... другие зависимости
)

replace github.com/HollyEllmo/my-proto-repo/gen/go/prod_service => /Users/vadim/go/src/github.com/HollyEllmo/my-proto-repo/gen/go/prod_service
replace github.com/HollyEllmo/my-proto-repo/gen/go/filter => /Users/vadim/go/src/github.com/HollyEllmo/my-proto-repo/gen/go/filter
```

### 2. Интеграция в app.go

В файле `internal/app/app.go`:

- Добавлен импорт gRPC клиента
- В структуру `App` добавлены поля для gRPC клиента и соединения
- В функции `NewApp()` добавлена инициализация gRPC клиента
- Добавлены методы для работы с gRPC сервисом
- Добавлено тестирование соединения при запуске

### 3. HTTP обработчики

Создан файл `internal/handler/product.go` с HTTP обработчиками, которые используют gRPC клиент:

- `GET /api/v1/products` - получение всех продуктов
- `POST /api/v1/products` - создание нового продукта

### 4. Методы для работы с gRPC

В `app.go` добавлены методы:

- `GetProductServiceClient()` - получение gRPC клиента
- `CallProductService()` - пример использования gRPC клиента
- `testProductServiceConnection()` - тестирование соединения
- `Shutdown()` - корректное закрытие gRPC соединения

## Использование

### Пример HTTP запросов

1. **Получить все продукты:**

```bash
curl -X GET http://localhost:8080/api/v1/products
```

2. **Создать новый продукт:**

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "Test Description",
    "price": "99.99",
    "category_id": "category-1"
  }'
```

### Использование gRPC клиента в коде

```go
// Получить клиент
client := app.GetProductServiceClient()
if client == nil {
    // Сервис недоступен
    return
}

// Вызвать метод
ctx := context.Background()
req := &pb_prod_products.AllProductsRequest{}
resp, err := client.AllProducts(ctx, req)
if err != nil {
    // Обработка ошибки
    return
}

// Использовать результат
products := resp.Product
```

## Настройка

### Адрес gRPC сервера

В настоящее время адрес gRPC сервера захардкожен в коде: `localhost:50051`

Для продакшена рекомендуется вынести это в конфигурацию.

### Таймауты

Установлены следующие таймауты:

- Тестирование соединения: 5 секунд
- HTTP обработчики: 10 секунд

### Обработка ошибок

- Если gRPC сервис недоступен при запуске, приложение продолжит работу без него
- HTTP обработчики возвращают соответствующие HTTP статус коды
- Все ошибки логируются

## Зависимости

- `google.golang.org/grpc` - основная библиотека gRPC
- `google.golang.org/grpc/credentials/insecure` - для незащищенного соединения
- Сгенерированные proto файлы из `my-proto-repo`

## Важные моменты

1. **Replace директивы**: Используются для локальной разработки. В продакшене нужно будет публиковать proto-репозиторий.

2. **Конфигурация**: Адрес gRPC сервера стоит вынести в конфигурацию.

3. **Graceful shutdown**: Добавлен корректный shutdown для закрытия gRPC соединения.

4. **Тестирование**: Добавлено автоматическое тестирование соединения при запуске.
