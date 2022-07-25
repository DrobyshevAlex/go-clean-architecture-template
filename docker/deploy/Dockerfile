FROM golang:1.18.1-buster as builder

WORKDIR /app

COPY ./src /app

RUN go mod download && go build -o bin/api-service ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin /app

CMD [ "./api-service" ]
