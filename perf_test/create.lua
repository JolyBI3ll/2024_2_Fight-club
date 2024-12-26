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
