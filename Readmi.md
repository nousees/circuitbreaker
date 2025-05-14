## Структура проекта

client.go - Клиентский сервис, который взаимодействует с двумя серверами
apiclient.go - Абстракция ApiClient с реализацией Circuit Breaker
server.go - Первый сервер на порту 8080 с 50% вероятностью ошибки
server2.go - Второй сервер на порту 8081 с 25% вероятностью ошибки 

## Установка

1. Склонируйте репозиторий:
   ```bash
   git clone https://github.com/nousees/circuitbreaker.git
   cd circuitbreaker
   ```

### Запуск через Docker


1. Запустите все сервисы:
   ```bash
   docker-compose up --build
   ```

2. Для остановки контейнеров:
   Нажмите `Ctrl+C` в терминале, затем выполните:
   ```bash
   docker-compose down
   ```

## Пример вывода

```
2025/05/14 09:13:24 Attempt 1 failed: server error: 500
2025/05/14 09:13:24 Attempt 2 failed: server error: 500
2025/05/14 09:13:25 Attempt 3 failed: server error: 500
2025/05/14 09:13:25 Client 1 error: all retries failed: server error: 500
2025/05/14 09:13:37 Client 2 response: Success response from server 2
```

Клиент выполняет до 3 попыток при сбоях.

Circuit Breaker переходит в состояние "открыт" после трёх последовательных ошибок, блокируя дальнейшие запросы до восстановления.