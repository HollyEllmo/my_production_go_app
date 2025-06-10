# Переменные
APP_BIN = app/build/app

# Основные цели
.PHONY: lint build clean swagger migrate migrate.down migrate.up

# Линтинг кода
lint:
	golangci-lint run

# Сборка приложения
build: clean $(APP_BIN)

# Компиляция бинарного файла
$(APP_BIN):
	go build -o $(APP_BIN) ./app/cmd/app/main.go

# Очистка собранных файлов
clean:
	rm -rf ./app/build || true

# Генерация Swagger документации
swagger:
	swag init -g ./app/cmd/app/main.go -o ./app/docs

# Миграции базы данных
migrate:
	$(APP_BIN) migrate -version $(version)

migrate.down:
	$(APP_BIN) migrate -seq down

migrate.up:
	$(APP_BIN) migrate -seq up