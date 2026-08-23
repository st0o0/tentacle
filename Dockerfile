FROM golang:1.26.6-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /tentacle ./cmd/tentacle

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /tentacle /tentacle

EXPOSE 9594

HEALTHCHECK --interval=30s --timeout=10s \
           --start-period=15s --retries=3 \
  CMD ["/tentacle", "healthcheck"]

ENTRYPOINT ["/tentacle"]
