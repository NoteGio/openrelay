FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/blockmonitor /blockmonitor

COPY --from=corebuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/blockmonitor", "redis:6379", "ethnode:8545", "queue://newblocks"]
