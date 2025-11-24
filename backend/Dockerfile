# Builder layer
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o jmr_backend .

# Runner layer
FROM alpine:latest
RUN apk add --no-cache rsync
WORKDIR /app
COPY --from=builder /app/jmr_backend ./
EXPOSE 7749
CMD ["./jmr_backend"]


