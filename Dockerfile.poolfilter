FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/poolfilter /poolfilter

COPY --from=corebuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/poolfilter", "postgres://postgres@postgres", "env://PG_SECRET", "redis:6379", "ethnode:8545", "queue://poolfilter=>queue://released"]
