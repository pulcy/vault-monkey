FROM alpine:3.2

ADD vault-monkey /app/

ENTRYPOINT ["/app/vault-monkey"]
