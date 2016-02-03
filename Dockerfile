FROM alpine:3.2

RUN apk add --update ca-certificates
ADD vault-monkey /app/

ENTRYPOINT ["/app/vault-monkey"]
