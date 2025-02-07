# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o tofu-service ./cmd/server

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/tofu-service .
COPY --from=builder /app/config.yaml .
EXPOSE 8080
CMD ["./tofu-service", "-jwt-secret", "$JWT_SECRET"] 