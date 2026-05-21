FROM golang:1.24-alpine AS builder

ENV GOTOOLCHAIN=auto

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8081
CMD ["./server"]
