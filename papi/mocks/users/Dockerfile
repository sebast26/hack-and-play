FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY main.go .
RUN go build -o users-service main.go

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/users-service .
COPY openapi.yaml .

EXPOSE 8082
CMD ["./users-service"]
