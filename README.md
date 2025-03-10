# Криптовалютный API

Этот проект предоставляет API для получения исторических данных о криптовалютах с биржи Bybit. Данные доступны в форматах JSON и CSV.

---


# Запуск приложения

Для запуска приложения выполните следующие команды:

```bash
cd src/cmd
go run main.go
```

# Запросы к API

## Получение списка доступных криптовалют

- **URL:** `http://localhost:8080/symbols`
- **Метод:** `[GET]`
- **Описание:** Получение списка доступных криптовалют.

## Получение информации по криптовалюте в формате JSON

- **URL:** `http://localhost:8080/klines`
- **Метод:** `[GET]`
- **Описание:** Получение информации по выбранной криптовалюте.
- **Пример:** http://localhost:8080/klines?symbol=BTCUSDT&interval=30&days=7

### Параметры запроса

| Параметр | Описание                                                                                                                                 |
|-----------|------------------------------------------------------------------------------------------------------------------------------------------|
| `symbol`  | Валюта (например, `BTCUSDT`)                                                                                                             |
| `interval`| Отрезок времени, на который будет делиться выбранный промежуток в минутах, <br/>поддерживает значения 1, 3, 5, 15, 30, 60, 120, 240, 720 |
| `days`    | Параметр для выбора промежутка. Работает как текущая дата минус `days` (например, `7`)                                                   |

# Ответ сервера

При успешном запросе сервер возвращает ответ в формате JSON:

```json
[
  {
    "open_time": "2025-03-01 00:00:00",
    "open_price": "85716.08",
    "high_price": "88713.02",
    "low_price": "85047.52",
    "close_price": "87446.78",
    "volume_base": "1886.71598",
    "volume_quote": "164072437.87526441"
  }
]
```

### Описание полей ответа

| Поле          | Описание |
|---------------|----------|
| `open_time`   | Время свечи |
| `open_price`  | Цена при открытии свечи |
| `high_price`  | Высшая цена за интервал |
| `low_price`   | Низшая цена за интервал |
| `close_price` | Цена при закрытии интервала |
| `volume_base` | Объем торгов в базовой валюте |
| `volume_quote`| Объем торгов в котировочной валюте |

## Получение информации по криптовалюте в формате CSV

- **URL:** `http://localhost:8080/klines/csv`
- **Метод:** `[GET]`
- **Описание:** Получение информации по выбранной криптовалюте в формате CSV.
- **Пример:** http://localhost:8080/klines/csv?symbol=BTCUSDT&interval=30&days=7

### Параметры запроса

| Параметр | Описание |
|-----------|----------|
| `symbol`  | Валюта (например, `BTCUSDT`) |
| `interval`| Отрезок времени, на который будет делиться выбранный промежуток в минутах, <br/>поддерживает значения 1, 3, 5, 15, 30, 60, 120, 240, 720  |
| `days`    | Параметр для выбора промежутка. Работает как текущая дата минус `days` (например, `7`) |

# Ответ сервера

При успешном запросе сервер возвращает ссылку для скачивая csv файла:

```json
{
    "file_url": "http://localhost:8080/static/BTCUSDT_30_7d.csv"
}
```