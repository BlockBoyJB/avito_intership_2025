# avito_intership_2025

### Prerequisites

- Docker, Docker Compose
- or Golang 1.25 + postgresql

### Getting started

* Добавить репозиторий к себе
* Создать .env файл в директории с проектом и заполнить информацией из .env.example

### Usage

Запустить сервис можно с помощью `make compose-up` (или `docker-compose up -d --build`)
или `make run` (при наличии go1.25 и локально развернутого postgresql)

Помимо основного задания был реализован эндпоинт статистки, который возвращает количество участников в командах, 
количество PR по авторам, ревьюеров с количеством PR для проверки  

Примеры запросов

`GET /stats?filter=author`

```json
[
  {
    "author_id": "u2",
    "pr_count": 1
  },
  {
    "author_id": "u3",
    "pr_count": 1
  }
]
```

`GET /stats?filter=team`

```json
[
  {
    "team_name": "payments",
    "members": 5
  },
  {
    "team_name": "billing",
    "members": 2
  }
]
```

`GET /stats?filter=reviewers`

```json
[
  {
    "reviewer_id": "u4",
    "assignments": 2
  },
  {
    "reviewer_id": "u5",
    "assignments": 2
  }
]
```