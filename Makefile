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
	cd app && go build -o build/app ./cmd/app/main.go

# Очистка собранных файлов
clean:
	rm -rf ./app/build || true

# Генерация Swagger документации
swagger:
	swag init -g ./app/cmd/app/main.go -o ./app/docs

# Миграции базы данных
migrate:
	cd app && $(APP_BIN) migrate -version $(version)

migrate.down:
	cd app && ./build/app migrate -seq down

migrate.up:
	cd app && ./build/app migrate -seq up