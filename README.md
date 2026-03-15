# UptimeMonitor

## Структура проекта

```
uptime_monitor/
├── cmd/
│   └── main.go                    # Точка входа приложения
├── presentation/
│   └── controller.go             # HTTP контроллеры
├── infrastructure/
│   └── worker/
│       └── worker.go             # Функция проверки доступности сайтов
├── application/                   # (пустая, будет использоваться позже)
├── model/                         # (пустая, будет использоваться позже)
├── migrations/                    # (пустая, будет использоваться позже)
├── go.mod
└── README.md
```

## Функционал

- Базовый HTTP сервер на порту 8080 (или из переменной окружения PORT)
- Эндпоинт `GET /health` - проверка работоспособности сервера (возвращает "OK")
- Эндпоинт `GET /info` - информация о сервере (возвращает текст с временем запуска)
- **Новое:** Эндпоинт `GET /check?url=<website_url>` - проверка доступности сайта (возвращает статус и время ответа)
- **Новое:** Функция `CheckSiteAvailability()` в `infrastructure/worker/worker.go` для проверки доступности сайтов

## Новая функциональность

### CheckSiteAvailability

Функция для проверки доступности веб-сайта:

```go
status, responseTime, err := worker.CheckSiteAvailability("https://example.com")
```

**Параметры:**
- `url string` - URL сайта для проверки

**Возвращает:**
- `status string` - "up" (код 200-299) или "down" (остальные коды или ошибки)
- `responseTime int64` - время ответа в миллисекундах
- `err error` - ошибка, если запрос не удался

**Особенности:**
- Таймаут запроса: 5 секунд
- Автоматическое закрытие тела ответа
- Обработка сетевых ошибок и таймаутов

## Запуск

```bash
go run cmd/main.go
```

## Тестирование

```bash
# Проверка health endpoint
curl http://localhost:8080/health

# Проверка info endpoint
curl http://localhost:8080/info

# Проверка доступности сайта
curl "http://localhost:8080/check?url=https://google.com"
curl "http://localhost:8080/check?url=https://example.com"
```

**Пример ответа от /check:**
```
URL: https://google.com
Status: up
Response time: 245 ms
```

## Пример использования CheckSiteAvailability

```go
package main

import (
    "fmt"
    "log"
    "uptime_monitor/infrastructure/worker"
)

func main() {
    status, responseTime, err := worker.CheckSiteAvailability("https://google.com")
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    
    fmt.Printf("Status: %s, Response time: %d ms\n", status, responseTime)
}
```

## Примечания

- На этом этапе используются только простые текстовые ответы
- Функция проверки доступности готова к использованию в фоновом воркере
- Используется стандартная библиотека Go для HTTP клиента
