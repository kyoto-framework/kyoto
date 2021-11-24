# -------------
# build stage
# -------------
FROM python:3-alpine AS build

# Set workdir
WORKDIR /docs

# Python deps
RUN pip3 install mkdocs-material

# Attach sources
ADD . /docs

# Build
RUN mkdocs build

# -------------
# runtime stage
# -------------
FROM caddy:alpine AS runtime

# Set workdir
WORKDIR /usr/share/caddy

# Copy docs to serving directory
COPY --from=build /docs/site/ ./
