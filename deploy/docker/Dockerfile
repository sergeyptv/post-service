FROM golang:1.26.0 AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG SERVICE

RUN test -n "$SERVICE" \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/app ./${SERVICE}/cmd/main.go

FROM alpine:3.23

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /out/app /app/app

COPY migrations /app/migrations

USER appuser

ENTRYPOINT ["/app/app"]