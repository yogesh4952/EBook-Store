FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -0 /server ./cmd/api/

FROM alpine:3.22

WORKDIR /app 

COPY --from=builder /server ./server

EXPOSE 8080

CMD ["./server"]
