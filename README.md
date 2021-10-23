# Economic Calendar API & Parser

## Overview
The simple project impelementation of economic calendar news service. Contains following parts:
1. Database schema, data seeding and roles creation scripts  [/src/db](https://github.com/denis-gudim/economic-calendar/tree/main/src/db)
2. Loader services implementation on Golang [/src/loader](https://github.com/denis-gudim/economic-calendar/tree/main/src/loader)
2. REST api services implementation on Golang [/src/api](https://github.com/denis-gudim/economic-calendar/tree/main/src/api)

## Running Locally
Project contains docker compose script for local starting and test usages. Use following command for starting:
```bash
sudo docker-compose up --build
```

## Healthchecks & Metrics
Project contains HTTP healthcheck and prometeus exporter API.
**API service:**
```bash
curl -f 'http://localhost:8080/healtz'
curl -f 'http://localhost:8080/metrics
```
**Loader service:**
```bash
curl -f 'http://localhost:8081/healtz'
```

## Swagger
Use following link for swagger UI:
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html "http://localhost:8080/swagger/index.html")