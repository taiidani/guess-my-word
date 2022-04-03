FROM node:17-alpine AS build

COPY web/package*.json /app/
WORKDIR /app
RUN npm install --frozen-lockfile

COPY web/ /app
RUN npm run build

FROM golang:1.18-alpine AS base

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
FROM nginx:1-alpine AS dist

# Add pre-built application
COPY --from=base /app/main /app
COPY --from=build /app/dist /usr/share/nginx/html/dist
COPY --from=build /app/assets /usr/share/nginx/html/assets
COPY --from=build /app/index.html /usr/share/nginx/html/index.html

ENV GIN_MODE="release"
EXPOSE 3000
EXPOSE 80
LABEL org.opencontainers.image.source=https://github.com/taiidani/guess-my-word

FROM dist AS dev

COPY nginx.dev.conf /etc/nginx/conf.d/default.conf
