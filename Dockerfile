FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY .env .env

RUN go build -o server .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/server .

COPY --from=builder /app/.env .env

EXPOSE 8080

CMD ["./server"]