# syntax=docker/dockerfile:1.4
FROM --platform=$BUILDPLATFORM golang:1.18 AS build
WORKDIR /src

# Copy over our mod and sum files, then run go mod download so we can cache the dependancies as three container layers
COPY --link go.mod go.sum ./
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go mod download

# Now copy over the rest of the source code
COPY --link . .

# Now build the app for the target OS / architecture
ARG TARGETOS
ARG TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/server ./cmd/server

# Now we build the final image
FROM gcr.io/distroless/static
COPY --from=build /out/server /

LABEL org.opencontainers.image.authors="support@encore.dev" \
		org.opencontainers.image.vendor="Encoretivity AB" \
		org.opencontainers.image.title="Placeholder Service" \
		org.opencontainers.image.source="https://github.com/encoredev/placeholder-service" \
		org.opencontainers.image.description="This image is initially deployed into newly provisioned infrastructure as a placeholder servier which can respond to healthz requests while we build the encore service"

ENV HTTP_PORT=80

EXPOSE 80/tcp

USER nonroot:nonroot

ENTRYPOINT ["/server"]
