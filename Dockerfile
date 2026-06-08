# Build stage
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /app/mapminder-backend ./cmd/main.go
FROM alpine:3.19
RUN apk add --no-cache ca-certificates dumb-init
WORKDIR /app

COPY --from=builder /app/mapminder-backend .
COPY .env .
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080
ENTRYPOINT ["dumb-init", "--"]

CMD ["./mapminder-backend"]
