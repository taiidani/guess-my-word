FROM golang:1.26-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY --parents **/*.go app internal ./

RUN go build -o /guess-my-word main.go

# Deployable image
FROM alpine:latest

RUN adduser -D appuser
USER appuser

COPY --from=builder /guess-my-word /guess-my-word

EXPOSE 3000
ENTRYPOINT ["/guess-my-word"]
