FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o invoice-generator ./cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/invoice-generator .
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./invoice-generator"]
