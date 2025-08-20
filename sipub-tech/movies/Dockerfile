# --- Build Stage ---
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o movies_service -trimpath -ldflags '-w -s' ./main.go

# --- Final Stage ---
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/movies_service /app/movies_service

ENTRYPOINT ["/app/movies_service"]
