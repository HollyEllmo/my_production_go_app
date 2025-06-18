APP_BIN = app/build/app

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build: clean $(APP_BIN)

$(APP_BIN):
	go build -o $(APP_BIN) ./app/cmd/app/main.go

.PHONY: clean
clean:
	rm -rf ./app/build || true

.PHONY: swagger
swagger:
	swag init -g ./app/cmd/app/main.go -o ./app/docs

.PHONY: migrate
migrate:
	$(APP_BIN) migrate -version $(version)

.PHONY: migrate.down
migrate.down:
	$(APP_BIN) migrate -seq down

.PHONY: migrate.up
migrate.up:
	$(APP_BIN) migrate -seq up

# Docker commands
.PHONY: docker.up
docker.up:
	docker-compose up --build -d

.PHONY: docker.down
docker.down:
	docker-compose down

.PHONY: docker.restart
docker.restart:
	docker-compose restart

.PHONY: docker.logs
docker.logs:
	docker-compose logs -f

# App-specific commands
.PHONY: app.restart
app.restart:
	docker-compose restart app

.PHONY: app.rebuild
app.rebuild:
	docker-compose stop app
	docker-compose build app
	docker-compose up -d app

.PHONY: app.logs
app.logs:
	docker-compose logs -f app

.PHONY: db.restart
db.restart:
	docker-compose restart ps-psql

.PHONY: db.logs
db.logs:
	docker-compose logs -f ps-psql

# Development workflow
.PHONY: dev.restart
dev.restart: app.rebuild app.logs

# Development mode with hot reload
.PHONY: dev.up
dev.up:
	docker-compose --profile dev up --build -d app-dev

.PHONY: dev.down
dev.down:
	docker-compose --profile dev down

.PHONY: dev.logs
dev.logs:
	docker-compose logs -f app-dev

.PHONY: dev.restart-hot
dev.restart-hot:
	docker-compose --profile dev restart app-dev

# Status and monitoring
.PHONY: status
status:
	@echo "=== Docker Containers Status ==="
	docker ps --filter "name=ps-psql" --filter "name=my-production-service"
	@echo ""
	@echo "=== Services Health Check ==="
	@curl -s http://localhost:30001/api/heartbeat > /dev/null && echo "✅ HTTP Service: OK" || echo "❌ HTTP Service: DOWN"
	@docker exec ps-psql pg_isready -U postgres > /dev/null 2>&1 && echo "✅ Database: OK" || echo "❌ Database: DOWN"

