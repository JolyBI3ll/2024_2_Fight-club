wrk.method = "POST"
wrk.headers["Content-Type"] = "multipart/form-data"
wrk.headers["X-CSRF-Token"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWQiOiJyNVgwanVuakJuTnVSVlc1T2hIQVJQckpnQkg0ZXZkc0hVbm9hWUYvNGVjPSIsImV4cCI6MTczNTI1MTIwMiwiaWF0IjoxNzM1MTY0ODAyfQ.cAIcFmFYlB3C5A8TEoaNMeyqMeW9Ai0K6n3J3ik9qXE"
wrk.headers["Cookie"] = "session_id=r5X0junjBnNuRVW5OhHARPrJgBH4evdsHUnoaYF/4ec=; Path=/; SameSite=Strict; Domain=localhost"

local counter = 0

request = function()
    counter = counter + 1

    local cityName = "Moscow"
    local description = "description of ads " + counter
    local address = "address of ads " + counter
    local roomsNumber = counter
    local dateFrom = "2024-11-10T00:00:00Z"
    local dateTo = "2024-11-10T00:00:00Z"
    local body = string.format('{"metadata": {"cityName": "%s", "description": "%s", "address": "%s", "roomsNumber": "%s", "dateFrom": "%s", "dateTo": "%s"}}', cityName, description, address, roomsNumber, dateFrom, dateTo)
    return wrk.format(nil, "/api/housing", nil, body)
end