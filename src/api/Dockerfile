FROM golang:1.17-alpine3.14 AS build-api

WORKDIR /src/

COPY . .

RUN go get -d -v ./...
RUN go build -o ./out/calendar-api .

FROM golang:1.17-alpine3.14

WORKDIR /opt/app/

COPY --from=build-api /src/out/calendar-api calendar-api
COPY --from=build-api /src/config.env config.env

EXPOSE 8080

CMD ["/opt/app/calendar-api"]