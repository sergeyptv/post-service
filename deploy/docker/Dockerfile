FROM golang:1.26.3-bookworm AS builder

WORKDIR /src

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates git librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG SERVICE

ENV CGO_CFLAGS="-I/usr/include"
ENV CGO_LDFLAGS="-L/usr/lib -lrdkafka"

RUN test -n "$SERVICE" \
    && CGO_ENABLED=1 GOOS=linux \
    go build -tags dynamic -trimpath -ldflags="-s -w" -o /out/app ./"$SERVICE"

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates tzdata librdkafka-dev \
    && rm -rf /var/lib/apt/lists/* \
    && useradd -m -s /bin/false appuser

WORKDIR /app

COPY --from=builder /out/app /app/app

COPY migrations /app/migrations

COPY certs /app/certs

USER appuser

ENTRYPOINT ["/app/app"]