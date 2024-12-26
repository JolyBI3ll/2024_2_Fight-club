wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"


request = function()
    return wrk.format(nil, "/api/housing/ca1c71d5-a526-4799-8d3f-aa023d0bb041", nil, nil)
end