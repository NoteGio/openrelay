FROM corebuild

FROM scratch

COPY --from=corebuild /go/src/github.com/notegio/openrelay/bin/multisigmonitor /multisigmonitor

CMD ["/multisigmonitor", "redis:6379", "ethnode:8545", "queue://newblocks", "0x48bacb9266a570d521063ef5dd96e61686dbe788"]
