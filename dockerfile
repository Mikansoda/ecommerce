# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main .

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8007
CMD ["./main"]
