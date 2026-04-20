Dental Clinic Backend
Требования
Docker version 27.3.1, build ce12230
Docker Compose version v2.30.3-desktop.1

Создайте файл .env в корне проекта:
APP_PORT=8080
DB_DSN=postgres://postgres:1234@db:5432/dental_clinic?sslmode=disable
JWT_SECRET=supersecretkey123
SMTP_USER=dautovalisher33@gmail.com
SMTP_PASS=ppzv hnga kzhl xdlz
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587


Собрать и запустить контейнеры:

docker-compose -f docker/docker-compose.yml up --build

После запуска приложение будет доступно по адресу:

http://localhost:8080

Остановить и удалить контейнеры:

docker-compose -f docker/docker-compose.yml down