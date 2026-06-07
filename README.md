# Dental Clinic Backend

## Requirements

* Docker
* Docker Compose


```bash
Docker version 27.3.1, build ce12230
Docker Compose version v2.30.3-desktop.1
```

---

## Environment

Создай файл `.env` в корне проекта:

```env
APP_PORT=
DB_DSN=
JWT_SECRET=
SMTP_USER=
SMTP_PASS=
SMTP_HOST=
SMTP_PORT=
OPENAI_API_KEY=
```

---

## Run Project

Сборка и запуск контейнеров:

```bash
docker-compose -f docker/docker-compose.yml up --build
```

Приложение будет доступно:

```
http://localhost:8080
```

---

## Stop Project

Остановка и удаление контейнеров:

```bash
docker-compose -f docker/docker-compose.yml down
```

---
