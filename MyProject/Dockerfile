FROM golang:1.26 AS builder

WORKDIR /app


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o /crypto_bot ./cmd/crypto_bot


FROM alpine

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /crypto_bot .
COPY deploy/config/config.yaml ./config.yaml

EXPOSE 8080

CMD ["./crypto_bot"]