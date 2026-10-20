
# Stage 1: Build the binary
FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

# Stage 2: Run the binary
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app ./app

EXPOSE 8080

CMD ["./app"]
