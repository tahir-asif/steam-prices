# Build stage
FROM golang:1.25-alpine3.21 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Run stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 3000
CMD ["./api"]
