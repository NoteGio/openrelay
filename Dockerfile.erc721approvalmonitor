FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/erc721approvalmonitor /erc721approvalmonitor

COPY --from=corebuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/erc721approvalmonitor", "redis:6379", "ethnode:8545", "queue://newblocks", "queue://recordspend", "0x48bacb9266a570d521063ef5dd96e61686dbe788"]
