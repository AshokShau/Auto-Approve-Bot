FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache tzdata

COPY --from=builder /app/main .

CMD ["./main"]
