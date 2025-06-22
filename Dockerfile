# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

# Runtime stage
FROM alpine:3.20
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
