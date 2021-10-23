FROM golang:1.17-alpine3.14 AS build-ldr

WORKDIR /src/

COPY . .

RUN go get -d -v ./...
RUN go build -o ./out/calendar-loader .

FROM golang:1.17-alpine3.14

WORKDIR /opt/app/

RUN apk --update --no-cache add curl

COPY --from=build-ldr /src/out/calendar-loader calendar-loader
COPY --from=build-ldr /src/config.env config.env

EXPOSE 8080

CMD ["/opt/app/calendar-loader"]