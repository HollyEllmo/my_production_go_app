# Использование Proto Dependencies

## Текущая настройка

Ваш проект теперь настроен для получения proto зависимостей напрямую из Git репозитория `https://github.com/HollyEllmo/my_proto_repo.git`.

### Текущие зависимости

- `github.com/HollyEllmo/my-proto-repo/gen/go/prod_service` - gRPC сервис продуктов
- `github.com/HollyEllmo/my-proto-repo/gen/go/filter` - сервис фильтрации

### Как это работает

1. **Псевдоверсии**: Используются псевдоверсии формата `v0.0.0-20250618070144-d25e3c474bb3`, которые указывают на конкретный коммит
2. **Replace директивы**: Перенаправляют зависимости с `my-proto-repo` на `my_proto_repo` (исправляют несоответствие имен)
3. **Автоматическая загрузка**: Go модули загружаются автоматически из Git при сборке

## Обновление до последней версии

### Автоматическое обновление

Используйте скрипт `update-proto-deps.sh`:

```bash
./update-proto-deps.sh
```

Этот скрипт:

- Получает последний коммит из main ветки
- Создает новую псевдоверсию
- Обновляет go.mod файл
- Запускает `go mod tidy`

### Ручное обновление

1. Получите хеш последнего коммита:

```bash
git ls-remote https://github.com/HollyEllmo/my_proto_repo.git HEAD
```

2. Создайте псевдоверсию:

```
v0.0.0-YYYYMMDDHHMMSS-COMMIT_HASH
```

3. Обновите go.mod файл с новой версией

4. Запустите:

```bash
cd app && go mod tidy
```

## Docker сборка

Теперь ваш проект можно собирать в Docker:

```bash
docker-compose up --build
```

Docker будет:

1. Загружать зависимости из Git
2. Собирать проект
3. Создавать контейнер

## Получение последних изменений

### Всегда получать последние изменения

**⚠️ Предупреждение**: Это не рекомендуется для продакшена, но допустимо для разработки.

Текущая настройка позволяет получать изменения из main ветки через:

- Запуск скрипта `update-proto-deps.sh`
- Пересборку Docker образа

### Рекомендуемый подход для продакшена

Для продакшена рекомендуется:

1. Использовать теги в proto репозитории
2. Привязываться к конкретным версиям
3. Контролировать обновления вручную

## Структура файлов

```
my-production-service/
├── app/
│   ├── go.mod          # Go зависимости (обновлен)
│   ├── go.sum          # Чексуммы зависимостей
│   └── ...
├── Dockerfile          # Docker конфигурация
├── docker-compose.yml  # Docker Compose
└── update-proto-deps.sh # Скрипт обновления
```

## Troubleshooting

### Проблемы с именами модулей

Если возникают проблемы с именами модулей:

- Проверьте, что используете `my-proto-repo` (с дефисом) в импортах
- Убедитесь, что replace директивы указывают на `my_proto_repo` (с подчеркиванием)

### Проблемы с версиями

Если Go не может найти версию:

- Проверьте, что коммит существует в репозитории
- Убедитесь, что псевдоверсия правильно сформирована
- Очистите модульный кэш: `go clean -modcache`

### Проблемы с Docker

Если Docker не может загрузить зависимости:

- Убедитесь, что репозиторий публичный
- Проверьте доступ к интернету в контейнере
- Убедитесь, что go.mod файл правильно скопирован в контейнер
