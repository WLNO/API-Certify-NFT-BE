# syntax=docker/dockerfile:1
FROM golang:1.23.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o app main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/.gitignore .
EXPOSE 4002
CMD ["./app"] 