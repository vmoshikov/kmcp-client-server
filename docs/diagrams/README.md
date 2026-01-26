# Диаграммы архитектуры DAAI

Этот каталог содержит PlantUML диаграммы, описывающие архитектуру системы DAAI.

## Файлы

### 1. `sequence.puml` - Sequence диаграмма

Показывает последовательность взаимодействий между компонентами системы при обработке запроса от пользователя.

**Компоненты:**
- Пользователь
- DAAI (Master Service)
- Postgres (Pangolin) / Clickhouse
- KaaS (REST/gRPC)
- MCP Client (внутри DAAI)
- MCP Server (kmcp) в Parent Cluster
- kagent tools
- Parent Cluster Kubernetes
- Children Cluster Kubernetes

**Поток:**
1. Пользователь отправляет запрос в DAAI
2. DAAI логирует запрос
3. DAAI выполняет авторизацию и бизнес-логику
4. DAAI получает информацию о кластере от KaaS
5. DAAI через MCP Client обращается к MCP Server
6. MCP Server запускает kagent tools
7. Tools анализируют Parent или Children Cluster
8. Результаты возвращаются обратно пользователю

### 2. `process.puml` - Диаграмма процесса

Показывает детальный процесс обработки запроса с ветвлениями и условиями.

**Основные этапы:**
- Логирование запроса
- Авторизация и валидация
- Получение информации о кластере от KaaS
- Формирование MCP запроса
- Выполнение tools на MCP Server
- Анализ кластера (Parent или Children)
- Обработка результатов
- Возврат ответа

### 3. `architecture.puml` - Архитектурная диаграмма

Показывает общую архитектуру системы с компонентами и их связями.

**Основные пакеты:**
- External (Пользователь)
- DAAI Service (Master Service, MCP Client, Business Logic, Logging)
- Data Storage (Postgres, Clickhouse)
- KaaS Service
- Parent Cluster Kubernetes (MCP Server, kagent tools)
- Children Clusters Kubernetes

## Как использовать

### Онлайн просмотр

1. Перейдите на [PlantUML Online Server](http://www.plantuml.com/plantuml/uml/)
2. Скопируйте содержимое нужного `.puml` файла
3. Вставьте в редактор
4. Диаграмма будет автоматически отрендерена

### Локальная установка

```bash
# Установка PlantUML
# macOS
brew install plantuml

# Linux
sudo apt-get install plantuml

# Или через Java
java -jar plantuml.jar diagram.puml
```

### Генерация изображений

```bash
# PNG
plantuml -tpng sequence.puml

# SVG
plantuml -tsvg sequence.puml

# PDF
plantuml -tpdf sequence.puml
```

### Интеграция с VS Code

Установите расширение "PlantUML" для VS Code:
- Автоматический preview
- Экспорт в различные форматы
- Подсветка синтаксиса

## Обновление диаграмм

При изменении архитектуры системы обновите соответствующие диаграммы:

1. **Sequence диаграмма** - при изменении последовательности взаимодействий
2. **Process диаграмма** - при изменении бизнес-логики или процесса обработки
3. **Architecture диаграмма** - при добавлении/удалении компонентов

## Примечания

- Все диаграммы используют тему `plain` для лучшей читаемости
- Компоненты сгруппированы по логическим пакетам
- Связи показывают направление потока данных
- Примечания (notes) добавлены для важных моментов
