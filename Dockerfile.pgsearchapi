FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/searchapi /searchapi

COPY --from=corebuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/searchapi", "redis:6379", "topic://newblocks", "postgres://postgres@postgres", "secret"]
