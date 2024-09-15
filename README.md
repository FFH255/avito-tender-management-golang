## Сбор и развертывание приложения

```
docker-compose up --build -d
```

После запуска, приложение будет доступно по `localhost:8080`

### Переменные окружения
```
SERVICE_ADDRESS=0.0.0.0:8080
SERVICE_READ_TIMEOUT=30
SERVICE_WRITE_TIMEOUT=30
SERVICE_IDLE_TIMEOUT=30

POSTGRES_DATABASE=tms
POSTGRES_USERNAME=admin
POSTGRES_PASSWORD=admin
POSTGRES_CONN=postgres://admin:admin@postgres:5432/tms
POSTGRES_AUTO_MIGRATE=true
POSTGRES_MIGRATION=migration/001_initial_migrations.sql

```

