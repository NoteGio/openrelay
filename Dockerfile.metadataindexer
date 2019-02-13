FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/metadataindexer /metadataindexer

COPY --from=corebuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/metadataindexer", "redis:6379", "queue://metadataindexer", "postgres://metadata@postgres", "secret", "http://ethnode:8545"]
