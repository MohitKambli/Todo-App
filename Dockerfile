FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/main /app/todo-app

EXPOSE 8080

CMD ["./todo-app"]