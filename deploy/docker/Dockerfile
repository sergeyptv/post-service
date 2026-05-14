FROM golang:1.26.3-alpine3.23 AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG SERVICE

RUN test -n "$SERVICE" \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/app ./"$SERVICE"

FROM alpine:3.23

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /out/app /app/app

COPY migrations /app/migrations

COPY certs /app/certs

USER appuser

ENTRYPOINT ["/app/app"]