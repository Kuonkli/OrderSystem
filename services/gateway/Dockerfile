FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем ВСЕ из корня проекта в контейнер
COPY . .

# Переходим в папку gateway и собираем
WORKDIR /app/services/gateway

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /main .
EXPOSE 8080
CMD ["./main"]