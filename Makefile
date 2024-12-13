.PHONY: build build-migrator build-ads build-auth build-city build-webapp run-migrator run-ads run-auth run-city run-webapp

build: build-migrator build-ads build-auth build-city build-webapp

build-migrator:
	go build -o bin/migrator ./cmd/migrator/

build-ads:
	go build -o bin/ads_service ./microservices/ads_service/cmd/main.go

build-auth:
	go build -o bin/auth_service ./microservices/auth_service/cmd/main.go

build-city:
	go build -o bin/city_service ./microservices/citiy_service/cmd/main.go

build-webapp:
	go build -o bin/webapp ./cmd/webapp/

run-migrator: build-migrator
	./bin/migrator

run-ads: build-ads
	./bin/ads_service

run-auth: build-auth
	./bin/auth_service

run-city: build-city
	./bin/city_service

run-webapp: build-webapp
	./bin/webapp

run: run-migrator run-ads run-auth run-city run-webapp