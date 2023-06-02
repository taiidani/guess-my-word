FROM alpine:3.18

RUN apk add --no-cache ca-certificates

# Add pre-built application
COPY guess_my_word /app
RUN /app --help

ENV GIN_MODE="release"
EXPOSE 3000
EXPOSE 80
LABEL org.opencontainers.image.source=https://github.com/taiidani/guess-my-word
