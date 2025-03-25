# Take a Pill

Сервис для управления расписанием приема лекарств. Позволяет создавать расписания приема лекарств и получать информацию о ближайших приемах.

## Функциональность

- Создание расписания приема лекарств
- Получение деталей расписания
- Получение списка ближайших приемов
- Автоматическое распределение времени приема в течение дня
- Фильтрация прошедших приемов

## Требования

- Go 1.16 или выше
- gorilla/mux (для маршрутизации)

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/take-a-pill.git
cd take-a-pill
```

2. Установите зависимости:
```bash
go mod download
```

3. Запустите сервер:
```bash
go run main.go
```

Сервер будет доступен по адресу: http://localhost:8081

## API Endpoints

### Создание расписания
```http
POST /schedule
Content-Type: application/json

{
    "user_id": "string",
    "medicine_name": "string",
    "frequency": integer,
    "duration": integer
}
```

### Получение деталей расписания
```http
GET /schedule?user_id=string&schedule_id=uuid
```

### Получение списка ближайших приемов
```http
GET /next_takings?user_id=string
```

## Примеры использования

### Создание расписания
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"user_id":"test123","medicine_name":"Аспирин","frequency":3,"duration":7}' \
  http://localhost:8081/schedule
```

### Получение деталей расписания
```bash
curl "http://localhost:8081/schedule?user_id=test123&schedule_id=your-schedule-id"
```

### Получение списка ближайших приемов
```bash
curl "http://localhost:8081/next_takings?user_id=test123"
```

## OpenAPI документация

Полная документация API доступна в файле `openapi.yaml`. Вы можете использовать этот файл с инструментами вроде Swagger UI для просмотра и тестирования API.

## Тестирование

Для запуска тестов выполните:
```bash
go test -v
```

## Структура проекта

```
take-a-pill/
├── main.go           # Основной файл приложения
├── main_test.go      # Тесты
├── models/           # Модели данных
│   └── models.go
├── storage/          # Реализация хранилища
│   └── storage.go
├── openapi.yaml      # OpenAPI спецификация
└── README.md         # Документация
# Take a Pill

Простой сервис для создания расписания приема лекарств.

## Что умеет сервис

Сервис может создавать расписание приема лекарств. Для этого нужно указать:
- ID пользователя (обязательно)
- Название лекарства
- Сколько раз в день принимать
- На сколько дней назначен курс

## Как установить

1. Скачайте код:
```bash
git clone https://github.com/yourusername/take-a-pill.git
cd take-a-pill
```

2. Установите нужные пакеты:
```bash
go mod download
```

3. Запустите сервер:
```bash
go run main.go
```

## Как использовать

Чтобы создать расписание, отправьте POST запрос на `/schedule` с такими данными:
```json
{
    "user_id": "123",
    "medicine_name": "Аспирин",
    "frequency": 3,
    "duration": 7
}
```

Сервер ответит ID созданного расписания:
```json
{
    "schedule_id": "6f9bf3d1-a2c1-4992-b116-7c27ce448302"
}
```

Если что-то пошло не так, сервер вернет ошибку:
- Если не указан user_id: "не указан user_id"
- Если данные в неправильном формате: "Неверный формат данных"
- Если используется не POST запрос: "Нужно использовать POST метод"

## Из чего состоит проект

```
take-a-pill/
├── main.go           # Основной файл с кодом сервера
├── models/
│   └── models.go     # Структуры данных
├── storage/
│   └── storage.go    # Хранение данных в памяти
└── README.md         # Этот файл
```

## Как протестировать

Можно использовать curl для проверки:

```bash
# Создать расписание
curl -X POST -H "Content-Type: application/json" \
     -d '{"user_id": "123", "medicine_name": "Аспирин", "frequency": 3, "duration": 7}' \
     http://localhost/schedule

# Проверить ошибку с пустым user_id
curl -X POST -H "Content-Type: application/json" \
     -d '{"user_id": "", "medicine_name": "Аспирин", "frequency": 3, "duration": 7}' \
     http://localhost/schedule
```

## API

### Создание расписания

```http
POST /schedule
Content-Type: application/json

{
    "user_id": "123",
    "medicine_name": "Аспирин",
    "frequency": 3,
    "duration": 7
}
```

#### Ответ

```json
{
    "schedule_id": "6f9bf3d1-a2c1-4992-b116-7c27ce448302"
}
```

#### Ошибки

- `400 Bad Request` - неверный формат данных или пустой user_id
- `405 Method Not Allowed` - неверный метод запроса (только POST)

## Структура проекта

```
take-a-pill/
├── main.go           # Точка входа приложения
├── models/
│   └── models.go     # Определения структур данных
├── storage/
│   └── storage.go    # Реализация хранилища данных
└── README.md         # Документация проекта
```

## Тестирование

Для тестирования API можно использовать curl:

```bash
# Создание расписания
curl -X POST -H "Content-Type: application/json" \
     -d '{"user_id": "123", "medicine_name": "Аспирин", "frequency": 3, "duration": 7}' \
     http://localhost/schedule

# Проверка ошибок
curl -X POST -H "Content-Type: application/json" \
     -d '{"user_id": "", "medicine_name": "Аспирин", "frequency": 3, "duration": 7}' \
     http://localhost/schedule
```

## API Endpoints

### POST /schedule
Создание нового расписания приема лекарств.

Пример запроса:
```json
{
    "user_id": "123",
    "medicine_name": "Аспирин",
    "frequency": 3,
    "duration": 14
}
```

### GET /schedule?user_id=
Получение списка идентификаторов расписаний пользователя.

### GET /schedule?user_id=&schedule_id=
Получение детальной информации о конкретном расписании.

### GET /next_takings?user_id=
Получение списка следующих приемов лекарств в ближайший час.

## Особенности

- Прием лекарств только в дневное время (с 8:00 до 22:00)
- Время приема округляется до ближайших 15 минут
- Поддержка различных периодичностей приема (от одного раза в день до ежечасного)
- Возможность указать продолжительность курса лечения
