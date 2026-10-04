# Ручная проверка API через curl

Проверено 4 октября 2026 года, Go 1.25.4, Windows. Сервер запущен с пустым хранилищем. Ниже — фактические ответы curl; время в HTTP-заголовках указано в GMT. Тела запросов передавались через UTF-8 файл request.json.

## Список задач

```sh
curl.exe -i -X GET http://localhost:8080/tasks
```

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:25 GMT
Content-Length: 3

[]
```

## Неверный метод коллекции

```sh
curl.exe -i -X DELETE http://localhost:8080/tasks
```

```http
HTTP/1.1 405 Method Not Allowed
Allow: GET, POST
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:25 GMT
Content-Length: 31

{"error":"method not allowed"}
```

## Создание

Содержимое `request.json`:
```json
{"title":"Learn Go","done":false}
```

```sh
curl.exe -i -X POST http://localhost:8080/tasks -H "Content-Type: application/json" --data-binary @request.json
```

```http
HTTP/1.1 201 Created
Content-Type: application/json
Location: /tasks/1
Date: Sun, 04 Oct 2026 07:53:25 GMT
Content-Length: 77

{"id":1,"title":"Learn Go","done":false,"created_at":"2026-10-04T07:53:25Z"}
```

## Пустой заголовок

Содержимое `request.json`:
```json
{"title":" "}
```

```sh
curl.exe -i -X POST http://localhost:8080/tasks -H "Content-Type: application/json" --data-binary @request.json
```

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:25 GMT
Content-Length: 52

{"error":"title is required and must not be blank"}
```

## Получение

```sh
curl.exe -i -X GET http://localhost:8080/tasks/1
```

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:25 GMT
Content-Length: 77

{"id":1,"title":"Learn Go","done":false,"created_at":"2026-10-04T07:53:25Z"}
```

## Несуществующая задача

```sh
curl.exe -i -X GET http://localhost:8080/tasks/999
```

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:25 GMT
Content-Length: 27

{"error":"task not found"}
```

## Обновление

Содержимое `request.json`:
```json
{"title":"Learn Go REST","done":true}
```

```sh
curl.exe -i -X PUT http://localhost:8080/tasks/1 -H "Content-Type: application/json" --data-binary @request.json
```

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:25 GMT
Content-Length: 81

{"id":1,"title":"Learn Go REST","done":true,"created_at":"2026-10-04T07:53:25Z"}
```

## Обновление отсутствующей задачи

Содержимое `request.json`:
```json
{"title":"Missing","done":false}
```

```sh
curl.exe -i -X PUT http://localhost:8080/tasks/999 -H "Content-Type: application/json" --data-binary @request.json
```

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:25 GMT
Content-Length: 27

{"error":"task not found"}
```

## Удаление

```sh
curl.exe -i -X DELETE http://localhost:8080/tasks/1
```

```http
HTTP/1.1 204 No Content
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:25 GMT

```

## Повторное удаление

```sh
curl.exe -i -X DELETE http://localhost:8080/tasks/1
```

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:25 GMT
Content-Length: 27

{"error":"task not found"}
```

## Проверка здоровья

```sh
curl.exe -i -X GET http://localhost:8080/health
```

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:26 GMT
Content-Length: 16

{"status":"ok"}
```

## Неверный метод health

```sh
curl.exe -i -X POST http://localhost:8080/health
```

```http
HTTP/1.1 405 Method Not Allowed
Allow: GET
Content-Type: application/json
Date: Sun, 04 Oct 2026 07:53:26 GMT
Content-Length: 31

{"error":"method not allowed"}
```
