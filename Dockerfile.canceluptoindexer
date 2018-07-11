FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/canceluptoindexer /canceluptoindexer

COPY --from=corebuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/canceluptoindexer", "redis:6379", "queue://pgordersfilled", "postgres://postgres@postgres", "/run/secrets/postgress_password"]
