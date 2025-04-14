FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/main \
    ./cmd/app

FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache bash && \
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

COPY --from=builder /app/main /app/main
COPY --from=builder /app/go.mod /app/go.sum /app/
COPY --from=builder /app/cmd /app/cmd
COPY --from=builder /app/internal /app/internal
COPY --from=builder /app/.golangci.yml /app/.golangci.yml
COPY --from=builder /app/.env /app/.env

SHELL ["/bin/bash", "-c"]

CMD ["/app/main"]
