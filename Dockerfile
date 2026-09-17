# Сборка Go-приложения
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники и собираем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Минимальный финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и статику из стадии сборки
COPY --from=builder /app/main .
COPY --from=builder /app/index.html .

EXPOSE 8080

CMD ["./main"]