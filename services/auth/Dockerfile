FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY libs/go.mod libs/go.sum ./libs/
COPY gen/go.mod gen/go.sum ./gen/
COPY services/auth/go.mod services/auth/go.sum ./services/auth/

COPY libs ./libs
COPY gen ./gen
COPY services/auth ./services/auth

RUN go work init ./services/auth ./libs ./gen

WORKDIR /app/services/auth
RUN go build -o /auth-service .

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /auth-service /usr/local/bin/auth-service

EXPOSE 50001

CMD ["auth-service"]
