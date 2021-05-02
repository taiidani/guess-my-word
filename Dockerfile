FROM golang:1.16.2-alpine AS base

# Build the app, dependencies first
RUN apk add --no-cache git
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download

COPY . /app
ENV CGO_ENABLED=0
RUN go build -o main

# ---
FROM base AS test

RUN go test ./...

# ---
FROM alpine:3.13.5 AS dist

# Dependencies
RUN apk add --no-cache ca-certificates

# Add pre-built application
COPY --from=base /app/main /app

EXPOSE 3000
ENTRYPOINT [ "/app" ]
LABEL org.opencontainers.image.source=https://github.com/taiidani/guess-my-word
