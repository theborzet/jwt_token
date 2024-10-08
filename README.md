# Описание проекта

JWT_token - это проект на Go, который использует PostgreSQL в качестве базы данных и Docker для контейнеризации. Проект предоставляет REST API для работы с аутентификацией, позволяя генерировать новые токены для пользователя и обновлять их.

## Стек технологий

- Go
- PostgreSQL
- JWT
- Docker
- Docker Compose

## Установка и запуск
1. Клонируйте репозиторий:
   git clone "https://github.com/theborzet/jwt_token"

2. Измените файл .env.example:

    Введите конфигурацию для БД и подписи для JWT токенов в этих файлах.

    Уберите ".example" из названия файла
3. При желании можете изменить ttl'ы токенов в файле configs/config.yaml:

### Запуск через Docker

1. Запустите Docker Compose:
   docker-compose up --build

2. Проверьте, что контейнеры запущены:
   Убедитесь, что контейнеры app и db запущены и работают корректно.

### Локальный запуск

1. Установка зависимостей
    go get ./...

2. Запуск приложения
    go run main.go или make build(сгенерирует exe-шник) make run(запустит этот exe-шник)

3. Проверка приложения
    Откройте веб-браузер и перейдите на http://localhost:8080 (или другой порт, указанный в файле .env)
