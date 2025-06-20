# My Production Service

Производственный Go сервис с чистой архитектурой.

## Структура проекта

```
production_service/
├── app/
│   ├── cmd/
│   │   └── app/
│   │       └── main.go          # Точка входа в приложение
│   ├── internal/                # Приватная логика приложения
│   │   ├── config/              # Конфигурация
│   │   ├── handler/             # HTTP обработчики
│   │   ├── service/             # Бизнес-логика
│   │   └── repository/          # Слой доступа к данным
│   └── pkg/                     # Публичные пакеты
│       ├── database/            # Работа с базой данных
│       └── logger/              # Логирование
├── migrations/                  # Миграции базы данных
├── .dockerignore
├── .gitignore
├── .gitlab-ci.yml              # CI/CD конфигурация
├── Dockerfile                  # Docker конфигурация
├── Makefile                    # Команды для сборки и разработки
└── go.mod                      # Go модули
```

## Быстрый старт

### Предварительные требования

- Go 1.24+
- PostgreSQL (опционально)
- Docker (для контейнеризации)

### Установка зависимостей

```bash
make deps
```

### Запуск приложения

```bash
# В режиме разработки
make run

# Сборка и запуск
make build
./bin/main
```

### Переменные окружения

Приложение поддерживает конфигурацию через переменные окружения и файлы конфигурации (YAML, ENV).

**Основные переменные:**

- `IS_DEBUG` - режим отладки (по умолчанию: false)
- `IS_DEV` - режим разработки (по умолчанию: false)
- `LISTEN_TYPE` - тип прослушивания (по умолчанию: port)
- `BIND_IP` - IP адрес для привязки (по умолчанию: 0.0.0.0)
- `PORT` - порт для запуска сервера (по умолчанию: 10000)
- `ADMIN_EMAIL` - email администратора (обязательно)
- `ADMIN_PASSWORD` - пароль администратора (обязательно)

**Способы конфигурации:**

1. Переменные окружения
2. Файл `config.yaml`
3. Файл `.env`

**Примеры файлов конфигурации:**

- `app/config.yaml.example` - пример YAML конфигурации
- `app/.env.example` - пример ENV конфигурации

## Docker

### Сборка образа

```bash
make docker-build
```

### Запуск в контейнере

```bash
make docker-run
```

## Доступные команды

Все доступные команды можно посмотреть:

```bash
make help
```

## API Endpoints

- `GET /health` - проверка здоровья сервиса

## Разработка

### Тестирование

```bash
make test
```

### Линтинг

```bash
make lint
```

### Форматирование кода

```bash
make fmt
```

## Миграции

Для работы с миграциями используется [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
# Применить миграции
make migrate-up

# Откатить миграции
make migrate-down

# Создать новую миграцию
make migrate-create name=new_migration_name
```

## CI/CD

Проект настроен для работы с GitLab CI/CD. Пайплайн включает:

- Тестирование
- Линтинг
- Сборку
- Деплой

## Архитектура

Проект следует принципам чистой архитектуры:

- **cmd/app** - точка входа
- **internal** - приватная логика приложения
- **pkg** - переиспользуемые пакеты
- **migrations** - миграции базы данных

Слои:

- **Handler** - HTTP обработчики
- **Service** - бизнес-логика
- **Repository** - доступ к данным
