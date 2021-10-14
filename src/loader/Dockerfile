FROM golang:1.17-alpine3.14

WORKDIR /go/src/calendar-loader
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["calendar-loader"]