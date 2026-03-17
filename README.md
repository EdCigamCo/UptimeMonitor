# UptimeMonitor

## Описание

Эта версия проекта добавляет REST API для управления мониторируемыми сайтами с использованием JSON для запросов и ответов.

## Что добавлено

### Новые компоненты

1. **DTO (Data Transfer Objects)** (`application/dto/dto.go`):
    - `CreateSiteRequest` - структура для запроса на создание сайта
    - `SiteResponse` - структура для ответа с информацией о сайте
    - `SiteListResponse` - структура для ответа со списком сайтов
    - `ErrorResponse` - структура для ответа с ошибкой

2. **Application Layer** (`application/uptime_monitor.go`):
    - `UptimeMonitor` - структура бизнес-логики
    - `CreateSite()` - создание сайта с валидацией URL
    - `GetAllSites()` - получение всех сайтов
    - `DeleteSite()` - удаление сайта по ID
    - `validateURL()` - валидация URL перед сохранением

3. **Mapper** (`presentation/mapper/mapper.go`):
    - `ToSiteResponse()` - преобразование модели в DTO для ответа
    - `ToSiteListResponse()` - преобразование списка моделей в DTO

### Изменения в существующих компонентах

- **presentation/controller.go**:
    - Добавлена функция `InitHandlers()` для инициализации с application layer
    - Добавлен `SitesHandler()` - обработчик для POST и GET `/api/sites`
    - Добавлен `SiteHandler()` - обработчик для DELETE `/api/sites/:id`
    - Добавлены вспомогательные функции `respondWithJSON()` и `respondWithError()`
    - Все новые эндпоинты возвращают JSON

- **cmd/main.go**:
    - Инициализация application layer с передачей БД
    - Инициализация presentation layer с передачей application
    - Регистрация новых API маршрутов

## Архитектура

Проект теперь следует принципам Clean Architecture с четким разделением слоев:

```
presentation (HTTP handlers)
    ↓
application (business logic)
    ↓
repository (data access)
    ↓
database
```

## API Эндпоинты

### POST /api/sites
Создает новый сайт для мониторинга.

**Request Body:**
```json
{
  "url": "https://example.com"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "url": "https://example.com",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid URL: URL must include scheme (http:// or https://)"
}
```

### GET /api/sites
Получает список всех мониторируемых сайтов.

**Response (200 OK):**
```json
{
  "sites": [
    {
      "id": 1,
      "url": "https://example.com",
      "created_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "url": "https://google.com",
      "created_at": "2024-01-15T11:00:00Z"
    }
  ]
}
```

### DELETE /api/sites/:id
Удаляет сайт по ID.

**Response (204 No Content):** - успешное удаление

**Response (404 Not Found):**
```json
{
  "error": "Site with id 123 not found"
}
```

## Валидация

При создании сайта выполняется валидация URL:
- URL не может быть пустым
- URL должен иметь правильный формат
- URL должен содержать схему (http:// или https://)
- URL должен содержать хост

## Обработка ошибок

API возвращает соответствующие HTTP коды статуса:
- `200 OK` - успешный запрос
- `201 Created` - успешное создание ресурса
- `204 No Content` - успешное удаление
- `400 Bad Request` - неверный запрос (валидация, парсинг)
- `404 Not Found` - ресурс не найден
- `500 Internal Server Error` - внутренняя ошибка сервера

Все ошибки возвращаются в формате JSON:
```json
{
  "error": "Описание ошибки"
}
```

## Старые эндпоинты

Все эндпоинты из предыдущих версий остаются доступными:
- `GET /health` - проверка работоспособности сервера (текст)
- `GET /info` - информация о сервере (текст)
- `GET /check?url=<website_url>` - проверка доступности сайта (текст)

## Как запустить

1. Установите зависимости:
```bash
go mod tidy
```

2. Запустите сервер:
```bash
go run ./cmd/main.go
```

3. Протестируйте API:

**Создать сайт:**
```bash
curl -X POST http://localhost:8080/api/sites \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

**Получить все сайты:**
```bash
curl http://localhost:8080/api/sites
```

**Удалить сайт:**
```bash
curl -X DELETE http://localhost:8080/api/sites/1
```

## Следующие шаги

В следующей версии будет добавлен фоновый воркер, который будет периодически проверять доступность всех сохраненных сайтов и сохранять результаты в базу данных.
