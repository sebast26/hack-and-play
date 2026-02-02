FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY main.go .
RUN go build -o orders-service main.go

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/orders-service .
COPY openapi.yaml .

EXPOSE 8081
CMD ["./orders-service"]
