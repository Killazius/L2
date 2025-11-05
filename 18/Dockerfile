FROM golang:1.25.1-alpine AS builder

WORKDIR /cmd
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app ./cmd/




FROM alpine:3.22 AS app
WORKDIR /app
COPY --from=builder /cmd/app .
RUN mkdir -p /config
COPY --from=builder /cmd/config/ ./config/


EXPOSE 8080

CMD ["./app"]
