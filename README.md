### Запуск контейнера с БД

```
make compose-up
```

### Запуск миграций

```
make goose install
make goose-up
```

### Запуск приложения

```
make run
```

API будет доступен на `localhost:8080/v1/auth`, swagger с описанием путей на `localhost:8080/swagger/index.html`.

### Запуск тестов

```
make test
```
