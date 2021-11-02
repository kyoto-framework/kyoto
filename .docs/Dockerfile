# -------------
# build stage
# -------------
FROM alpine:latest AS build

# Set workdir
WORKDIR /docs

# System deps
RUN apk add --no-cache npm

# Node deps
ADD ./package.json ./
RUN npm i

# Attach sources
ADD . /docs

# Build
RUN npm run build

# -------------
# runtime stage
# -------------
FROM caddy:alpine

# Set workdir
WORKDIR /usr/share/caddy

# Copy docs to serving directory
COPY --from=build /docs/.vuepress/dist/ ./
