# UptimeMonitor

## Описание

Эта версия проекта добавляет API эндпоинт для получения истории проверок сайта, а также расширяет список сайтов данными о последней проверке.

## Что добавлено

### Новые компоненты

1. **Расширенные DTO** (`application/dto/dto.go`):
   - `SiteResponse` дополнен полями:
      - `status` - последний статус сайта (`up` или `down`)
      - `response_time` - время ответа последней проверки в миллисекундах
      - `last_checked` - время последней проверки в формате RFC3339
   - `CheckResponse` - структура одной проверки
   - `SiteHistoryResponse` - структура ответа с историей проверок сайта

2. **Application Layer** (`application/uptime_monitor.go`):
   - `GetAllSitesWithStatus()` - получение всех сайтов с их последней проверкой
   - `GetSiteHistory(siteID, limit)` - получение истории проверок конкретного сайта

3. **Mapper** (`presentation/mapper/mapper.go`):
   - `ToCheckResponse()` - преобразование `model.Check` в `dto.CheckResponse`
   - `ToSiteHistoryResponse()` - формирование ответа истории проверок
   - `ToSiteResponseWithCheck()` - формирование сайта с данными проверки
   - `ToSiteListResponseWithChecks()` - формирование списка сайтов с последними статусами

### Изменения в существующих компонентах

- **presentation/controller.go**:
   - `ListSitesHandler()` теперь возвращает сайты с полями последней проверки
   - Добавлен `GetSiteHistoryHandler()` для `GET /api/sites/:id/history`
   - Добавлена поддержка query параметра `limit` для ограничения размера истории

- **cmd/main.go**:
   - Зарегистрирован новый эндпоинт `GET /api/sites/:id/history`
   - Обновлено логирование доступных API маршрутов

## API эндпоинты

### GET /api/sites

Возвращает список всех сайтов с данными о последней проверке (если проверки есть).

**Пример ответа (200 OK):**
```json
{
  "sites": [
    {
      "id": 1,
      "url": "https://example.com",
      "created_at": "2026-03-20T10:00:00Z",
      "status": "up",
      "response_time": 142,
      "last_checked": "2026-03-20T10:05:00Z"
    },
    {
      "id": 2,
      "url": "https://offline.example",
      "created_at": "2026-03-20T10:01:00Z"
    }
  ]
}
```

### GET /api/sites/:id/history

Возвращает историю проверок для конкретного сайта.

**Query параметры:**
- `limit` (опционально) - максимальное количество проверок (по умолчанию 50)

**Пример ответа (200 OK):**
```json
{
  "site_id": 1,
  "url": "https://example.com",
  "checks": [
    {
      "id": 12,
      "status": "up",
      "response_time": 130,
      "checked_at": "2026-03-20T10:10:00Z"
    },
    {
      "id": 11,
      "status": "down",
      "response_time": 0,
      "checked_at": "2026-03-20T09:40:00Z"
    }
  ]
}
```

**Пример ошибки (404 Not Found):**
```json
{
  "error": "Site with id 999 not found"
}
```

## Как запустить

1. Установите зависимости:
```bash
go mod tidy
```

2. Запустите сервер:
```bash
go run ./cmd/main.go
```

## Как проверить

1. Получить сайты со статусами:
```bash
curl http://localhost:8080/api/sites
```

2. Получить историю проверок:
```bash
curl http://localhost:8080/api/sites/1/history
```

3. Получить историю с лимитом:
```bash
curl "http://localhost:8080/api/sites/1/history?limit=10"
```

## Старые эндпоинты

Сохраняются все эндпоинты из предыдущих шагов:
- `GET /health`
- `GET /info`
- `GET /check?url=<website_url>`
- `POST /api/site`
- `DELETE /api/site/:id`

## Следующие шаги

В следующей версии будет добавлена интеграция с фронтендом и поддержка CORS для браузерных запросов.
