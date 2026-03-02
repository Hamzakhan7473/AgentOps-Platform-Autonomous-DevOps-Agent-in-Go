# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
# Copy go.mod first so download works even without go.sum (e.g. fresh clone)
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /agent ./cmd/agent

# Runtime stage — single binary
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /agent .
ENTRYPOINT ["./agent"]
