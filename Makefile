.PHONY: build run test clean docker-build docker-run

# Переменные
BINARY_NAME=kagent-operator
DOCKER_IMAGE=kagent-operator
DOCKER_TAG=latest

# Сборка бинарника
build:
	go build -o $(BINARY_NAME) .

# Запуск локально
run:
	go run .

# Тесты
test:
	go test -v ./...

# Очистка
clean:
	go clean
	rm -f $(BINARY_NAME)

# Сборка Docker образа
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Запуск Docker контейнера
docker-run:
	docker run -p 8080:8080 \
		-e PORT=8080 \
		-e MCP_SERVER_URL=http://host.docker.internal:3000/mcp \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Установка зависимостей
deps:
	go mod download
	go mod tidy
