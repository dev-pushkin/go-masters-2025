FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./review_manager/cmd/service/*

FROM debian:bullseye-slim
COPY --from=builder /app/server /server

