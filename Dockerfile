FROM golang:1.15.6-alpine

# Build the app, dependencies first
RUN apk add --no-cache git
RUN  go get github.com/markbates/pkger/cmd/pkger
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download

COPY . /app
RUN pkger
ENV CGO_ENABLED=0
RUN go build -o main
RUN go test ./...

# ---
FROM alpine:3.11 AS dist

# Dependencies
RUN apk add --no-cache ca-certificates

# Add pre-built application
COPY --from=0 /app/main /app

EXPOSE 3000
ENTRYPOINT [ "/app" ]