# KAgent Operator

Go сервис для анализа состояния Kubernetes кластеров через MCP Server с использованием tools из [kagent-dev/tools](https://github.com/kagent-dev/tools).

## Архитектура

Этот проект реализует **MCP Client**, который используется внутри **DAAI (Master Service)** для взаимодействия с MCP Server.

### Полная архитектура DAAI System

Система состоит из следующих компонентов:

- **DAAI (Master Service)** - центральный сервис с эндпоинтами `agent_action_1`, `agent_action_2`, ...
- **MCP Client** - клиент для взаимодействия с MCP Server (реализован в этом проекте)
- **KaaS** - сервис управления кластерами (REST/gRPC), предоставляет информацию о кластерах
- **MCP Server (kmcp)** - сервер в Parent Cluster, запускающий kagent tools
- **kagent tools** - инструменты для анализа Kubernetes кластеров из [kagent-dev/tools](https://github.com/kagent-dev/tools)
- **Parent Cluster** - основной Kubernetes кластер, содержащий множество Children Clusters
- **Children Clusters** - дочерние Kubernetes кластеры для анализа
- **Postgres (Pangolin) / Clickhouse** - базы данных для логирования метаданных и метрик

### Диаграммы

Подробная архитектура описана в [docs/DAAI_ARCHITECTURE.md](docs/DAAI_ARCHITECTURE.md)

Визуальные диаграммы доступны в [docs/diagrams/](docs/diagrams/):
- **Sequence диаграмма** (`sequence.puml`) - последовательность взаимодействий между компонентами
- **Process диаграмма** (`process.puml`) - детальный процесс обработки запроса с ветвлениями
- **Architecture диаграмма** (`architecture.puml`) - общая архитектура системы с компонентами

См. [docs/diagrams/README.md](docs/diagrams/README.md) для инструкций по просмотру и генерации диаграмм.

## Компоненты

### 1. Go HTTP Сервис (kagent-operator)

Принимает HTTP запросы от пользователя и формирует запросы к MCP Server.

**Основные функции:**
- Прием HTTP запросов на эндпоинт `/analyze`
- Валидация входных данных
- Формирование запросов к MCP Server
- Обработка и возврат результатов анализа

### 2. MCP Client

Клиент для взаимодействия с MCP Server по протоколу MCP (Model Context Protocol).

**Основные функции:**
- Подключение к MCP Server
- Вызов одного или нескольких tools
- Обработка ответов и ошибок

### 3. MCP Server (kagent-tools)

Сервер, который запускает tools для анализа Kubernetes кластера. Должен быть развернут отдельно на стороне клиента.

## Установка и запуск

### Требования

- Go 1.21 или выше
- Доступ к MCP Server (kagent-tools)

### Установка зависимостей

```bash
go mod download
```

### Сборка

```bash
go build -o kagent-operator .
```

### Запуск

```bash
# С переменными окружения по умолчанию
./kagent-operator

# С кастомными настройками
PORT=8080 MCP_SERVER_URL=http://localhost:3000/mcp ./kagent-operator
```

### Переменные окружения

- `PORT` - порт для HTTP сервера (по умолчанию: `8080`)
- `MCP_SERVER_URL` - URL MCP Server (по умолчанию: `http://localhost:3000/mcp`)

## API

### POST /analyze

Выполняет анализ кластера с использованием указанных tools.

**Request Body:**

```json
{
  "clusterId": "cluster-1",
  "toolNames": ["kubectl_get_pods", "kubectl_get_nodes"],
  "parameters": {
    "namespace": "default",
    "labelSelector": "app=myapp"
  }
}
```

**Response (200 OK):**

```json
{
  "clusterId": "cluster-1",
  "toolResults": {
    "kubectl_get_pods": {
      "pods": [...],
      "count": 5
    },
    "kubectl_get_nodes": {
      "nodes": [...],
      "count": 3
    }
  }
}
```

**Response (400 Bad Request):**

```json
{
  "error": "toolNames cannot be empty"
}
```

**Response (500 Internal Server Error):**

```json
{
  "error": "Failed to call MCP Server: connection refused"
}
```

### GET /health

Проверка здоровья сервиса.

**Response (200 OK):**

```json
{
  "status": "healthy",
  "service": "kagent-operator"
}
```

## Примеры использования

### Пример 1: Анализ подов в namespace

```bash
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "clusterId": "prod-cluster",
    "toolNames": ["kubectl_get_pods"],
    "parameters": {
      "namespace": "production"
    }
  }'
```

### Пример 2: Комплексный анализ кластера

```bash
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "clusterId": "prod-cluster",
    "toolNames": [
      "kubectl_get_nodes",
      "kubectl_get_pods",
      "get_events"
    ],
    "parameters": {
      "namespace": "default"
    }
  }'
```

### Пример 3: Использование с Prometheus tools

```bash
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "clusterId": "prod-cluster",
    "toolNames": ["prometheus_query"],
    "parameters": {
      "query": "up",
      "prometheus_url": "http://prometheus:9090"
    }
  }'
```

## Логика работы

### Поток выполнения запроса

1. **Получение HTTP запроса**
   - Пользователь отправляет POST запрос на `/analyze`
   - Сервис валидирует структуру запроса и обязательные поля

2. **Формирование запроса к MCP Server**
   - Для каждого tool из `toolNames` формируется MCP запрос
   - Параметры передаются в формате MCP Protocol

3. **Вызов MCP Server**
   - MCP Client отправляет HTTP запрос к MCP Server
   - MCP Server получает запрос и выполняет указанные tools
   - Tools анализируют состояние Kubernetes кластера

4. **Обработка результатов**
   - MCP Server возвращает результаты выполнения tools
   - Go сервис собирает результаты всех tools
   - Формируется единый ответ с результатами анализа

5. **Возврат ответа пользователю**
   - Результаты форматируются в JSON
   - Возвращается HTTP ответ с результатами анализа

### Обработка ошибок

- **Ошибки валидации**: возвращается HTTP 400 с описанием ошибки
- **Ошибки подключения к MCP Server**: возвращается HTTP 500
- **Ошибки выполнения tools**: результаты содержат поле `error` для каждого tool

### Параллельное выполнение

В текущей реализации tools выполняются последовательно. Для улучшения производительности можно реализовать параллельное выполнение tools.

## Интеграция с MCP Server

MCP Server должен быть развернут отдельно и доступен по указанному URL. 

### Формат запроса к MCP Server

```json
{
  "method": "tools/call",
  "params": {
    "name": "kubectl_get_pods",
    "arguments": {
      "namespace": "default"
    }
  }
}
```

### Формат ответа от MCP Server

```json
{
  "result": {
    "pods": [...],
    "count": 5
  },
  "error": null
}
```

## Развертывание

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o kagent-operator .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/kagent-operator .
EXPOSE 8080
CMD ["./kagent-operator"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kagent-operator
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kagent-operator
  template:
    metadata:
      labels:
        app: kagent-operator
    spec:
      containers:
      - name: kagent-operator
        image: kagent-operator:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: MCP_SERVER_URL
          value: "http://mcp-server:3000/mcp"
```

## Разработка

### Структура проекта

```
kagent-operator/
├── main.go           # HTTP сервер и обработчики
├── mcp/
│   └── client.go     # MCP клиент
├── go.mod
├── go.sum
└── README.md
```

### Добавление новых эндпоинтов

1. Добавить новый handler в `main.go`
2. Зарегистрировать route в `main()`
3. Обновить документацию

### Тестирование

```bash
# Запуск тестов
go test ./...

# Запуск с покрытием
go test -cover ./...
```

## Безопасность

- Валидация всех входных данных
- Ограничение размера тела запроса
- Таймауты для HTTP запросов к MCP Server
- Логирование всех операций

## Дальнейшее развитие

- [ ] Поддержка SSE транспорта для MCP
- [ ] Параллельное выполнение tools
- [ ] Кэширование результатов
- [ ] Метрики и мониторинг
- [ ] Аутентификация и авторизация
- [ ] Rate limiting
- [ ] WebSocket поддержка для real-time обновлений

## Лицензия

Apache-2.0
