# Быстрый старт

## Установка

```bash
# Клонировать репозиторий
git clone <repository-url>
cd kagent-operator

# Установить зависимости
go mod download
```

## Запуск

### Локальный запуск

```bash
# Сборка
go build -o kagent-operator .

# Запуск
./kagent-operator
```

Или с переменными окружения:

```bash
PORT=8080 MCP_SERVER_URL=http://localhost:3000/mcp ./kagent-operator
```

### Docker

```bash
# Сборка образа
docker build -t kagent-operator:latest .

# Запуск контейнера
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e MCP_SERVER_URL=http://host.docker.internal:3000/mcp \
  kagent-operator:latest
```

### Kubernetes

```bash
# Применить манифест
kubectl apply -f examples/kubernetes-deployment.yaml
```

## Проверка работы

### Health check

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:
```json
{
  "status": "healthy",
  "service": "kagent-operator"
}
```

### Тестовый запрос

```bash
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "clusterId": "test-cluster",
    "toolNames": ["kubectl_get_pods"],
    "parameters": {
      "namespace": "default"
    }
  }'
```

## Требования

- Go 1.21+
- Доступ к MCP Server (kagent-tools)
- MCP Server должен быть запущен и доступен по указанному URL

## Настройка MCP Server

Убедитесь, что MCP Server (kagent-tools) запущен и доступен:

```bash
# Пример запуска kagent-tools
./kagent-tools
# Сервер должен быть доступен на http://localhost:3000/mcp
```

## Troubleshooting

### Ошибка подключения к MCP Server

```
Failed to call MCP Server: connection refused
```

**Решение:** Проверьте, что MCP Server запущен и доступен по указанному URL.

### Ошибка валидации

```
toolNames cannot be empty
```

**Решение:** Убедитесь, что в запросе указан хотя бы один tool в массиве `toolNames`.

### Таймаут запроса

```
MCP Server returned status 500: timeout
```

**Решение:** Увеличьте таймаут в `mcp/client.go` или проверьте производительность MCP Server.
