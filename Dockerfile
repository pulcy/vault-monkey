FROM alpine:3.4

RUN apk add --update ca-certificates
ADD bin/vault-monkey-linux-amd64 /app/vault-monkey

ENTRYPOINT ["/app/vault-monkey"]
