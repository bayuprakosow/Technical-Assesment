# build
FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# run
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -H -u 65532 nonroot
WORKDIR /app
COPY --from=build /out/server /app/server
COPY migrations /app/migrations
USER nonroot:nonroot
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/health || exit 1
ENTRYPOINT ["/app/server"]
