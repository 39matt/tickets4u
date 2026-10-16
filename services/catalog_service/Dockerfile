FROM golang:1.27.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags="-w -s" \
    -o /app/bin/catalog_service \
    ./cmd/api

FROM scratch

COPY --from=builder /app/bin/catalog_service /catalog_service

EXPOSE 8080

ENTRYPOINT ["/catalog_service"]