FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN GONOSUMCHECK=* GOFLAGS="-mod=mod" go mod download

COPY . .
RUN go build -o server ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8081
CMD ["./server"]
