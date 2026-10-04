# Tasks API

Учебный REST API на стандартной библиотеке Go

## Запуск


```sh
go run ./cmd/server
```

Адрес: `http://localhost:8080`. 

## API

| Метод | Путь | Назначение | Успех |
|---|---|---|---|
| GET | /tasks | Список задач (пустой список — `[]`) | 200 |
| POST | /tasks | Создать задачу | 201 + Location |
| GET | /tasks/{id} | Получить задачу | 200 |
| PUT | /tasks/{id} | Заменить title и done | 200 |
| DELETE | /tasks/{id} | Удалить задачу | 204 |
| GET | /health | Проверить доступность | 200 |


## Примеры curl


```sh
curl -i http://localhost:8080/tasks
curl -i -X POST http://localhost:8080/tasks -H 'Content-Type: application/json' -d '{"title":"Learn Go","done":false}'
curl -i http://localhost:8080/tasks/1
curl -i -X PUT http://localhost:8080/tasks/1 -H 'Content-Type: application/json' -d '{"title":"Learn Go REST","done":true}'
curl -i -X DELETE http://localhost:8080/tasks/1
curl -i http://localhost:8080/health
```

Пример тела ответа POST (ID и время зависят от состояния сервера):

```json
{"id":1,"title":"Learn Go","done":false,"created_at":"2026-10-04T09:00:00Z"}
```

Ошибочные сценарии:

```sh
curl -i -X DELETE http://localhost:8080/tasks
curl -i -X POST http://localhost:8080/tasks -H 'Content-Type: application/json' -d '{"title":" "}'
curl -i http://localhost:8080/tasks/999
curl -i -X PUT http://localhost:8080/tasks/999 -H 'Content-Type: application/json' -d '{"title":"Missing","done":false}'
curl -i -X DELETE http://localhost:8080/tasks/999
curl -i -X POST http://localhost:8080/health
```



## Проверка

```sh
go test ./...
go test -race ./...
go vet ./...
```


## Структура

```text
cmd/server/main.go                 запуск сервера и тайм-ауты
internal/models/task.go            модель Task
internal/storage/storage.go        интерфейс Storage
internal/storage/memory/           потокобезопасное хранилище и тесты
internal/handlers/                 маршруты, JSON, валидация и тесты
internal/http/middleware.go        логирование
```

