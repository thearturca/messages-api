# Тестовое задание

Этот репозиторий - результат выполнения тестового задания по созданию сервиса для сообщений.

# Запуск

Для запуска в контейнере необходим `.env` файл для конфигурации. Пример файла [`.env.example`](./configs/.env.example).

## Запуск локально в контейнере
```shell
docker compose -f docker-compose.yml -f docker-compose.local.yml up --build

```

# Эндпоинты

- `POST /messages` - создать новое сообщение. Принимается JSON-объект с полем `text`. 
Пример запроса:
```json
{
    "text": "Hello, World!"
}
```

Пример ответа:
```json
{
    "id": "fbc4136e-d7e9-4e9f-9b09-c7008a24244b",
    "text": "Hello, World!",
    "isProcessed": false,
    "processedAt": null,
    "createdAt": "2024-07-30T10:39:43.997814Z"
}

```

- `GET /messages/{id}` - получить сообщение по его id. Пример:
```shell
curl -X GET http://localhost:3000/messages/fbc4136e-d7e9-4e9f-9b09-c7008a24244b
```
Пример ответа:
```json
{
    "id": "fbc4136e-d7e9-4e9f-9b09-c7008a24244b",
    "text": "Hello, World!",
    "isProcessed": false,
    "processedAt": null,
    "createdAt": "2024-07-30T10:39:43.997814Z"
}
```


- `GET /stats` - получить статистику по сообщениям. Эндпоинт поддерживает фильтрацию по дате. Для фильтрации нужно указать `from` и `to` параметры. `from` и `to` являются необязательными параметрами. Пример:
```shell
curl -X GET http://localhost:3000/stats?from=2024-07-30T00:00:00Z&to=2024-07-30T10:39:43Z
```

```shell
curl -X GET http://localhost:3000/stats?from=2024-07-30T00:00:00Z
```

```shell
curl -X GET http://localhost:3000/stats?to=2024-07-30T10:39:43Z
```

Пример ответа:
```json
{
    "totalMessages": 7,
    "processedMessages": 5
}
```
