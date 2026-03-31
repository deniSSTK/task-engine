FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY libs/go.mod libs/go.sum ./libs/
COPY gen/go.mod gen/go.sum ./gen/
COPY gateway/go.mod gateway/go.sum ./gateway/

COPY libs ./libs
COPY gen ./gen
COPY gateway ./gateway

RUN go work init ./gateway ./gen ./libs

WORKDIR /app/gateway
RUN go build -o /gateway .

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /gateway /usr/local/bin/gateway
COPY --from=builder /app/gen/openapiv2 ./openapi

EXPOSE 8080

CMD ["gateway"]
