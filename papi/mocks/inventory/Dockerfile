FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY main.go .
RUN go build -o inventory-service main.go

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/inventory-service .

EXPOSE 8083
CMD ["./inventory-service"]
