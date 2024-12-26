# Домашнее задание №3

Основная сущность нашего приложения - объявление. Мы провели нагрузочное тестирование: создали 100000 объявлений, затем прочитали 100000 объявлений, потом оценили результаты. Для проведения нагрузочного тестирования мы выбрали https://github.com/wg/wrk.

## Создание объявлений

Для генерации объявлений мы написали скрипт [create.lua](./create.lua):
```lua
wrk.method = "POST"
wrk.headers["X-CSRF-Token"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWQiOiJnNm8zZm5qaSs1aS9SdTlhellybkc0anN2R01OeFJ3QUhMeExEQXI3WG9ZPSIsImV4cCI6MTczNTI1MjExOSwiaWF0IjoxNzM1MTY1NzE5fQ.tmdTaJm_krJ7I2UNa1kRsA_LDbwa9os5h1noFPwNYv4"
wrk.headers["Cookie"] = "session_id=g6o3fnji+5i/Ru9azYrnG4jsvGMNxRwAHLxLDAr7XoY=; Path=/; SameSite=Strict; Domain=127.0.0.1"

local boundary = "----WebKitFormBoundary7MA4YWxkTrZu0gW"
wrk.headers["Content-Type"] = "multipart/form-data; boundary=" .. boundary

local function read_file(file_path)
    local file = io.open(file_path, "rb")  -- Открываем файл в бинарном режиме
    if not file then
        error("Failed to open file: " .. file_path)
    end
    local content = file:read("*all")  -- Считываем содержимое файла
    file:close()
    return content
end

local counter = 0
local image_data = read_file("../img.png")  -- Укажите путь к вашему файлу изображения

request = function()
    counter = counter + 1

    local cityName = "Moscow"
    local description = "description of ads " .. counter
    local address = "address of ads " .. counter
    local roomsNumber = 2
    local dateFrom = "2024-11-10T00:00:00Z"
    local dateTo = "2024-11-10T00:00:00Z"

    local metadata = string.format(
        '{"cityName": "%s", "description": "%s", "address": "%s", "roomsNumber": %s, "dateFrom": "%s", "dateTo": "%s"}',
        cityName, description, address, roomsNumber, dateFrom, dateTo
    )

    local body =
        "--" .. boundary .. "\r\n" ..
        'Content-Disposition: form-data; name="metadata"\r\n\r\n' ..
        metadata .. "\r\n" ..
        "--" .. boundary .. "\r\n" ..
        'Content-Disposition: form-data; name="images"; filename="image.jpg"\r\n' ..
        "Content-Type: image/jpeg\r\n\r\n" ..
        image_data .. "\r\n" ..
        "--" .. boundary .. "--\r\n"

    return wrk.format(nil, "/api/housing", nil, body)
end
```

Запускали вот такой командой:
```bash
wrk -t4 -c100 -d300s --script=create.lua --latency http://127.0.0.1:8008
```

Вывод команды:
![img.png](images/img_1.png)

## Чтение объявлений

Для чтения объявлений мы написали скрипт [read.lua](./read.lua):
```lua
wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"


request = function()
    return wrk.format(nil, "/api/housing/ca1c71d5-a526-4799-8d3f-aa023d0bb041", nil, nil)
end
```

Запускали вот такой командой:
```bash
wrk -t4 -c100 -d300s --script=read.lua --latency http://127.0.0.1:8008
```

Вывод команды:
![img_2.png](images/img_2.png)

## Анализ результатов

При создании объявлений получили около 283 rps - это небольшой показатель, но и не плохой.

При чтении объявлений получили также около 150 rps.

## Оптимизации

Можно оптимизировать запрос на получение объявлений: можно денормализовать схему данных, избавившись от таблицы доступных дат. Вместо таблицы доступных дат можно будет сохранять эти данные внутри каждого объявления. Это поможет избавиться от JOIN в запросе. Также можно поступить с остальными JOIN.

