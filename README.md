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
## Swagger
Use following link for swagger UI:
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html "http://localhost:8080/swagger/index.html")
