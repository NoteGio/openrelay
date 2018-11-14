FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/affiliatemonitor /affiliatemonitor

COPY --from=corebuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/affiliatemonitor", "redis:6379", "ethnode:8545", "queue://newblocks", "0x48bacb9266a570d521063ef5dd96e61686dbe788"]
