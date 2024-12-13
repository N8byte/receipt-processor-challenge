# Build stage
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o receipt-processor

# Run stage
FROM alpine:latest
LABEL name="receipt-processor-challenge"

WORKDIR /root/

COPY --from=builder /app/receipt-processor .

EXPOSE 8080

CMD ["./receipt-processor"]