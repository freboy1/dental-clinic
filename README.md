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
APP_PORT=8080
DB_DSN=postgres://postgres:1234@localhost:5432/dental_clinic?sslmode=disable
JWT_SECRET=supersecretkey123
SMTP_USER=dautovalisher33@gmail.com
SMTP_PASS=ppzv hnga kzhl xdlz
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
OPENAI_API_KEY=sk-22c0691832294091bfbd9265546844c1
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
