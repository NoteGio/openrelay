FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/indexer /indexer

COPY --from=corebuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/indexer", "redis:6379", "queue://indexer", "postgres://postgres@postgres", "/run/secrets/postgress_password"]
