# -------------
# build stage
# -------------
FROM golang:alpine AS build

# Attach sources
WORKDIR /src
ADD . /src

# System deps
RUN apk add --no-cache git npm

# Node deps
RUN npm i -g vuepress vuepress-theme-default-prefers-color-scheme

# Build
RUN go build -o ssceng-web
RUN (cd static; npm i; npm run build)
RUN vuepress build docs

# -------------
# runtime stage
# -------------
FROM alpine

# Copy app
WORKDIR /app
COPY --from=build /src/ssceng-web /app/
COPY --from=build /src/*.html /app/
COPY --from=build /src/static /app/static
COPY --from=build /src/docs /app/docs

# Entrypoint
ENTRYPOINT PORT=25025 ./ssceng-web
