FROM node:17-alpine AS build

COPY web/package*.json /app/
WORKDIR /app
RUN npm install --frozen-lockfile

COPY web/ /app
RUN npm run build

# ---
FROM nginx:1-alpine AS dist

# Add pre-built application
COPY guess-my-word /app
COPY --from=build /app/dist /usr/share/nginx/html

ENV GIN_MODE="release"
EXPOSE 3000
EXPOSE 80
LABEL org.opencontainers.image.source=https://github.com/taiidani/guess-my-word

FROM dist AS dev

COPY nginx.dev.conf /etc/nginx/conf.d/default.conf
